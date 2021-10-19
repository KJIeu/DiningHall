package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type DiningHallWeb struct {
	diningHallServer  http.Server
	diningHallHandler DiningHallHandler
	diningHallClient  http.Client
	latestDelivery    *Delivery
	connectionError   error
	deliveryMutex     sync.Mutex
}

func (dhw *DiningHallWeb) getDelivery() *Delivery {
	dhw.deliveryMutex.Lock()
	defer dhw.deliveryMutex.Unlock()

	ret := dhw.latestDelivery
	dhw.latestDelivery = getNilDelivery()

	return ret
}

func (dhw *DiningHallWeb) setDelivery(delivery *Delivery) {
	dhw.deliveryMutex.Lock()
	defer dhw.deliveryMutex.Unlock()

	dhw.latestDelivery = delivery

}

func (dhw *DiningHallWeb) start() {
	dhw.diningHallServer.Addr = diningHallPort
	dhw.diningHallServer.Handler = &dhw.diningHallHandler

	fmt.Println(time.Now())
	fmt.Println("DiningHallServer is listening and serving on port:" + diningHallPort)
	if err := dhw.diningHallServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (dhw *DiningHallWeb) sendOrder(order *Order) bool {
	requestBody := order.getPayload()
	request, _ := http.NewRequest(http.MethodPost, kitchenServerHost+kitchenServerPort+"/order", bytes.NewBuffer(requestBody))
	response, err := dhw.diningHallClient.Do(request)

	if err != nil {
		fmt.Println("Could not send order to kitchen.")
		log.Fatal(err)
		return false
	}
	var responseBody = make([]byte, response.ContentLength)
	response.Body.Read(responseBody)
	if string(responseBody) != "OK" {
		return false
	}

	return true
}

func (dhw *DiningHallWeb) establishConnection() bool {
	if diningHall.connected == true {
		return false
	}
	request, _ := http.NewRequest(http.MethodConnect, kitchenServerHost+kitchenServerPort+"/", bytes.NewBuffer([]byte{}))
	response, err := dhw.diningHallClient.Do(request)
	if err != nil {
		dhw.connectionError = err
		return false
	}
	var responseBody = make([]byte, response.ContentLength)
	response.Body.Read(responseBody)
	if string(responseBody) != "OK" {
		return false
	}

	return true
}
