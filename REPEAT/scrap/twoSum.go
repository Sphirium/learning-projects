package main

import "fmt"

func main() {

	nums := []int{20, 35, 56, 72, 89, 92, 98}
	target := 55
	result := twoSum(nums, target) // вызываем функцию
	fmt.Println(result)

}

func twoSum(nums []int, target int) []int {
	hashMap := make(map[int]int)

	for i, num := range nums {
		// target = addNumber + num ===> addNumber = target - num
		addNumber := target - num
		if j, ok := hashMap[addNumber]; ok {
			return []int{j, i}
		}
		hashMap[num] = i
	}
	return nil
}
