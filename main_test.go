package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckout(t *testing.T) {
	myRouter = mux.NewRouter().StrictSlash(true)

	testCases := []struct {
		testID        string
		testType      string
		testDesc      string
		payload       CheckoutRequest
		expHTTPStatus int
		expTotal      float64
		expItems      []ScannedItem
	}{
		{
			testID:   "1",
			testDesc: "buy 3 google home, pay 2 ",
			testType: "P",
			payload: CheckoutRequest{
				Items: []ScannedItem{
					{
						Sku: GOOGLE_HOME,
						Qty: 3,
					},
				},
			},
			expHTTPStatus: 200,
			expTotal:      99.98,
			expItems: []ScannedItem{
				{
					Sku: GOOGLE_HOME,
					Qty: 3,
				},
			},
		},
		{
			testID:   "2",
			testDesc: "buy 1 macbook free 1 raspberry pi ",
			testType: "P",
			payload: CheckoutRequest{
				Items: []ScannedItem{
					{
						Sku: MACBOOK_PRO,
						Qty: 1,
					},
				},
			},
			expHTTPStatus: 200,
			expTotal:      5399.99,
			expItems: []ScannedItem{
				{
					Sku: MACBOOK_PRO,
					Qty: 1,
				},
				{
					Sku: RASPBERRY_PI_B,
					Qty: 1,
				},
			},
		},
		{
			testID:   "3",
			testDesc: "buy 3 alexa speaker, discount 10%",
			testType: "P",
			payload: CheckoutRequest{
				Items: []ScannedItem{
					{
						Sku: ALEXA_SPEAKER,
						Qty: 3,
					},
				},
			},
			expHTTPStatus: 200,
			expTotal:      295.65,
			expItems: []ScannedItem{
				{
					Sku: ALEXA_SPEAKER,
					Qty: 3,
				},
			},
		},
		{
			testID:   "4",
			testDesc: "buy 99 google home, out of stock",
			testType: "N",
			payload: CheckoutRequest{
				Items: []ScannedItem{
					{
						Sku: GOOGLE_HOME,
						Qty: 99,
					},
				},
			},
			expHTTPStatus: 400,
		},
	}

	Convey("Test checkout items", t, func() {
		for _, tc := range testCases {
			Convey(fmt.Sprintf("%s - [%s] : %s", tc.testID, tc.testType, tc.testDesc), func() {
				body, err := json.Marshal(tc.payload)
				if err != nil {
					t.Error("Marshal Body Req", err.Error())
				}

				req, _ := http.NewRequest("POST", "/checkout", bytes.NewBuffer(body))
				myRouter.HandleFunc("/checkout", checkoutItems).Methods("POST")
				response := executeRequest(req)

				So(response.Code, ShouldEqual, tc.expHTTPStatus)

				if tc.testType == "P" {
					// assert resp
					var result CheckoutResponse
					rawData, _ := ioutil.ReadAll(response.Body)
					if err := json.Unmarshal(rawData, &result); err != nil {
						t.Error("Unmarshal Resp ", err.Error())
					}

					So(result.Total, ShouldEqual, tc.expTotal)
					So(result.Items, ShouldResemble, tc.expItems)
				}
			})
		}
	})
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	myRouter.ServeHTTP(rr, req)
	return rr
}
