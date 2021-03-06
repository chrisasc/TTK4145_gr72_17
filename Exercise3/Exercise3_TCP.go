package main

import (
	"fmt"
	"net"
	//"time"
	"bufio"
	"os"
)

const bcast1 = "129.241.187.255:34933"
const bcast2 = "129.241.187.255:20002"


func CheckError (err error){
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main(){
	
	Conn, err := net.Dial("tcp",  bcast1)
	CheckError(err)
	
	defer Conn.Close()
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Text to send: ")
	
	text, _ := reader.ReadString('\n')
	fmt.Fprint(Conn, text+"\n")
	
	msg, _ := bufio.NewReader(Conn).ReadString('\n')
	fmt.Print("Message for server: " + msg)
	
	
	
	

}


