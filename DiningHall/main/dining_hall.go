package main

import (
	"time"
)

type DiningHall struct {
	diningHallWeb DiningHallWeb
	waiterList    *WaiterList
	tableList     *TableList
	connected     bool
}

func (dh *DiningHall) start() {
	dh.waiterList = NewWaiterList()
	dh.tableList = NewTableList()
	go dh.tryConnectKitchen()
	dh.diningHallWeb.start()
}

func (dh *DiningHall) connectionSuccessful() {
	if dh.connected {
		return
	}
	dh.connected = true
	dh.tableList.start()
	dh.waiterList.start()
}

func (dh *DiningHall) tryConnectKitchen() {
	dh.connected = false
	for !dh.connected {
		if dh.diningHallWeb.establishConnection() {
			dh.connectionSuccessful()
			break
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (dh *DiningHall) sendOrder(order *Order) bool {
	return dh.diningHallWeb.sendOrder(order)
}

func (dh *DiningHall) getStatus() string {
	ret := "Waiters:"
	for _, waiter := range dh.waiterList.waiterList {
		ret += makeDiv(waiter.getStatus())
	}
	ret += "Tables:"
	for _, table := range dh.tableList.tableList {
		ret += makeDiv(table.getStatus())
	}

	return ret
}
