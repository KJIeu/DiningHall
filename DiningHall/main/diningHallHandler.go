package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DiningHallHandler struct {
	packetsReceived int32
	postReceived    int32
}

func (d DiningHallHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		{
			response := "OK"
			latestDelivery := diningHall.diningHallWeb.getDelivery()
			if latestDelivery != nil {
				diningHall.diningHallWeb.setDelivery(latestDelivery)
				response = "NOT OK"
				fmt.Fprint(w, response)
				return
			} else {
				latestDelivery = new(Delivery)
				var requestBody = make([]byte, r.ContentLength)
				r.Body.Read(requestBody)
				json.Unmarshal(requestBody, latestDelivery)
				diningHall.diningHallWeb.setDelivery(latestDelivery)

				//Respond with "OK"
				fmt.Fprint(w, response)
			}
		}
	case http.MethodGet:
		{
			fmt.Fprintln(w, "<head><meta http-equiv=\"refresh\" content=\"1\" /></head>")
			if diningHall.connected {
				fmt.Fprintln(w, makeDiv(""))
			} else {
				fmt.Fprintln(w, makeDiv("DiningHall did not establish connection to kitchen on address:"+kitchenServerHost+kitchenServerPort))
				err := diningHall.diningHallWeb.connectionError
				if err != nil {
					fmt.Fprintln(w, makeDiv("Connection error: "+err.Error()))
				}
			}
			fmt.Fprintln(w, makeDiv(diningHall.getStatus()))
		}
	case http.MethodConnect:
		{
			diningHall.connectionSuccessful()
			fmt.Fprint(w, "OK")
		}
	default:
		{
			fmt.Fprintln(w, "UNSUPPORTED METHOD")
		}
	}
}
