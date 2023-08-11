package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"bytes"
	"time"
)

func main() {
	fmt.Print("Enter Server IP: ")
	var serverIP string
	var connected bool

	// Taking input from user
	fmt.Scanln(&serverIP)

	fmt.Print("Connecting")
	for !connected {
		serverIP = "127.0.0.1:8080"
		urlLive := "http://" + serverIP + "/live"
		resp, err := http.Get(urlLive)

		fmt.Print(".")

		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			sb := string(body)
			fmt.Println(sb)

			if sb != "OK\n\n" {
				fmt.Println("Server is live")
				connected = true
			}
		}

		time.Sleep(1 * time.Second)
	}

	for connected {
		displayMessages(serverIP)
		var message string
		fmt.Print("Enter message: ")
		fmt.Scanln(&message)

		resp, err := http.Post("http://"+serverIP+"/message", "text/plain", bytes.NewBufferString(message))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(resp.Status)
	}
}

func displayMessages(ip string) {
	fmt.Println("\nFetching messages")
	url := "http://" + ip + "/read"
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	fmt.Println(sb)
}