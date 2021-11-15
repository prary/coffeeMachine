package main

import (
	"fmt"
	"strconv"
	"sync"
)

type Outlet struct {
	Count int `json:"count_n"`
}

type order struct {
	name string
	item map[string]int
}

type CoffeeMachine struct {
	OutletNo  Outlet                 `json:"outlets"`
	Inventory map[string]interface{} `json:"total_items_quantity"'`
	Beverages map[string]interface{} `json:"beverages"`

	// private variable which are initialised after parsing is done
	order     []order
	inventory map[string]int
	// lock to prevent race condition
	lock sync.RWMutex
}

func (c *CoffeeMachine) getValue(key string) int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	temp := c.inventory[key]
	return temp
}

func (c *CoffeeMachine) setValue(key string, val int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.inventory[key] -= val
}

// function to initialise inventory
func (c *CoffeeMachine) setInventory() {
	c.inventory = make(map[string]int)
	for item, value := range c.Inventory {
		x := fmt.Sprint(value)
		y, _ := strconv.Atoi(x)
		c.inventory[item] = y
	}
	//fmt.Println(c.inventory)
}

// api to refill ingredient which is low
func (c *CoffeeMachine) refill(ingredient string, value int) {
	if _, ok := c.inventory[ingredient]; ok == false {
		fmt.Println("Invalid ingredient")
	}
	c.inventory[ingredient] += value
}

// get the snapshot of inventory
func (c *CoffeeMachine) getInventory() {
	for key, value := range c.inventory {
		fmt.Println(key, value)
	}
}

func (c *CoffeeMachine) setBeverageOrder() {
	for key, value := range c.Beverages {
		//fmt.Println(key,value)
		r := order{name: key}
		itemMap := make(map[string]int)
		for i, r := range value.(map[string]interface{}) {
			quantity := fmt.Sprint(r)
			q, _ := strconv.Atoi(quantity)
			//fmt.Println(i,q)
			itemMap[i] = q
		}
		r.item = itemMap
		c.order = append(c.order, r)
		//fmt.Println("Beverage Order", r.name,"successfully placed")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// function that parallelly make beverage
func (c *CoffeeMachine) makeBeverage() {
	wg := new(sync.WaitGroup)
	c.lock = sync.RWMutex{}

	chnl := make(chan int, c.OutletNo.Count)
	for outlet := 0; outlet < len(c.order); outlet++ {
		wg.Add(1)
		go c.cookBeverage(c.order[outlet], wg, chnl)
	}
	wg.Wait()
}

func (c *CoffeeMachine) cookBeverage(order order, wg *sync.WaitGroup, chnl chan int) {
	defer wg.Done()
	chnl <- 1
	key := []string{}
	for item, _ := range order.item {
		key = append(key, item)
		if _, ok := c.inventory[item]; ok == false {
			fmt.Printf("%s cannot be prepare because %s is not available\n", order.name, item)
			<-chnl
			return
		}
	}
	for idx := 0; idx < len(key); idx++ {
		item := key[idx]
		value := order.item[key[idx]]
		if value > c.getValue(key[idx]) {
			fmt.Printf("%s cannot be prepare because %s is not sufficient\n", order.name, item)
			// Rollback all the changes
			for idx > 0 {
				c.setValue(item, -value)
				idx--
			}
			<-chnl
			return
		}
		c.setValue(item, value)
	}
	fmt.Printf("%s is prepared\n", order.name)
	<-chnl
}
