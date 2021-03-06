package main

import (
	"fmt"
	"time"
	"os/exec"
	"../network/network/bcast"
)

const udpPort = 20002
//const myIP = "129.241.187.255"


type Counter struct{
	State int
}

type Message struct{
	Data int
}

func Backup(counter Counter){

	toBackup := make(chan Message, 1)

	spawnBackup := exec.Command("gnome-terminal", "-x", "sh", "-c", "go run backup.go")
	
	spawnBackup.Start()
	
	go bcast.Transmitter(udpPort, toBackup)
	
	
	for {
		fmt.Printf("Counter: %d \n", counter.State)
		msg := Message{counter.State}
		toBackup <- msg
		counter.State++
		time.Sleep(1*time.Second)
	}
	
}

func main(){
	var isBackup bool = true
	timer := time.NewTimer(3*time.Second)
	fmt.Print("Hello, I'm backup\n")
	fromMaster := make(chan Message, 1)
	
	masterCounter := Counter{0}
	go bcast.Receiver(udpPort, fromMaster)
	
	go func(){
		<- timer.C
		isBackup = false
		fmt.Print("Backup becomes the master\n")
		masterCounter.State++
		Backup(masterCounter)
	}()	

	for {
		if isBackup {
			msg := <- fromMaster
			masterCounter.State = msg.Data
			timer.Reset(3*time.Second)
		}
	}	
}	
	





