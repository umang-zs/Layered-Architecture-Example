package store

import "github.com/umang01-hash/layered-architecture/model"

type Store interface {
	Get(id int ) (model.Customer,error)
	Create(c model.Customer) (model.Customer,error)
	Delete(id int) error
	Update(c model.Customer) error
}
