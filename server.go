package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"database/sql"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

const (
	port = ":3456"
)

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

var db *sql.DB
var errSQLOpen error

func main() {
	db, errSQLOpen = sql.Open("postgres", "postgres://postgres:secret@localhost:5432/steakhouse?sslmode=disable")
	if errSQLOpen != nil {
		log.Fatal("error occurs on connect postgres", errSQLOpen)
	}
	defer db.Close()
	router := CreateRouter()
	log.Fatal(http.ListenAndServe(port, router))
}

// CreateRouter ...
func CreateRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/healthz", HealthCheck)
	router.GET("/menus", GetMenus)

	router.GET("/dishes", GetDishes)
	router.GET("/dishes/:id", GetDish)
	router.POST("/dishes", CreateDish)
	router.PUT("/dishes/:id", UpdateDish)
	router.DELETE("/dishes/:id", DeleteDish)
	return router
}

// Index is welcome page
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "OK")
}

// HealthCheck used to provide health check without cache
func HealthCheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h := w.Header()
	h.Add("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Add("Pragma", "no-cache")
	h.Add("Expires", "0")
	fmt.Fprintf(w, "OK")
}

// GetMenus ...
func GetMenus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// To Be Implement
}

// UpdateDish ...
func UpdateDish(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// id := p.ByName("id")
}

// GetDish ...
func GetDish(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	h := w.Header()
	h.Add("Content-Type", "application/json")
	pID := p.ByName("id")
	var name, description, id string
	var price float32
	db.QueryRow(`
  SELECT
  name, description, price, id
  FROM dish
  WHERE id = $1`, pID).Scan(&name, &description, &price, &id)
	dish := Dish{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          id,
	}
	bs, _ := json.Marshal(dish)
	fmt.Fprintf(w, "%s", bs)
}

// GetDishes ...
func GetDishes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	h := w.Header()
	h.Add("Content-Type", "application/json")
	rows, ErrSQLQuery := db.Query("SELECT name, description, price, id FROM dish")
	if ErrSQLQuery != nil {
		log.Fatal("ErrSQLQuery => ", ErrSQLQuery)
	}
	defer rows.Close()
	dishes := []Dish{}
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
	bs, _ := json.Marshal(dishes)
	fmt.Fprintf(w, "%s", bs)
}

// CreateDish ...
func CreateDish(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, ErrReadAll := ioutil.ReadAll(r.Body)
	if ErrReadAll != nil {
		log.Println("error occurs on parsing body", ErrReadAll)
	}

	dish := Dish{}
	if err := json.Unmarshal(body, &dish); err != nil {
		panic(err)
	}
	u1 := uuid.NewV4()
	_, ErrSQLQuery := db.Exec(`
  INSERT INTO
  dish (name, description, price, id)
  VALUES($1, $2, $3, $4)`, dish.Name, dish.Description, dish.Price, u1)
	if ErrSQLQuery != nil {
		log.Fatal("ErrSQLQuery => ", ErrSQLQuery)
	}
	fmt.Fprintf(w, string(body))
}

// DeleteDish ...
func DeleteDish(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	_, ErrSQLQuery := db.Exec(`DELETE FROM dish WHERE id = $1`, id)
	if ErrSQLQuery != nil {
		log.Fatal(ErrSQLQuery)
	}
	fmt.Fprintf(w, "OK")
}
