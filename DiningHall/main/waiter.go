package main

import (
	"strconv"
	"time"
)

var waiterStatus = [...]string{"Waiting.", "Sending order id:", "Delivering delivery id:"}

type Waiter struct {
	id         int
	atWork     int
	statusId   int
	modifierId int
}

func NewWaiter(id int, atWork int, statusId int, modifierId int) *Waiter {
	ret := new(Waiter)
	ret.id = id
	ret.atWork = atWork
	ret.statusId = statusId
	ret.modifierId = modifierId
	return ret
}

func (w *Waiter) startWorking() {
	w.atWork = 1
	for w.atWork == 1 {
		didATask := false

		//Look for table that needs their order taken
		table := diningHall.tableList.lookUpTable()

		if table != nil {
			//Get order
			order := table.getOrder(w)

			//Send order
			success := diningHall.sendOrder(order)
			if success {
				didATask = true
				w.modifierId = order.Id
				w.statusId = 1
				time.Sleep(time.Second)
			} else {
				go table.waitForOrderList()
				didATask = false
			}
		}

		//Receive delivery
		delivery := diningHall.diningHallWeb.getDelivery()
		if delivery != nil {
			didATask = true
			//Serve delivery to the required table
			w.statusId = 2
			w.modifierId = delivery.OrderId
			time.Sleep(time.Second)
			go diningHall.tableList.tableList[delivery.TableId].deliver(delivery)
		}

		if !didATask {
			//Wait one second because there are no tasks
			w.statusId = 0
			time.Sleep(time.Second)
		}

	}
}
func (w *Waiter) getStatus() string {
	status := "Waiter id:" + strconv.Itoa(w.id) + " Status:" + waiterStatus[w.statusId]
	if w.statusId != 0 {
		return status + strconv.Itoa(w.modifierId)
	}
	return status
}
