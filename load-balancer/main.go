package main

import (
	"fmt"
	"log"
)

func main() {
	lb := NewLoadBalancer()
	fmt.Println(lb.getServer())
	fmt.Println(lb.getServer())
	fmt.Println(lb.getServer())
	fmt.Println(lb.getServer())
	fmt.Println(lb.getServer())
	err := lb.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
