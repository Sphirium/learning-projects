package main

import "fmt"

func main() {

tr1:= []int{1,2,3}
tr2:=[]int{4,5,6}
tr1 = append(tr1, tr2...)
fmt.Println(tr1)
for i, value := range tr1 {
fmt.Println(i, value)
}
	
transactions := []float64{}
	fmt.Println("Введите ваши транзакции за текущий месяц: ")

	for {
		transaction := getUserInput()
		if transaction == 0 {
			break
		}

		transactions = append(transactions, transaction)
	}
	fmt.Println(transactions)
}

func getUserInput() float64 {
	var transaction float64
	fmt.Print("Введите транзакцию: (n для выхода)")
	fmt.Scan(&transaction)
	return transaction
}
