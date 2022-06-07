package main

import (
	"bufio"
	mA "control/manejo_archivos"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CH   = "channel"
	SUBS = "subscribe"
	SND  = "send"
)

func conectar() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")

	if err != nil {
		fmt.Printf("net. Dial() function execution error, error is:%v \n", err)
		os.Exit(0)
	}

	return conn
}

func main() {

	//Initiate connection request actively

	//defer conn.Close()
	conn := conectar()

	path := "./"

	go listenServer(conn)

	var scanner = bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var msg = scanner.Text()
		fmt.Println()
		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")

		if len(args) == 3 || len(args) == 4 {

			cmd := strings.TrimSpace(args[0])

			switch cmd {
			case SUBS:
				//codigo para subscripcion
				if strings.TrimSpace(args[1]) == CH && len(args) == 3 {
					fmt.Println("successfully subscribed")

					conn.Write([]byte(CH + " " + args[2]))

				}

			case SND:
				//codigo para enviar
				if args[1] == CH && len(args) == 4 {
					fmt.Println("successfully sent")
					path += args[3]
					nameFile := args[3]
					infoChannel := args[2]

					//verificando que el archivo existe
					_, err := os.Stat(path)
					if err != nil {
						fmt.Printf("OS. Stat() function execution error, error is:%v \n", err)
						path = "./"
						continue
					}

					infoFile := "send " + nameFile + " " + infoChannel
					conn.Write([]byte(infoFile))

					mA.EnviarArchivo(conn, path)

					path = "./"

				}

			default:
				fmt.Printf("wrong entry: %s\n", msg)
			}

		} else {
			fmt.Printf("wrong entry: %s\n", msg)
		}

	}

}

func listenServer(cnn net.Conn) {
	//read message
	for {

		var buff = make([]byte, 2048)
		n, err := cnn.Read(buff)

		if err != nil {

			return
		}

		//se ejecuta una vez cada vez que se hace enter
		//var texto = scanner.Text()
		var texto = string(buff[:n])

		//texto = "./" + texto

		if texto != "" {
			mA.RecibirArchivo(cnn, texto)
		}

	}

}
