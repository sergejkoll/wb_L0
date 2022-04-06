package services

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func (s *Service) getOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := s.cache.Get(id)
	if err != nil {
		return fiber.ErrBadRequest
	}
	return c.Render("order", fiber.Map{
		"OrderUID":           order.OrderUID,
		"track_number":       order.TrackNumber,
		"entry":              order.Entry,
		"locale":             order.Locale,
		"internal_signature": order.InternalSignature,
		"customer_id":        order.CustomerId,
		"delivery_service":   order.DeliveryService,
		"shardkey":           order.ShardKey,
		"sm_id":              order.SmId,
		"date_created":       order.DateCreated,
		"oof_shard":          order.OofShard,

		"name":    order.Delivery.Name,
		"phone":   order.Delivery.Phone,
		"zip":     order.Delivery.Zip,
		"city":    order.Delivery.City,
		"address": order.Delivery.Address,
		"region":  order.Delivery.Region,
		"email":   order.Delivery.Email,

		"transaction":   order.Payment.Transaction,
		"request_id":    order.Payment.RequestId,
		"currency":      order.Payment.Currency,
		"provider":      order.Payment.Provider,
		"amount":        order.Payment.Amount,
		"payment_dt":    order.Payment.PaymentDt,
		"bank":          order.Payment.Bank,
		"delivery_cost": order.Payment.DeliveryCost,
		"goods_total":   order.Payment.GoodsTotal,
		"custom_fee":    order.Payment.CustomFee,

		"items": order.Items,
	})
}

func (s *Service) Run() {
	engine := html.New("./web", ".html")

	server := fiber.New(fiber.Config{
		Views: engine,
	})

	server.Get("/api/v1/:id", s.getOrder)

	err := server.Listen(":3000")
	if err != nil {
		log.Fatal(err)
		return
	}
}
