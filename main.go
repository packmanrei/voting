package main

import (
	controller "voting/controllers"
)

func main() {
	router := controller.StartServer()
	router.Run()
}
