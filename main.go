package main

import (
	"log"
	"time"

	"github.com/chandan/trading-dashboard/routes"
	"github.com/chandan/trading-dashboard/services"
)

func main() {

	go services.StartPriceSimulator(time.Second * 2)

	router := routes.SetupRouter()

	log.Println("Server running on 8080")
	router.Run(":8080")
}
