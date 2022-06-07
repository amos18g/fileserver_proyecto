package manejoArchivos

import (
	"fmt"
	"io"
	"net"
	"os"
)

func EnviarArchivo(conn net.Conn, filePath string) {

	bandera := true
	//	}

	if bandera {
		//Read only open file
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("OS. Open() function execution error, error is:%v \n", err)
			return
		}

		defer file.Close()

		buf := make([]byte, 4096)

		var cont int
		for {
			//Read the data from the local file and write it to the network receiver. How much to read, how much to write
			n, erro := file.Read(buf)

			//Write to network socket
			_, err = conn.Write(buf[:n])

			cont += n
			if erro != nil || n < 4096 {
				if err == io.EOF || n < 4096 {
					//fmt.Println("El tamaño es ", cont)
					//fmt.Printf("sending file completed \n")

				} else {
					fmt.Printf("file. Read() method execution error, error is:%v \n", err)
				}
				return
			}

		}
	}

}

func RecibirArchivo(conn net.Conn, fileName string) {

	//Create a new file by file name
	file, err := os.Create(fileName)
	//fmt.Println("fileName is: ", fileName)
	if err != nil {
		fmt.Printf("OS. Create() function execution error, error is:% v \n", err)
		return
	}
	defer file.Close()

	//Read data from network and write to local file

	for {
		buf := make([]byte, 4096)
		n, err := conn.Read(buf)

		//Write to local file, read and write
		file.Write(buf[:n])
		if err != nil || n < 4096 {
			if err == io.EOF || n < 4096 {
				//fmt.Printf("receive file complete. \n")

			} else {
				fmt.Printf("conn.read() method execution error, error is:%v \n", err)
			}
			return
		}
	}

}
