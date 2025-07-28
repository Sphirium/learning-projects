package main

import (
	"fmt"
)

func main() {
	arr := [5]int{1, 2, 3, 4, 5}
	reverse(&arr)
	fmt.Println(arr)
}

func reverse(numbers *[5]int) {
	for index, value := range *numbers {
		numbers[len(numbers)-1-index] = value

	}

}
