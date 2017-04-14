package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var templ = template.Must(template.New("itemlist").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <title>Items list</title>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css">
  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.2.0/jquery.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
</head>
<body>
<div class="container">
  <h2>Items with prices</h2>
  <table class="table table-bordered table-striped table-hover">
    <thead class="thead-inverse">
      <tr>
        <th>Item</th>
        <th>Price</th>
      </tr>
    </thead>
    <tbody>
    {{range $item, $price := .}}
    <tr>
      <td>{{ $item }}</td><td>{{ $price | printf "%4.2f"}}</td>
    </tr>
    {{end}}
    </tbody>
  </table>
</div>
</body>
</html>
`))

type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database struct {
	sync.Mutex
	items map[string]dollars
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {
	if err := templ.Execute(w, db.items); err != nil {
		log.Fatal(err)
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

	w.WriteHeader(http.StatusNoContent)
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

	w.WriteHeader(http.StatusNoContent)
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

	w.WriteHeader(http.StatusNoContent)
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
