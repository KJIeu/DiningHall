package main

type WaiterList struct {
	waiterList      []*Waiter
	waiterIdCounter int
}

func NewWaiterList() *WaiterList {
	waiterList := make([]*Waiter, 0)
	waiterIdCounter := 0
	for i := 0; i < waiterN; i++ {
		waiterList = append(waiterList, NewWaiter(waiterIdCounter, 0, 0, 0))
		waiterIdCounter++
	}

	return &WaiterList{waiterList, waiterIdCounter}
}

func (wl WaiterList) start() {
	for _, waiter := range wl.waiterList {
		go waiter.startWorking()
	}
}
