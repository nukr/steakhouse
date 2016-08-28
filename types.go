package main

// Dish ...
type Dish struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	ID          string  `json:"id"`
}

// Menu ...
type Menu struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Dishes      []Dish `json:"dishes"`
	ID          string `json:"id"`
}
