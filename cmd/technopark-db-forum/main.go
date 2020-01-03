package main

import (
	"fmt"
	_ "github.com/lib/pq"
	apiserver "github.com/nozimy/technopark-db-forum/internal/app"
	"log"
)

func main() {
	fmt.Println("Running technopark-db-forum rest api")
	if err := apiserver.Start(); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
