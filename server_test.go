package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server *httptest.Server
	reader io.Reader
)

type FakeDAL struct{}

func (dal FakeDAL) FindDishes() []Dish {
	dishes := []Dish{}
	return dishes
}

func (dal FakeDAL) GetDishByID(id string) Dish {
	return Dish{}
}

func (dal FakeDAL) CreateDish(d *Dish) error {
	return nil
}

func (dal FakeDAL) DeleteDishByID(id string) error {
	return nil
}

func (dal FakeDAL) UpdateDish(d *Dish) error {
	return nil
}

func init() {
	fakeDAL := FakeDAL{}
	server = httptest.NewServer(CreateRouter(fakeDAL))
}

func TestIndex(t *testing.T) {
	res, err := http.Get(server.URL)
	if err != nil {
		t.Error("fail")
	}
	if res.StatusCode != 200 {
		t.Errorf("fail %d", res.StatusCode)
	}
}

func TestHealthCheck(t *testing.T) {
	healthzURL := fmt.Sprintf("%s/healthz", server.URL)
	res, _ := http.Get(healthzURL)
	if res.StatusCode != http.StatusOK {
		t.Errorf("fail %d", res.StatusCode)
	}
	expected := "no-cache, no-store, must-revalidate"
	header := res.Header.Get("Cache-Control")
	if header != expected {
		t.Errorf("we expected %s but got %s", expected, header)
	}

	expected = "no-cache"
	header = res.Header.Get("Pragma")
	if header != expected {
		t.Errorf("we expected %s but got %s", expected, header)
	}

	expected = "0"
	header = res.Header.Get("Expires")
	if header != expected {
		t.Errorf("we expected %s but got %s", expected, header)
	}
}

// func TestGetDishes(t *testing.T) {
// 	getDishesURL := fmt.Sprintf("%s/dishes", server.URL)
// 	expectedContentType := "application/json"
// 	res, _ := http.Get(getDishesURL)
// 	if res.StatusCode != http.StatusOK {
// 		t.Errorf("expected StatusCode is %d, but we got %d", http.StatusOK, res.StatusCode)
// 	}

// 	resContentType := res.Header.Get("Content-Type")
// 	if resContentType != expectedContentType {
// 		t.Errorf("expected content type is %s, but we got %s", expectedContentType, resContentType)
// 	}
// }
