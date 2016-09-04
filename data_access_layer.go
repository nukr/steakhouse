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
	CreateDish(*Dish) error
	DeleteDishByID(string) error
	UpdateDish(*Dish) error
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

// CreateDish ...
func (dal DAL) CreateDish(dish *Dish) error {
	_, err := dal.db.
		Exec(`INSERT INTO dish (name, description, price, id) VALUES($1, $2, $3, $4)`,
			dish.Name,
			dish.Description,
			dish.Price,
			dish.ID,
		)
	return err
}

// UpdateDish ...
func (dal DAL) UpdateDish(dish *Dish) error {
	_, err := dal.db.Exec(`UPDATE dish SET (name, description, price) = ($1, $2, $3) WHERE id = $4`,
		dish.Name,
		dish.Description,
		dish.Price,
		dish.ID,
	)
	return err
}

// DeleteDishByID ...
func (dal DAL) DeleteDishByID(id string) error {
	_, err := dal.db.Exec(`DELETE FROM dish WHERE id=$1`, id)
	return err
}

// NewDAL ...
func NewDAL(url string) DAL {
	db, errSQLOpen := sql.Open("postgres", url)
	if errSQLOpen != nil {
		log.Fatal("error occurs on connect postgres", errSQLOpen)
	}
	return DAL{db}
}
