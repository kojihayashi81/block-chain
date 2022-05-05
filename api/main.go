package main

import (
	"app/server"
	"flag"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 8080, "TCP Port Number Start Blockchain Server")
	flag.Parse()
	app := server.NewBlockchainServer(uint16(*port))
	fmt.Println("start")
	app.Run()
}
