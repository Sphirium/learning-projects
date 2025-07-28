package main

import "fmt"

type Person struct { //объявляем структуру (набор данных по смысловому значению)
	Age     int
	Name    string
	Surname string
}

func (p *Person) Talk() { //реализация метода Talk() - то есть "говорит"

	fmt.Printf("Привет, меня зовут %s и мне %d лет\n", p.Name, p.Age)

}

func (p *Person) IncAge() {

	p.Age = p.Age + 2
}

func main() {
	ivan := Person{
		Age:     18,
		Name:    "Ivan",
		Surname: "Borisov",
	}
	fmt.Println(ivan)
	ivan.Talk()
	ivan.IncAge() // данный метод увеличивает возраст на 2 года
	fmt.Println(ivan)
	ivan.Talk()

}
