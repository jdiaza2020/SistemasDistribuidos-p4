/*
ESTE ES EL ÚNICO ARCHIVO QUE SE PUEDE MODIFICAR

RECOMENDACIÓN: Solo modicar a partir de la parte
				donde se encuentran la explicación
				de las otras variables.

*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
	msg    string
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			msg = string(buf[:n])
			/*
				Desde aquí debería salir la información a una goroutine o a una función ordinaria según se requiera

				n 	:=	Longitud del string enviado por el servidor
				msg := 	Mensaje recibido por la conexión del server

				Recomendación: Para usar la conversión entre string e int se recomienda usar strconv.Atoi
				Más info en: https://pkg.go.dev/strconv#Atoi

			*/
			fmt.Println("len: " + strconv.Itoa(n) + " msg: " + msg)
		}
	}
}
