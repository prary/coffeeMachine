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
	OutletNo Outlet						`json:"outlets"`
	Inventory map[string]interface{}	`json:"total_items_quantity"'`
	Beverages map[string]interface{}    `json:"beverages"`

	// private variable which are initialised after parsing is done
	order	  []order
	inventory map[string]int
	// lock to prevent race condition
	sync.RWMutex
}

// function to initialise inventory
func(c *CoffeeMachine) setInventory()  {
	c.inventory = make(map[string]int)
	for item,value :=range c.Inventory {
		x := fmt.Sprint(value)
		y, _ :=strconv.Atoi(x)
		c.inventory[item] = y
	}
	//fmt.Println(c.inventory)
}
// api to refill ingredient which is low
func(c *CoffeeMachine) refill(ingredient string, value int) {
	if _,ok := c.inventory[ingredient]; ok==false {
		fmt.Println("Invalid ingredient")
	}
	c.inventory[ingredient] += value
}
// get the snapshot of inventory
func(c *CoffeeMachine) getInventory()  {
	for key,value := range c.inventory {
		fmt.Println(key,value)
	}
}

func(c *CoffeeMachine) setBeverageOrder() {
	for key,value := range c.Beverages {
		//fmt.Println(key,value)
		r := order{name: key}
		itemMap := make(map[string]int)
		for i,r := range value.(map[string]interface{}){
			quantity := fmt.Sprint(r)
			q, _ := strconv.Atoi(quantity)
			//fmt.Println(i,q)
			itemMap[i] = q
		}
		r.item = itemMap
		c.order = append(c.order,r)
		//fmt.Println("Beverage Order", r.name,"successfully placed")
	}
}

func min(a,b int) int {
	if a<b {
		return a
	}
	return b
}
// function that parallelly make beverage
func (c *CoffeeMachine) makeBeverage() {
	wg := new(sync.WaitGroup)
	//fmt.Println(min(c.OutletNo.Count,len(c.order)))
	wg.Add(min(c.OutletNo.Count,len(c.order)))
	for outlet:=0; outlet<c.OutletNo.Count && outlet<len(c.order); outlet++ {
		// spins up new goroutine for every order and with empty outlet
		c.Lock()
		go c.cookBeverage(c.order[outlet], wg)
		c.Unlock()
	}
	wg.Wait()
	//time.Sleep(5*time.Second)
}

func (c *CoffeeMachine) cookBeverage(order order, wg *sync.WaitGroup){
	defer wg.Done()
	mapMutex := sync.Mutex{}
	//fmt.Printf("Cooking %s\n",order.name)
	for item,value := range order.item {
		mapMutex.Lock()
		if _,ok := c.inventory[item]; ok == false {
			fmt.Printf("%s cannot be prepare because %s is not available\n", order.name,item)
			mapMutex.Unlock()
			return
		}else{
			if value > c.inventory[item] {
				mapMutex.Unlock()
				//fmt.Println(order.name, item, value,c.inventory[item])
				fmt.Printf("%s cannot be prepare because %s is not sufficient\n", order.name,item)
				return
			}
			//mapMutex.RUnlock()
			//mapMutex.Lock()
			c.inventory[item] = c.inventory[item] - value
			mapMutex.Unlock()
			//c.Unlock()
		}
	}
	fmt.Printf("%s is prepared\n", order.name)
}
