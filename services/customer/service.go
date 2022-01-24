package customer

import (
	"github.com/umang01-hash/layered-architecture/model"
	"github.com/umang01-hash/layered-architecture/store"
)

type service struct {
	store store.Store
}

func New(store store.Store) service {
	return service{store: store}
}

func (s service) GetByID(id int) (model.Customer, error) {

	return s.store.GetByID(id)
}

func (s service) Create(c model.Customer) (model.Customer, error) {
	return s.store.Create(c)
}

func (s service) Update(c model.Customer) error {
	return s.store.Update(c)
}

func (s service) Delete(id string) error {
	return s.store.Delete(id)
}
