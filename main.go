package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/mux"
)

type Item struct {
	Sku   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   float64 `json:"qty"`
}

type CheckoutRequest struct {
	Items []ScannedItem `json:"items"`
}

type CheckoutResponse struct {
	Total float64       `json:"total"`
	Items []ScannedItem `json:"items"`
}

type ScannedItem struct {
	Sku string  `json:"sku"`
	Qty float64 `json:"qty"`
}

const (
	GOOGLE_HOME    = `120P90`
	MACBOOK_PRO    = `43N23P`
	ALEXA_SPEAKER  = `A304SD`
	RASPBERRY_PI_B = `234234`
)

var (
	myRouter *mux.Router

	Items = map[string]Item{
		GOOGLE_HOME: {
			Sku:   GOOGLE_HOME,
			Name:  "Google Home",
			Price: 49.99,
			Qty:   10,
		},
		MACBOOK_PRO: {
			Sku:   MACBOOK_PRO,
			Name:  "Macbook Pro",
			Price: 5399.99,
			Qty:   5,
		},
		ALEXA_SPEAKER: {
			Sku:   ALEXA_SPEAKER,
			Name:  "Alexa Speaker",
			Price: 109.50,
			Qty:   10,
		},
		RASPBERRY_PI_B: {
			Sku:   RASPBERRY_PI_B,
			Name:  "Raspberry Pi B",
			Price: 30.00,
			Qty:   2,
		},
	}
)

func checkoutItems(w http.ResponseWriter, r *http.Request) {
	// read request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var payload CheckoutRequest
	json.Unmarshal(reqBody, &payload)

	// validate sku
	for _, item := range payload.Items {
		if item.Qty < 1 {
			fmt.Printf("SKU %v qty at least 1", item.Sku)
			w.WriteHeader(400)
			return
		}

		if val, ok := Items[item.Sku]; !ok || val.Qty < item.Qty {
			fmt.Printf("SKU %v not exist or out of stock", item.Sku)
			w.WriteHeader(400)
			return
		}
	}

	result := countPrice(payload.Items)

	json.NewEncoder(w).Encode(result)
}

func countPrice(items []ScannedItem) CheckoutResponse {
	result := CheckoutResponse{}

	for _, item := range items {
		result.Total += Items[item.Sku].Price * item.Qty

		result.Items = append(result.Items, ScannedItem{
			Qty: item.Qty,
			Sku: item.Sku,
		})

		if item.Sku == MACBOOK_PRO {
			result.Items = append(result.Items, ScannedItem{
				Qty: item.Qty,
				Sku: RASPBERRY_PI_B,
			})
		} else if item.Sku == ALEXA_SPEAKER && item.Qty >= 3 {
			result.Total -= Items[item.Sku].Price * 0.1 * item.Qty
		} else if item.Sku == GOOGLE_HOME && int(item.Qty)%3 == 0 {
			result.Total -= item.Qty / 3 * Items[item.Sku].Price
		}
	}

	result.Total = math.Round(result.Total*100) / 100

	return result
}

func main() {
	myRouter = mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/checkout", checkoutItems).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
