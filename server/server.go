package main

import (
	mA "control/manejo_archivos"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	CH   = "channel"
	SUBS = "subscribe"

	SND = "send"

	CLIENTE = "client "
)

var contadorCliente = 1
var bandera = false

type clienteStruct struct {
	nombre    string
	direccion string
	conn      net.Conn
	channel   int
}

func (this *clienteStruct) DetalleCliente(nombre string, conn net.Conn, direccion string) {
	this.nombre = nombre
	this.conn = conn
	this.direccion = direccion
}

var clients []clienteStruct

func main() {

	//inicia server
	server, err := net.Listen("tcp", "127.0.0.1:8000")

	if err != nil {
		panic(err)
	}
	fmt.Println("./server start")
	fmt.Println()

	for {
		//se conecta un cliente

		client, err := server.Accept()
		if err != nil {
			fmt.Printf("Se ha producido un error: %v \n", err)
			return
		}
		cliente1 := new(clienteStruct)

		if err != nil {
			fmt.Printf("Se ha producido un error: %v \n", err)
			return
		}

		if bandera {
			actualClient := CLIENTE + strconv.Itoa(contadorCliente)
			cliente1.DetalleCliente(actualClient, client, client.RemoteAddr().String())
			clients = append(clients, *cliente1)
			contadorCliente++
		} else {
			actualClient := CLIENTE + strconv.Itoa(contadorCliente)
			cliente1.DetalleCliente(actualClient, client, client.RemoteAddr().String())
			clients = append(clients, *cliente1)
			bandera = true
			contadorCliente++
		}

		go listenClients(client)

	}

	//mA.PreparEnvio(client, "./algo.txt")

	//go managerConnection(client)

}

func listenClients(client net.Conn) {
	for {
		var buff = make([]byte, 2048)
		n, err := client.Read(buff)
		dirClientAc := client.RemoteAddr().String()

		if err != nil {

			return
		}

		//se ejecuta una vez cada vez que se hace enter
		//var texto = scanner.Text()
		var texto = string(buff[:n])
		args := strings.Split(texto, " ")

		switch args[0] {
		case CH:
			var canal int
			canal, _ = strconv.Atoi(args[1])
			for i, cs := range clients {
				if cs.conn.RemoteAddr().String() == dirClientAc {
					clients[i].channel = canal
					//texto en consola client subscrito
					fmt.Printf("./%s subscribed -channel %d\n\n", clients[i].nombre, clients[i].channel)

					break
				}
			}

		case SND:
			nameFile := args[1]
			canal, _ := strconv.Atoi(args[2])
			canalNum := args[2]

			//./client send miarchivo.png -channel 1
			var clientSend string
			for i, cs := range clients {
				if client.RemoteAddr().String() == cs.conn.RemoteAddr().String() {
					//fmt.Printf("./%s send %s -channel %d\n", clients[i].nombre, nameFile, canal)
					clientSend = "./" + clients[i].nombre + " send " + nameFile + "-channel " + canalNum
				}

				if client.RemoteAddr().String() != cs.conn.RemoteAddr().String() && cs.channel == canal {
					fmt.Printf("./%s receive -channel %d\n\n", clients[i].nombre, canal)

				}
			}
			fmt.Println(clientSend)
			fmt.Println()

			nameFile = "./" + nameFile

			mA.RecibirArchivo(client, nameFile)

			writeMessageAllUsers(nameFile, client, canal)

		default:
			fmt.Printf("wrong entry: %s\n", texto)
		}

	}
}

func writeMessageAllUsers(ruta string, conn net.Conn, canal int) {
	//ruta = "go.png"
	for i, _ := range clients {
		if (clients[i].conn.RemoteAddr().String() != conn.RemoteAddr().String()) && clients[i].channel == canal {
			clients[i].conn.Write([]byte(ruta))
			mA.EnviarArchivo(clients[i].conn, ruta)
		}

	}

}
