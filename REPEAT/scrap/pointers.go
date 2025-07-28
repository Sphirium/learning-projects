package main

import "fmt"

func main() {
	a := 10
	changeValue(&a)
	fmt.Println(a)
}

func changeValue(c *int) {
	*c = 50

}
