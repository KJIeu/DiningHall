package main

var Dishes = []Menu{Pizza, Salad, Zeama, ScallopSashimi, IslandDuck, Waffles, Aubergine, Lasagna, Burger, Gyros}
var Machines = map[string]int{"": 0, "oven": 1, "stove": 2}
var MachinesId = map[int]string{0: "nil", 1: "oven", 2: "stove"}

type Menu struct {
	Id              	int
	Name            	string
	PreparationTime  	int
	Complexity       	int
	CookingApparatus 	string
}

var Pizza = Menu{1, "Pizza", 20, 2, "oven"}
var Salad = Menu{2, "Salad", 10, 1, ""}
var Zeama = Menu{3, "Zeama", 7, 1, "stove"}
var ScallopSashimi = Menu{4, "Scallop Sashimi with Meyer Lemon Confit", 32, 3, ""}
var IslandDuck = Menu{5, "Island Duck with Mulberry Mustard", 35, 3, "oven"}
var Waffles = Menu{6, "Waffles", 10, 1, "stove"}
var Aubergine = Menu{7, "Aubergine", 20, 2, ""}
var Lasagna = Menu{8, "Lasagna", 30, 2, "oven"}
var Burger = Menu{9, "Burger", 15, 1, "oven"}
var Gyros = Menu{10, "Gyros", 15, 1, ""}