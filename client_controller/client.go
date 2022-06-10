package client_controller

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	conFiles "control/controller_files"
)

const (
	CH   = "channel"
	SUBS = "subscribe"
	SND  = "send"
)

func Client() {
	//Initiate connection request actively
	conn := connect()
	go listenServer(conn)
	OptionClient(conn)

}

func connect() net.Conn {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")

	if err != nil {
		fmt.Printf("net. Dial() function execution error, error is:%v \n", err)
		os.Exit(0)
	}

	return conn
}

func OptionClient(conn net.Conn) {
	path := "./"

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
				//subscription code
				if strings.TrimSpace(args[1]) == CH && len(args) == 3 {
					fmt.Println("successfully subscribed")

					conn.Write([]byte(CH + " " + args[2]))

				}

			case SND:
				//code to send
				if args[1] == CH && len(args) == 4 {

					path += args[3]
					nameFile := args[3]
					infoChannel := args[2]

					//verifying that the file exists
					_, err := os.Stat(path)
					if err != nil {
						fmt.Printf("OS. Stat() function execution error, error is:%v \n", err)
						path = "./"
						continue
					} else {
						fmt.Println("successfully sent")
					}

					infoFile := "send " + nameFile + " " + infoChannel
					conn.Write([]byte(infoFile))

					conFiles.SendFiles(conn, path)
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

		var texto = string(buff[:n])

		if texto != "" {
			conFiles.ReceiveFiles(cnn, texto)
		}

	}

}
