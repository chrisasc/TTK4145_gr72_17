package order_handler

import (
	"../driver"
	//. "../network"
	"../structs"
	"fmt"
	"time"
)

//for all orders in queue, sends new floor if order is command or matches direction
func other_orders_in_dir(OrderQueue []structs.Order, new_target_floor chan<- int, State structs.Elev_state) {
	var floorSignal = driver.GetFloorSignal()
	if floorSignal != -1 {
		for _, order := range OrderQueue {
			if order.Floor == floorSignal && (int(order.Type) == int(State.Current_direction)+1 || int(order.Type) == 1) {
				time.Sleep(10 * time.Millisecond)
				new_target_floor <- order.Floor
			}
		}
	}
}

func is_duplicate(order structs.Order, OrderQueue []structs.Order) bool {
	for _, order_iter := range OrderQueue {
		if order == order_iter {
			return true
		}
	}
	return false
}

func get_new_order(OrderQueue []structs.Order, new_target_floor chan<- int, localIP string) {
	for _, order := range OrderQueue {
		if order.IP == localIP {
			fmt.Printf("Sending new target floor from get_new_order\n")
			new_target_floor <- OrderQueue[0].Floor
		}
	}
}

func printOrderQueue(OrderQueue []structs.Order) {
	fmt.Printf("---------------------------------")
	fmt.Printf("PRINTING FROM FUNCTION\n")
	fmt.Printf("Order Queue: \n")
	for i, order := range OrderQueue {
		fmt.Printf("Element: %d \n", i+1)
		fmt.Print("Button  Floor\n")
		fmt.Print(order.Type, order.Floor+1)
		fmt.Printf("\n")
	}
	fmt.Printf("Length: ")
	fmt.Print(len(OrderQueue))
	fmt.Printf("\n")
	fmt.Printf("---------------------------------")
}

func add_order(order structs.Order, OrderQueue []structs.Order, new_target_floor chan<- int, localIP string) []structs.Order {
	if is_duplicate(order, OrderQueue) == false {
		OrderQueue = append(OrderQueue, order)
		printOrderQueue(OrderQueue)
		if len(OrderQueue) == 1 {
			get_new_order(OrderQueue, new_target_floor, localIP)
		}
		//driver.SetButtonLamp(order.Type, order.Floor, 1)
	}
	return OrderQueue
}

func remove_order(order structs.Order, OrderQueue []structs.Order) []structs.Order {
	for index, order_iter := range OrderQueue {
		if order_iter == order {
			OrderQueue = append(OrderQueue[:index], OrderQueue[index+1:]...)
			//driver.SetButtonLamp(OrderQueue[index].Type, OrderQueue[index].Floor, 0)
			return OrderQueue
		}
	}
	return nil
}

func remove_all(floor int, elev_send_remove_order chan<- structs.Order, OrderQueue []structs.Order) []structs.Order {
	for _, order := range OrderQueue {
		if order.Floor == floor {
			fmt.Printf("Found order in floor, removing\n")
			OrderQueue = remove_order(order, OrderQueue)
			//elev_send_remove_order <- order
		}
	}
	return OrderQueue
}

/*
func merge_orders(){
	//Ved nettverksbrud, merger ordre med våre
	//Kanskje ikke nødvendig da alle ekstern bestillingen allerede er i køen vår
}*/

func Order_handler_init(State structs.Elev_state,
	OrderQueue []structs.Order,
	localIP string,
	floor_completed <-chan int,
	button_event <-chan driver.OrderButton,
	assignedNewOrder <-chan structs.Order,
	newOrder chan<- structs.Order,
	elev_send_new_order chan<- structs.Order,
	elev_send_remove_order chan<- structs.Order,
	elev_receive_new_order <-chan structs.Order,
	elev_receive_remove_order <-chan structs.Order,
	new_target_floor chan<- int) {

	for {
		select {
		case floor := <-floor_completed:
			fmt.Printf("Floor completed message received\n")
			OrderQueue = remove_all(floor, elev_send_remove_order, OrderQueue)
			fmt.Printf("Removed from order queue\n")
			if len(OrderQueue) != 0 {
				printOrderQueue(OrderQueue)
				fmt.Printf("Retrieving new order\n")
				get_new_order(OrderQueue, new_target_floor, localIP)
			}

		case order_button := <-button_event:
			if order_button.Type == driver.ButtonCallCommand {
				new_order := structs.Order{Type: order_button.Type, Floor: order_button.Floor, Internal: true, IP: localIP}
				OrderQueue = add_order(new_order, OrderQueue, new_target_floor, localIP)

			} else { // if external, send to order_distribution
				new_order := structs.Order{Type: order_button.Type, Floor: order_button.Floor, Internal: false, IP: localIP}
				newOrder <- new_order
				elev_send_new_order <- new_order // for å sende til network
			}
		case new_order := <-assignedNewOrder:
			OrderQueue = add_order(new_order, OrderQueue, new_target_floor, localIP)

		case order := <-elev_receive_remove_order:
			OrderQueue = remove_order(order, OrderQueue)

		case new_order := <-elev_receive_new_order:
			OrderQueue = add_order(new_order, OrderQueue, new_target_floor, localIP)
		default:
			other_orders_in_dir(OrderQueue, new_target_floor, State)
		}

	}
}
