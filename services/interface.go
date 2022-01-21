package services

import "github.com/umang01-hash/layered-architecture/model"

type Service interface {
	Get(id int) (model.Customer, error)
	Create(c model.Customer) (model.Customer, error)
	Update(c model.Customer) error
	Delete(id int) error
}
