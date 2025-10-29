package orders

import (
	"encoding/json"
	"fmt"
	"gophermart/internal/db"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Worker struct {
	db          *gorm.DB
	accrualAddr string
	interval    time.Duration
	client      *http.Client
}

func NewWorker(db *gorm.DB, accrualAddr string, interval time.Duration) *Worker {
	return &Worker{
		db:          db,
		accrualAddr: accrualAddr,
		interval:    interval,
		client:      &http.Client{Timeout: 2 * time.Second},
	}
}

func (w *Worker) processOrders() {
	var orders []db.Order
	if err := w.db.Where("status IN ?", []string{"NEW", "PROCESSING", "REGISTERED"}).Find(&orders).Error; err != nil {
		log.Println("Error get orders: ", err)
		return
	}

	for _, order := range orders {
		url := fmt.Sprintf("%s/api/orders/%s", w.accrualAddr, order.Number)
		resp, err := w.client.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			continue
		}

		var data struct {
			Order   string   `json:"order"`
			Status  string   `json:"status"`
			Accrual *float64 `json:"accrual,omitempty"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			resp.Body.Close()
			log.Println(err)
			continue
		}

		order.Status = data.Status
		order.Accrual = data.Accrual
		w.db.Save(&order)

		if data.Status == "PROCESSED" && data.Accrual != nil {
			w.db.Model(&db.User{}).Where("id = ?", order.UserID).Update("balance", gorm.Expr("balance + ?", *data.Accrual))
			log.Printf("Add %.2f balance %d", *data.Accrual, order.UserID)
		}
	}
}
func (w *Worker) Start() {
	ticker := time.NewTicker(w.interval)
	go func() {
		for range ticker.C {
			w.processOrders()
		}
	}()
}
