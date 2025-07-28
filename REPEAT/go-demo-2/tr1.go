package main

import "fmt"

// Задача:
// 1. В цикле спрашиваем ввод транзакций: -10, 10, 40.5 и и тд
// 2. Добавлять каждую в массив транзакций
// 3. Вывести массив
// 4. Вывести сумму баланса

func main() {
	transactions2 := make([]float64, 3, 6)
	fmt.Println(len(transactions2), cap(transactions2))
	transactions2 = append(transactions2, 1,1,1,2)
	fmt.Println(transactions2)
	fmt.Println(len(transactions2), cap(transactions2))

	transactions := []float64{}
	fmt.Println("Введите ваши транзакции за текущий месяц: ")

	for {
		transaction := getUserInput()
		if transaction == 0 {
			break
		}

		transactions = append(transactions, transaction)
	}
	balance := calculateBalance(transactions)
	fmt.Printf("Ваш баланс: %.2f", balance)

}

func getUserInput() float64 {
	var transaction float64
	fmt.Println("Введите транзакцию: (n для выхода): ")
	fmt.Scan(&transaction)
	return transaction
}

func calculateBalance(transactions []float64) float64 {
	balance := 0.0
	for _, value := range transactions {
		balance += value
	}
	return balance
}

