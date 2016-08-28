package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// DataAccessLayer ...
type DataAccessLayer interface {
	FindDishes() []Dish
	GetDishByID(string) Dish
}

// DAL ...
type DAL struct {
	db *sql.DB
}

// Close ...
func (dal DAL) Close() {
	dal.db.Close()
}

// FindDishes ...
func (dal DAL) FindDishes() []Dish {
	var dishes = []Dish{}
	rows, ErrSQLQuery := dal.db.Query("SELECT name, description, price, id FROM dish")
	if ErrSQLQuery != nil {
		log.Fatal("ErrSQLQuery => ", ErrSQLQuery)
	}
	defer rows.Close()
	for rows.Next() {
		var name, description, id string
		var price float32
		if err := rows.Scan(&name, &description, &price, &id); err != nil {
			log.Fatal(err)
		}
		dishes = append(dishes, Dish{
			Name:        name,
			Price:       price,
			Description: description,
			ID:          id,
		})
	}
	return dishes
}

// GetDishByID ...
func (dal DAL) GetDishByID(dishID string) Dish {
	var name, description, id string
	var price float32
	dal.
		db.
		QueryRow("SELECT name, description, price, id FROM dish WHERE id=$1", dishID).
		Scan(&name, &description, &price, &id)
	return Dish{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
	}
}

// NewDAL ...
func NewDAL(url string) DAL {
	db, errSQLOpen := sql.Open("postgres", url)
	if errSQLOpen != nil {
		log.Fatal("error occurs on connect postgres", errSQLOpen)
	}
	return DAL{db}
}
