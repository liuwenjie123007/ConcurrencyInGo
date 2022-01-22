package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	wg.Add(4)
	go sayHello()

	go func() {
		defer wg.Done()
		fmt.Println("hello")
	}()

	sayHello := func() {
		defer wg.Done()
		fmt.Println("hello")
	}

	go sayHello()

	salutation := "hello"
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()

	wg.Wait()
	fmt.Println(salutation)

	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()

	for _, salutation := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()

}

func sayHello() {
	defer wg.Done()
	fmt.Println("hello")
}
