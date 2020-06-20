package main

import "github.com/tachanokkik/finalexam/customer"

func main() {
	r := customer.SetupRouter()
	r.Run(":2019")
}
