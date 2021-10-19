package main

import (
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

)

var tableStatus = [...]string{"Waiting for delivery.", "Waiting for customers.", "Eating.", "Waiting for kitchen's order list to empty.", "Waiting for waiter."}

type Table struct {
	id        int
	ordered   int32
	occupied  int32
	available int32
	statusId  int
	order     *Order
}

func NewTable(id int, ordered int32, occupied int32, available int32, statusId int, order *Order) *Table {
	ret := new(Table)
	ret.id = id
	ret.ordered = ordered
	ret.occupied = occupied
	ret.available = available
	ret.statusId = statusId
	ret.order = order
	return ret
}

func (t *Table) startAvailability() {
	atomic.StoreInt32(&t.available, 1)
	t.waitCustomers()
}

func (t *Table) deliver(delivery *Delivery) {
	//Wait based on the delivery size
	t.statusId = 2
	time.Sleep(time.Second * time.Duration(len(delivery.Items)))


	atomic.StoreInt32(&t.ordered, 0)
	atomic.StoreInt32(&t.occupied, 0)
	t.waitCustomers()
}

func (t *Table) waitCustomers() {
	if t.available == 1 && t.occupied == 0 {
		syncMutex := sync.Mutex{}
		atomic.StoreInt32(&t.ordered, 1)
		t.statusId = 1
		time.Sleep(time.Second * time.Duration(rand.Intn(10)))

		syncMutex.Lock()
		if t.order == nil {
			t.order = generateOrder(t)
		}
		syncMutex.Unlock()
		//addr := (*unsafe.Pointer)(unsafe.Pointer(t.order))
		//newOrder := unsafe.Pointer(generateOrder(t))
		//atomic.StorePointer(addr, newOrder)
		atomic.StoreInt32(&t.ordered, 0)
		atomic.StoreInt32(&t.occupied, 1)
		t.statusId = 4
	}
}

func (t *Table) getOrder(waiter *Waiter) *Order {
	t.statusId = 0
	t.order.WaiterId = waiter.id
	return t.order
}

func (t *Table) stopAvailability() {
	atomic.StoreInt32(&t.available, 0)
}

func (t *Table) waitForOrderList() {
	atomic.StoreInt32(&t.ordered, 1)
	t.statusId = 3
	time.Sleep(time.Second * 2) //Wait 2 seconds for the order list to free
	atomic.StoreInt32(&t.ordered, 0)
}

func (t *Table) getStatus() string {
	waitStatus := ""
	if t.occupied == 1 && t.ordered == 1 && t.order != nil {
		waitStatus = " Waiting for:" + strconv.Itoa(int(time.Now().Unix()-t.order.PickUpTime)) + "sec" + " Max wait:" + strconv.Itoa(t.order.MaxWait)
	}
	return "Table id:" + strconv.Itoa(t.id) + " Status:" + tableStatus[t.statusId] + waitStatus
}
