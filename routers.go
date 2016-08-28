package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// CreateRouter ...
func CreateRouter(dal DataAccessLayer) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index(dal))
	router.GET("/healthz", HealthCheck(dal))
	// router.GET("/menus", GetMenus)

	router.GET("/dishes", GetDishes(dal))
	router.GET("/dishes/:id", GetDish(dal))
	// router.POST("/dishes", CreateDish)
	// router.PUT("/dishes/:id", UpdateDish)
	// router.DELETE("/dishes/:id", DeleteDish)
	return router
}

// Index is welcome page
func Index(dal DataAccessLayer) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "OK")
	}
}

// HealthCheck used to provide health check without cache
func HealthCheck(dal DataAccessLayer) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h := w.Header()
		h.Add("Cache-Control", "no-cache, no-store, must-revalidate")
		h.Add("Pragma", "no-cache")
		h.Add("Expires", "0")
		fmt.Fprintf(w, "OK")
	}
}

// // GetMenus ...
// func GetMenus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	// To Be Implement
// }

// // UpdateDish ...
// func UpdateDish(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	// id := p.ByName("id")
// }

// GetDish ...
func GetDish(dal DataAccessLayer) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		h := w.Header()
		h.Add("Content-Type", "application/json")
		pID := p.ByName("id")
		dish := dal.GetDishByID(pID)
		bs, _ := json.Marshal(dish)
		fmt.Fprintf(w, "%s", bs)
	}
}

// GetDishes ...
func GetDishes(dal DataAccessLayer) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h := w.Header()
		h.Add("Content-Type", "application/json")
		dishes := dal.FindDishes()
		bs, _ := json.Marshal(dishes)
		fmt.Fprintf(w, "%s", bs)
	}
}

// // CreateDish ...
// func CreateDish(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	body, ErrReadAll := ioutil.ReadAll(r.Body)
// 	if ErrReadAll != nil {
// 		log.Println("error occurs on parsing body", ErrReadAll)
// 	}

// 	dish := Dish{}
// 	if err := json.Unmarshal(body, &dish); err != nil {
// 		panic(err)
// 	}
// 	u1 := uuid.NewV4()
// 	_, ErrSQLQuery := db.Exec(`INSERT INTO dish (name, description, price, id) VALUES($1, $2, $3, $4)`, dish.Name, dish.Description, dish.Price, u1)

// 	if ErrSQLQuery != nil {
// 		log.Fatal("ErrSQLQuery => ", ErrSQLQuery)
// 	}
// 	fmt.Fprintf(w, string(body))
// }

// // DeleteDish ...
// func DeleteDish(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
// 	id := p.ByName("id")
// 	_, ErrSQLQuery := db.Exec(`DELETE FROM dish WHERE id = $1`, id)
// 	if ErrSQLQuery != nil {
// 		log.Fatal(ErrSQLQuery)
// 	}
// 	fmt.Fprintf(w, "OK")
// }
