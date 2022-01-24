package customer

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/umang01-hash/layered-architecture/errors"
	"github.com/umang01-hash/layered-architecture/model"
)

type mockService struct{}
type mockReader struct{}

func (m mockReader) Read(p []byte) (n int, err error) {
	return 0, errors.TestError{}
}

func (m mockService) Create(c model.Customer) (model.Customer, error) {
	switch c.Name {
	case "Umang":
		c.ID = 1
		return c, nil
	case "Manav":
		return model.Customer{}, errors.TestError{}
	default:
		return model.Customer{}, nil

	}
}

func (m mockService) GetByID(id int) (model.Customer, error) {
	switch id {
	case 1:
		return model.Customer{ID: 1, Name: "Umang", PhoneNo: "8527029558", Address: "Shakti Khand 2 Indirapuram UP-201014"}, nil
	case 15:
		return model.Customer{}, errors.TestError{}
	default:
		return model.Customer{}, nil

	}
}

func (m mockService) Update(c model.Customer) error {
	return nil
}

func (m mockService) Delete(id string) error {

	switch id {
	case "1":
		return nil
	case "aab":
		return errors.TestError{}
	case "15":
		return errors.TestError{}
	default:
		return nil

	}

}

func TestHandler_Create(t *testing.T) {
	h := New(mockService{})

	cases := []struct {
		desc       string
		body       io.Reader
		statusCode int
		res        []byte
	}{
		{"Success", bytes.NewReader([]byte(`{"id":1,"name":"Umang","phone":"8527029558","address":"Shakti Khand 2 Indirapuram UP-201014"}`)), http.StatusCreated, nil},
		{"Missing or Invalid Parameter", bytes.NewReader([]byte(`{"name":"Manav","phone":"9953113063"}`)), http.StatusBadRequest, nil},
		{"Reader Error", mockReader{}, http.StatusInternalServerError, nil},
		{"Unmarshal Error", bytes.NewReader([]byte(`"invalid body"`)), http.StatusBadRequest, nil},
	}

	for i, tc := range cases {
		r := httptest.NewRequest(http.MethodPost, "http://dummy", tc.body)
		w := httptest.NewRecorder()

		h.Create(w, r)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {

			t.Errorf("[TEST%d] Failed Desc : %v. Got %v\tExpected %v\n", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}

func TestHandler_GetByID(t *testing.T) {
	h := New(mockService{})

	cases := []struct {
		desc       string
		id         string
		statusCode int
		res        []byte
	}{
		{"invalid id", "abc", http.StatusBadRequest, []byte(``)},
		{"Id Not Found", "15", http.StatusNotFound, []byte(``)},
		{"Success", "1", http.StatusOK, []byte(`{"id":"1","name":"Umang","phone":"8527029558","address":"Shakti Khand 2 Indirapuram UP"}`)},
	}

	for i, tc := range cases {
		r, _ := http.NewRequest(http.MethodGet, "http://dummy", nil)
		r = mux.SetURLVars(r, map[string]string{"id": tc.id})

		w := httptest.NewRecorder()

		h.GetByID(w, r)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {

			t.Errorf("[TEST%d] Failed Desc : %v. Got %v\tExpected %v\n", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}

}

func TestHandler_DeleteByID(t *testing.T) {
	h := New(mockService{})

	cases := []struct {
		desc       string
		id         string
		statusCode int
	}{
		{"Success", "1", http.StatusNoContent},
		{"Invalid Id", "aab", http.StatusBadRequest},
		{"Server Error", "15", http.StatusInternalServerError},
	}

	for i, tc := range cases {
		r, _ := http.NewRequest(http.MethodGet, "http://dummy", nil)
		r = mux.SetURLVars(r, map[string]string{"id": tc.id})

		w := httptest.NewRecorder()

		h.DeleteByID(w, r)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {

			t.Errorf("[TEST%d] Failed Desc : %v. Got %v\tExpected %v\n", i, tc.desc, resp.StatusCode, tc.statusCode)
		}
	}
}
