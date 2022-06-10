package main

import (
	conFiles "control/files_controller"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	CH   = "channel"
	SUBS = "subscribe"

	SND = "send"

	CLT = "client "
)

var counterClient = 1
var status = false

type clientStruct struct {
	name    string
	add     string
	conn    net.Conn
	channel int
}

func (this *clientStruct) detailClient(name1 string, conn1 net.Conn, address1 string) {
	this.name = name1
	this.conn = conn1
	this.add = address1
}

var clients []clientStruct

func main() {

	//start server
	server, err := net.Listen("tcp", "127.0.0.1:8000")

	if err != nil {
		panic(err)
	}
	fmt.Println("./server start")
	fmt.Println()

	for {
		//a client connects

		client, err := server.Accept()
		if err != nil {
			fmt.Printf("An error has occurred: %v \n", err)
			return
		}
		client1 := new(clientStruct)

		if err != nil {
			fmt.Printf("An error has occurred: %v \n", err)
			return
		}

		if status {
			currentClient := CLT + strconv.Itoa(counterClient)
			client1.detailClient(currentClient, client, client.RemoteAddr().String())
			clients = append(clients, *client1)
			counterClient++
		} else {
			currentClient := CLT + strconv.Itoa(counterClient)
			client1.detailClient(currentClient, client, client.RemoteAddr().String())
			clients = append(clients, *client1)
			status = true
			counterClient++
		}

		go listenClients(client)

	}

}

func listenClients(client net.Conn) {
	for {
		var buff = make([]byte, 2048)
		n, err := client.Read(buff)

		if err != nil {

			return
		}

		var text = string(buff[:n])

		optionServer(client, text)

	}

}

func optionServer(client net.Conn, text string) {
	dirClientCurrent := client.RemoteAddr().String()

	args := strings.Split(text, " ")

	switch args[0] {
	case CH:
		var canal int
		canal, _ = strconv.Atoi(args[1])
		for i, cs := range clients {
			if cs.conn.RemoteAddr().String() == dirClientCurrent {
				clients[i].channel = canal
				//texto en consola client subscrito
				fmt.Printf("./%s subscribed -channel %d\n\n", clients[i].name, clients[i].channel)

				break
			}
		}

	case SND:
		nameFile := args[1]
		channel, _ := strconv.Atoi(args[2])
		channelNum := args[2]

		var clientSend string
		for i, cs := range clients {
			if dirClientCurrent == cs.conn.RemoteAddr().String() {
				//fmt.Printf("./%s send %s -channel %d\n", clients[i].nombre, nameFile, canal)
				clientSend = "./" + clients[i].name + " send " + nameFile + "-channel " + channelNum
			}

			if dirClientCurrent != cs.conn.RemoteAddr().String() && cs.channel == channel {
				fmt.Printf("./%s receive -channel %d\n\n", clients[i].name, channel)

			}
		}
		fmt.Println(clientSend)
		fmt.Println()

		nameFile = "./" + nameFile

		conFiles.ReceiveFiles(client, nameFile)
		writeFilesAllUsers(nameFile, client, channel)

	default:
		fmt.Printf("wrong entry: %s\n", text)
	}

}

func writeFilesAllUsers(path string, conn net.Conn, channel int) {

	for i := range clients {
		if (clients[i].conn.RemoteAddr().String() != conn.RemoteAddr().String()) && clients[i].channel == channel {
			clients[i].conn.Write([]byte(path))
			conFiles.SendFiles(clients[i].conn, path)
		}

	}

}
