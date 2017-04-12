package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database struct {
	sync.Mutex
	items map[string]dollars
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db.items {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db *database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db.items[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db *database) update(w http.ResponseWriter, req *http.Request) {
	item, price, err := getParams(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	if _, ok := db.items[item]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	db.Lock()
	db.items[item] = price
	db.Unlock() // Don't use defer to release lock earlier.

	w.WriteHeader(http.StatusOK)
}

func (db *database) create(w http.ResponseWriter, req *http.Request) {
	item, price, err := getParams(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	if _, ok := db.items[item]; ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "item already exists %q\n", item)
		return
	}
	db.Lock()
	db.items[item] = price
	db.Unlock() // Don't use defer to release lock earlier.

	w.WriteHeader(http.StatusOK)
}

func (db *database) remove(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db.items[item]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %q\n", item)
		return
	}
	db.Lock()
	delete(db.items, item)
	db.Unlock() // Don't use defer to release lock earlier.

	w.WriteHeader(http.StatusOK)
}

func getParams(req *http.Request) (string, dollars, error) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	priceFloat, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return item, dollars(0), fmt.Errorf("invalid value for price: '%s'", price)
	}
	return item, dollars(priceFloat), nil
}

func main() {
	db := database{items: map[string]dollars{"shoes": 50, "socks": 5}}
	// Don't use here PUT, POST, DELETE due to task requirements.
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/remove", db.remove)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
