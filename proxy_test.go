package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

var end = make(chan bool)

func TestProxy(t *testing.T) {
	go runServer()
	time.Sleep(3 * time.Second) // Wait for server start

	msg, cod, err := getRequest("http://localhost:12345/?url=https%3A%2F%2Fwww.google.com%2F")
	if err != nil {
		t.Error("Service is not started? cannot getRequest", err)
	}
	if cod != 200 {
		t.Error("Query service isn't posible, HTTP error code: ", cod, msg)
	}
	if msg == "The URL is not in whitelist.lst.\n" {
		t.Error("Is not query www.google.com")
	}

	msg, cod, err = getRequest("http://localhost:12345/?url=https%3A%2F%2Fwww.gmail.com%2F")
	if err != nil {
		t.Error("Service is not started? cannot getRequest", err)
	}
	if cod != 200 {
		t.Error("Query service isn't posible, HTTP error code: ", cod, msg)
	}
	if msg != "The URL is not in whitelist.lst.\n" {
		t.Error("Is query www.gmail.com")
	}

	end <- true
	time.Sleep(1 * time.Second) // Wait for server end
}

func runServer() {
	cmd := exec.Command("go", "run", "proxy.go")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()

	time.Sleep(5 * time.Second)

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15) // note the minus sign
	}

	msg := <-end // wait for message
	fmt.Println("Exit go run proxy.go: ", msg)
	cmd.Wait() // Exit?
}

func getRequest(url string) (string, int, error) {
	client := http.Client{} // request client
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error GET: ", err)
		return "", -1, err
	}

	// execute petition
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error Client DO: ", err)
		return "", -1, err
	}
	defer response.Body.Close() // response can be nil

	// get response data
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error responseData: ", err)
		return "", response.StatusCode, err
	}

	return string(responseData), response.StatusCode, nil
}
