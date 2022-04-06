package services

import (
	"encoding/json"
	"log"

	"L0/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/nats-io/stan.go"
)

type Service struct {
	Client stan.Conn
	Sub    stan.Subscription
	cache  *Cache
	db     *DB
	valid  *validator.Validate
}

func NewService(cache *Cache, db *DB, valid *validator.Validate) *Service {
	return &Service{
		Client: nil,
		Sub:    nil,
		cache:  cache,
		db:     db,
		valid:  valid,
	}
}

func (s *Service) ConnectionToNats() (err error) {
	s.Client, err = stan.Connect("test-cluster", "sub")
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Subscribe() (err error) {
	s.Sub, err = s.Client.Subscribe("foo", s.handler)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) handler(m *stan.Msg) {
	var order models.Order

	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		log.Println("wrong data")
		return
	}

	err = s.valid.Struct(order)
	if err != nil {
		log.Println("wrong data")
		return
	}
	
	err = s.cache.Set(order)
	if err != nil {
		log.Println(err.Error())
		return
	}

	go func() {
		err = s.db.InsertOrder(order)
		if err != nil {
			log.Println(err)
			return
		}
	}()
}
