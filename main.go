package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Order struct {
	
	OrderID      uint      `json:"orderId" gorm:"primary_key"`
	CustomerName string    `json:"customerName"`
	OrderedAt    time.Time `json:"orderedAt"`
	Items        []Item    `json:"items" gorm:"foreignkey:OrderID"`
}


type Item struct {
	
	LineItemID  uint   `json:"lineItemId" gorm:"primary_key"`
	ItemCode    string `json:"itemCode"`
	Description string `json:"description"`
	Quantity    uint   `json:"quantity"`
	OrderID     uint   `json:"-"`
}

var db *gorm.DB

func initDB() {
	var err error
	dataSourceName := "root:@tcp(localhost:3306)/?parseTime=True"
	db, err = gorm.Open("mysql", dataSourceName)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	db.Exec("CREATE DATABASE orders")
	db.Exec("USE orders")

	db.AutoMigrate(&Order{}, &Item{})
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var orders []Order
	db.Preload("Items").Find(&orders)
	json.NewEncoder(w).Encode(orders)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputOrderID := params["orderId"]

	var order Order
	db.Preload("Items").First(&order, inputOrderID)
	json.NewEncoder(w).Encode(order)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	db.Create(&order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}


func updateOrder(w http.ResponseWriter, r *http.Request) {
	var updatedOrder Order
	json.NewDecoder(r.Body).Decode(&updatedOrder)
	db.Save(&updatedOrder)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inputOrderID := params["orderId"]
	id64, _ := strconv.ParseUint(inputOrderID, 10, 64)
	idToDelete := uint(id64)

	db.Where("order_id = ?", idToDelete).Delete(&Item{})
	db.Where("order_id = ?", idToDelete).Delete(&Order{})
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	router := mux.NewRouter()
	// membaca sesuai id
	router.HandleFunc("/orders/{orderId}", getOrder).Methods("GET")
	// membaca data keseluruhan
	router.HandleFunc("/orders", getOrders).Methods("GET")
	// membuat data
	router.HandleFunc("/orders", createOrder).Methods("POST")
	// Update
	router.HandleFunc("/orders/{orderId}", updateOrder).Methods("PUT")
	// menghapus data
	router.HandleFunc("/orders/{orderId}", deleteOrder).Methods("DELETE")
	initDB()
	log.Fatal(http.ListenAndServe(":8080", router))
}

