package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
)

// CreateRouter ...
func CreateRouter(dal DataAccessLayer) *httprouter.Router {
	router := httprouter.New()
	router.GET("/", Index(dal))
	router.GET("/healthz", HealthCheck(dal))
	// router.GET("/menus", GetMenus)

	router.GET("/dishes", GetDishes(dal))
	router.GET("/dishes/:id", GetDish(dal))
	router.POST("/dishes", CreateDish(dal))
	router.PUT("/dishes", UpdateDish(dal))
	router.DELETE("/dishes/:id", DeleteDish(dal))
	return router
}

// Index is welcome page
func Index(dal DataAccessLayer) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		req *http.Request,
		_ httprouter.Params,
	) {
		fmt.Fprintf(w, "OK")
	}
}

// HealthCheck used to provide health check without cache
func HealthCheck(dal DataAccessLayer) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		req *http.Request,
		_ httprouter.Params,
	) {
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

// UpdateDish ...
func UpdateDish(dal DataAccessLayer) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		req *http.Request,
		p httprouter.Params,
	) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}
		dish := Dish{}
		json.Unmarshal(body, &dish)
		err = dal.UpdateDish(&dish)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprint(w, "ok")
	}
}

// GetDish ...
func GetDish(dal DataAccessLayer) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		req *http.Request,
		p httprouter.Params,
	) {
		h := w.Header()
		h.Add("Content-Type", "application/json")
		pID := p.ByName("id")
		dish := dal.GetDishByID(pID)
		bs, _ := json.Marshal(dish)
		fmt.Fprintf(w, "%s", bs)
	}
}

// GetDishes ...
func GetDishes(dal DataAccessLayer) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		req *http.Request,
		_ httprouter.Params,
	) {
		h := w.Header()
		h.Add("Content-Type", "application/json")
		dishes := dal.FindDishes()
		bs, _ := json.Marshal(dishes)
		fmt.Fprintf(w, "%s", bs)
	}
}

// CreateDish ...
func CreateDish(dal DataAccessLayer) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		req *http.Request,
		_ httprouter.Params,
	) {
		body, ErrReadAll := ioutil.ReadAll(req.Body)
		if ErrReadAll != nil {
			log.Println("error occurs on parsing body", ErrReadAll)
		}
		dish := Dish{}
		if err := json.Unmarshal(body, &dish); err != nil {
			panic(err)
		}
		u1 := uuid.NewV4()
		dish.ID = u1.String()
		err := dal.CreateDish(&dish)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, u1.String())
	}
}

// DeleteDish ...
func DeleteDish(dal DataAccessLayer) httprouter.Handle {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		p httprouter.Params,
	) {
		id := p.ByName("id")
		err := dal.DeleteDishByID(id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "OK")
	}
}
