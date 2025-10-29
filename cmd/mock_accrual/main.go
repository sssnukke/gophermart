package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type AccrualResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

var (
	mu     sync.Mutex
	orders = make(map[string]*AccrualResponse)
)

func writeJSON(w http.ResponseWriter, data *AccrualResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func advanceStatus(o *AccrualResponse) {
	switch o.Status {
	case "REGISTERED":
		if rand.Intn(3) == 0 {
			o.Status = "PROCESSING"
		}
	case "PROCESSING":
		if rand.Intn(3) == 0 {
			if rand.Float32() < 0.8 {
				o.Status = "PROCESSED"
				accrual := float64(rand.Intn(500)) + rand.Float64()*100
				o.Accrual = &accrual
			} else {
				o.Status = "INVALID"
			}
		}
	case "PROCESSED", "INVALID":
	}
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Path[len("/api/orders/"):]

	if orderID == "" {
		http.Error(w, "order number required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if o, ok := orders[orderID]; ok {
		advanceStatus(o)
		writeJSON(w, o)
		return
	}

	status := "REGISTERED"
	orders[orderID] = &AccrualResponse{
		Order:  orderID,
		Status: status,
	}

	writeJSON(w, orders[orderID])
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/api/orders/", handleOrder)

	log.Println("🚀 Mock Accrual System запущен на :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
