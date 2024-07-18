package main

import ()

var receiptDB = map[string]Receipt{}

func main() {
	server := NewAPIServer(":8080")
	server.Run()
}
