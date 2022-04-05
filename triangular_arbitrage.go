package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type order_book struct {
	Timestamp      string     `json:'timestamp'`
	Microtimestamp string     `json:'microtimestamp'`
	Bids           [][]string `json:'bids'`
	Asks           [][]string `json:'asks'`
}

func main() {

	currency1 := "xrpeur"
	currency2 := "eurusd"
	currency3 := "xrpusd"
	url := "https://www.bitstamp.net/api/v2/order_book/"

	var order_book_xrpeur order_book
	var order_book_eurusd order_book
	var order_book_xrpusd order_book

	var counter int = 0
	var opportunity int = 0
	var profit float64 = 0

	for {
		//http.get calls in func get_order_book
		order_book_xrpeur = *get_order_book(url + currency1 + "/")
		order_book_eurusd = *get_order_book(url + currency2 + "/")
		order_book_xrpusd = *get_order_book(url + currency3 + "/")
		//string to float64 conversions
		var xrpeurAsk float64 = stringToFloat64(order_book_xrpeur.Asks[0][0])
		var xrpusdBid float64 = stringToFloat64(order_book_xrpusd.Bids[0][0])
		var eurusdAsk float64 = stringToFloat64(order_book_eurusd.Asks[0][0])
		var eurusdBid float64 = stringToFloat64(order_book_eurusd.Bids[0][0])
		var xrpusdAsk float64 = stringToFloat64(order_book_xrpusd.Asks[0][0])
		var xrpeurBid float64 = stringToFloat64(order_book_xrpeur.Bids[0][0])

		//calculating path1 : EUR -> XRP -> USD -> EUR
		var eur_xrp_usd_eur float64 = ((1 / xrpeurBid * 0.9995) * xrpusdAsk * 0.9995) / eurusdBid * 0.9995
		//calculating path2 : EUR -> USD -> XRP -> EUR
		var eur_usd_xrp_eur float64 = ((1 * eurusdAsk * 0.9995) / xrpusdBid * 0.9995) * xrpeurAsk * 0.9995

		if eur_xrp_usd_eur > 1 || eur_usd_xrp_eur > 1 {
			opportunity++
			//path1 trading quantity and profit calculation
			if eur_xrp_usd_eur > 1 {
				var eurusdBidQuantity float64 = stringToFloat64(order_book_eurusd.Bids[0][1])
				var xrpusdAskQuantity float64 = stringToFloat64(order_book_xrpusd.Asks[0][1])
				var xrpeurBidQuantity float64 = stringToFloat64(order_book_xrpeur.Bids[0][1])
				//minQuantity is in EUR
				var minQuantity float64 = eurusdBidQuantity
				//amount in xrp converted in EUR so it can be compared
				if minQuantity > xrpusdAskQuantity*xrpeurBid {
					minQuantity = xrpusdAskQuantity * xrpeurBid
				}
				//xrpeurBid conversion rate is used in path1
				if minQuantity > xrpeurBidQuantity*xrpeurBid {
					minQuantity = xrpeurBidQuantity * xrpeurBid
				}
				profit += minQuantity * (eur_xrp_usd_eur - 1)
			}
			//path2 trading quantity and profit calculation
			if eur_usd_xrp_eur > 1 {
				var eurusdAskQuantity float64 = stringToFloat64(order_book_eurusd.Asks[0][1])
				var xrpusdBidQuantity float64 = stringToFloat64(order_book_xrpusd.Bids[0][1])
				var xrpeurAskQuantity float64 = stringToFloat64(order_book_xrpeur.Asks[0][1])
				//minQuantity is in EUR
				var minQuantity float64 = eurusdAskQuantity
				//amount in xrp converted in EUR so it can be compared
				if minQuantity > xrpeurAskQuantity*xrpeurAsk {
					minQuantity = xrpeurAskQuantity * xrpeurAsk
				}
				//xrpeurAsk conversion rate is used in path2
				if minQuantity > xrpusdBidQuantity*xrpeurAsk {
					minQuantity = xrpusdBidQuantity * xrpeurAsk
				}
				profit += minQuantity * (eur_usd_xrp_eur - 1)
			}
		}

		counter++
		if counter%600 == 0 {
			fmt.Println("opportunity frequency:", float64(opportunity)/float64(counter), "profit in eur:", profit)
		}

		time.Sleep(880 * time.Millisecond)
	}
}

func get_order_book(url string) *order_book {

	resp, getErr := http.Get(url)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var order_book_currency order_book
	jsonErr := json.Unmarshal(body, &order_book_currency)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return &order_book_currency
}

func stringToFloat64(s string) float64 {

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
