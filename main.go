package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"slices"
	"strconv"
	"strings"
)

var valid_methods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE"}
var valid_protocol_versions = []string{"HTTP/1.1", "HTTP/1.0"}

func main() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error Listening on Port 8080")
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection on Port 8080")
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	// Close connection after function exits
	defer conn.Close()
	// reader := bufio.NewRead(conn)

	buffer := make([]byte, 1024)
	var active_buffer []byte
	// var active_buffer []byte

	delim := []byte("\r\n\r\n")
	var index int
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed connection")
			} else {
				log.Println("Read Error: ", err)
			}
			return
		}
		active_buffer = append(active_buffer, buffer[:n]...)
		index = bytes.Index(active_buffer, delim)
		if index != -1 {
			break
		}
	}
	// Removed including the delimiter because was messing up later, keep in mind
	request_data := active_buffer[:index]

	fields := strings.Split(string(request_data), "\r\n")

	request_line := strings.Split(fields[0], " ")
	// fmt.Println(requst)
	bad_request := false
	if len(request_line) == 3 && slices.Contains(valid_methods, request_line[0]) && len(request_line[1]) > 0 && slices.Contains(valid_protocol_versions, request_line[2]) {
		method := request_line[0]
		target := request_line[1]
		header_map := make(map[string]string)
		// fmt.Println(fields)
		if len(fields) > 1 {
			for _, header := range fields[1:] {
				// fmt.Println(header)
				key, value, found := strings.Cut(header, ":")
				// fmt.Println(header)
				// fmt.Println(found)
				if !found {
					bad_request = true
					break

				} else {
					header_map[key] = value
				}

			}

			if !bad_request {

				response_body := "Method: " + method + "\n" + "URI: " + target + "\n"
				for key, value := range header_map {
					response_body += key + ": " + value + "\n"
				}
				// fmt.Println(len(response_body))
				response := "HTTP/1.1 200 OK\r\n" +
					"Content-Type: text/plain\r\n" +
					"Content-Length: " + strconv.Itoa(len(response_body)) + "\r\n" +
					"Connection: close\r\n\r\n" + response_body

				_, err := conn.Write([]byte(response))

				if err != nil {
					log.Println("Write Error: ", err)
				}
				return

			}
		}
	}
	// fmt.Printf("Headers: %s\n", request_line)
	// httpRequest := "HTTP/1.1 200 OK\r\n" +
	// 	"Content-Type: text/plain\r\n" +
	// 	"Content-Length: 4\r\n" +
	// 	"Connection: close\r\n" +
	// 	"\r\n" +
	// 	"pong"

	httpError := "HTTP/1.1 400 Bad Request\r\n\r\n"
	_, err := conn.Write([]byte(httpError))

	if err != nil {
		log.Println("Write Error: ", err)
	}

}
