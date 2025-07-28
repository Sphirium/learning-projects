package main

import "fmt"

func main() {
	a := 10
	changeValue(&a) // изменили значение по ссылке в память
	fmt.Println(a)  // выводится значение
	fmt.Println(&a) // выводится адрес переменной

	double(&a)      // удвоили значение
	fmt.Println(a)  // выводится значение
	fmt.Println(&a) // выводится адрес переменной

}

func changeValue(c *int) {
	*c = 50
}

func double(num *int) {
	*num = *num * 2
}
