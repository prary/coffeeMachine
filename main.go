package main

import (
	"encoding/json"
	"io/ioutil"
)

type Machine struct {
	CoffeeMachine CoffeeMachine `json:"machine"`
}

func main(){
	// Reading input from input.json file
	file, _ := ioutil.ReadFile("input.json")
	var machine Machine
	// Parsing the inpyt json
	_ = json.Unmarshal([]byte(file),&machine)

	var coffeeMachine CoffeeMachine
	coffeeMachine = machine.CoffeeMachine
	// Intialising inventory
	coffeeMachine.setInventory()
	// Intialising all the orders
	coffeeMachine.setBeverageOrder()
	// making beverage in parrallel
	coffeeMachine.makeBeverage()
}
