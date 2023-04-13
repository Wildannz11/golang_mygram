package main

import (
	"fmt"
	"os"
	database "project4/databases"
	"project4/helpers"
	router "project4/routers"
)

func main() {
	helpers.LoadEnv()
	r := router.StartApp()
	database.StartDB()

	serverPort := os.Getenv("SERVER_PORT")
	r.Run(fmt.Sprintf(":%v", serverPort) )
}
