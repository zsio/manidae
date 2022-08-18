package main

import "manidae/client"

func main() {

	// 关机
	shutdown := make(chan struct{})

	go client.Start()

	<-shutdown
}
