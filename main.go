package main

import "fmt"

func main() {
	// alertHandler := AlertHandler{}
	fmt.Println("Hello world")
}

type AlertHandler struct {
	AlertsByService map[ServiceKey][]Alert
}
