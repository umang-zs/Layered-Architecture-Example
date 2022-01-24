package store

import "github.com/umang01-hash/layered-architecture/model"

type Store interface {
	GetByID(id int) (model.Customer, error)
	Create(c model.Customer) (model.Customer, error)
	Delete(id string) error
	Update(c model.Customer) error
}
