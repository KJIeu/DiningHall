package main

import "sync"

type TableList struct {
	tableList        []*Table
	tableListCounter int
	tableMutex       sync.Mutex
}

func NewTableList() *TableList {
	tableListCounter := 0
	tableList := make([]*Table, 0)
	for i := 0; i < tableN; i++ {
		tableList = append(tableList, NewTable(tableListCounter, 0, 0, 0, 1, nil))
		tableListCounter++
	}
	return &TableList{tableList, tableListCounter, sync.Mutex{}}
}

func (tl *TableList) start() {
	for _, table := range tl.tableList {
		go table.startAvailability()
	}
}

//func (tl *TableList) lookUpForOrder() *Order {
//	for _, table := range tl.tableList {
//		if table.available == 1 && table.occupied == 1 && table.ordered == 0 {
//			return table.order
//		}
//	}
//	return nil
//}

func (tl *TableList) lookUpTable() *Table {
	tl.tableMutex.Lock()
	defer tl.tableMutex.Unlock()

	for _, table := range tl.tableList {
		if table.available == 1 && table.occupied == 1 && table.ordered == 0 {
			table.ordered = 1
			return table
		}
	}
	return nil
}
