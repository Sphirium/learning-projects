package main

import (
	"fmt"
)

func main() {
	massive1 := []int{5, 10, 15, 20, 25, 45}
	demoMassive := massive1
	massive1 = append(massive1, 100)
	massive1[4] = 500
	slice1 := massive1[:2]
	slice1[1] = 30
	slice1 = slice1[:5]

	persons := [5]string{"one", "two", "three", "four", "five"}

	fmt.Println(massive1)
	fmt.Println(demoMassive)
	fmt.Println(massive1[1:])
	fmt.Println(len(persons), cap(persons))
	fmt.Println(persons)
	fmt.Println(slice1)
	fmt.Println(len(slice1), cap(slice1))

}
