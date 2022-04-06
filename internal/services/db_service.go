package services

import (
	"database/sql"
	"log"

	"L0/internal/models"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// check connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfully connected!")
	return &DB{db}, nil
}

func (db *DB) InsertOrder(order models.Order) (err error) {
	err = db.InsertOrderData(order)
	if err != nil {
		return err
	}

	err = db.InsertDeliveryData(order.OrderUID, order.Delivery)
	if err != nil {
		return err
	}

	err = db.InsertPaymentData(order.OrderUID, order.Payment)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		err = db.InsertItemData(order.OrderUID, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) InsertOrderData(data models.Order) (err error) {
	sqlQuery := "INSERT INTO order_data (order_uid, track_number, entry, locale, internal_signature, customer_id, " +
		"delivery_service, shard_key, sm_id, date_created, oof_shard) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	_, err = db.Exec(sqlQuery, data.OrderUID, data.TrackNumber, data.Entry, data.Locale, data.InternalSignature,
		data.CustomerId, data.DeliveryService, data.ShardKey, data.SmId, data.DateCreated,
		data.OofShard)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertDeliveryData(orderId string, deliveryData models.Delivery) (err error) {
	sqlQuery := "INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = db.Exec(sqlQuery, orderId, deliveryData.Name, deliveryData.Phone, deliveryData.Zip, deliveryData.City,
		deliveryData.Address, deliveryData.Region, deliveryData.Email)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertPaymentData(orderId string, paymentData models.Payment) (err error) {
	sqlQuery := "INSERT INTO payment (order_uid, transaction_id, request_id, currency, provider, amount, payment_dt, " +
		"bank, delivery_cost, goods_total, custom_fee) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	_, err = db.Exec(sqlQuery, orderId, paymentData.Transaction, paymentData.RequestId, paymentData.Currency,
		paymentData.Provider, paymentData.Amount, paymentData.PaymentDt, paymentData.Bank, paymentData.DeliveryCost,
		paymentData.GoodsTotal, paymentData.CustomFee)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertItemData(orderId string, item models.Item) (err error) {
	sqlQuery := "INSERT INTO item (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, " +
		"nm_id, brand, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
	_, err = db.Exec(sqlQuery, orderId, item.ChrtId, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale,
		item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetOrders() (orders []models.Order, err error) {
	sqlQuery := "SELECT * FROM order_data"
	rows, err := db.Query(sqlQuery)
	if err != nil {
		return []models.Order{}, err
	}
	defer rows.Close()

	for rows.Next() {
		order := models.Order{}
		err = rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerId, &order.DeliveryService, &order.ShardKey, &order.SmId, &order.DateCreated,
			&order.OofShard)
		if err != nil {
			return []models.Order{}, err
		}

		order.Delivery, err = db.GetDeliveryById(order.OrderUID)
		if err != nil {
			return []models.Order{}, err
		}

		order.Payment, err = db.GetPaymentById(order.OrderUID)
		if err != nil {
			return []models.Order{}, err
		}

		order.Items, err = db.GetItemsById(order.OrderUID)
		if err != nil {
			return []models.Order{}, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (db *DB) GetOrderById(orderId string) (order models.Order, err error) {
	sqlQuery := "SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, " +
		"shard_key, sm_id, date_created, oof_shard FROM order_data WHERE order_uid = $1"
	row := db.QueryRow(sqlQuery, orderId)
	err = row.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerId, &order.DeliveryService, &order.ShardKey, &order.SmId, &order.DateCreated, &order.OofShard)
	if err != nil {
		return models.Order{}, err
	}

	order.Delivery, err = db.GetDeliveryById(orderId)
	if err != nil {
		return models.Order{}, err
	}

	order.Payment, err = db.GetPaymentById(orderId)
	if err != nil {
		return models.Order{}, err
	}

	order.Items, err = db.GetItemsById(orderId)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (db *DB) GetDeliveryById(orderId string) (delivery models.Delivery, err error) {
	sqlQuery := "SELECT order_uid, name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1"
	row := db.QueryRow(sqlQuery, orderId)
	err = row.Scan(&orderId, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
		&delivery.Region, &delivery.Email)
	if err != nil {
		return models.Delivery{}, err
	}
	return delivery, nil
}

func (db *DB) GetPaymentById(orderId string) (payment models.Payment, err error) {
	sqlQuery := "SELECT order_uid, transaction_id, request_id, currency, provider, amount, payment_dt, " +
		"bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1"
	row := db.QueryRow(sqlQuery, orderId)
	err = row.Scan(&orderId, &payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider,
		&payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal,
		&payment.CustomFee)
	if err != nil {
		return models.Payment{}, err
	}
	return payment, err
}

func (db *DB) GetItemsById(orderId string) (items []models.Item, err error) {
	sqlQuery := "SELECT order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, " +
		"status FROM item WHERE order_uid = $1"
	rows, err := db.Query(sqlQuery, orderId)
	if err != nil {
		return []models.Item{}, err
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Item{}
		err = rows.Scan(&orderId, &item.ChrtId, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		if err != nil {
			return []models.Item{}, err
		}
		items = append(items, item)
	}

	return items, nil
}
