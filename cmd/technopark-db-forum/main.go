package main

import (
	"fmt"
	"log"
	apiserver "technopark-db-forum/internal/app"
)

func main() {
	fmt.Println("Running technopark-db-forum rest api")
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
