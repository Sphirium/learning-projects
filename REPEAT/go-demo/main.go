package main

import (
	"errors"
	"fmt"
	"math"
)

func main() {
	// for i := 0; i < 10; i++ {
	// 	if i == 5 {
	// 		break
	// 	}
	// 	fmt.Printf("%d\n", i)
	fmt.Println("___Калькулятор индекса массы тела___")
	// }
	// for i := 0; i < 10; i++ {
	for {

		height, weight := getUserInput()
		IMT, err := calculateIMT(height, weight)
		if err != nil {
			// fmt.Println("Неверные параметры расчета. Повторите попытку.")
			// continue
			panic("Неверные параметры расчета. Повторите попытку.")
		}
		outputResult(IMT)
		isRepeatCalculation := checkRepeatCalculation()
		if !isRepeatCalculation {
			break

		}
	}

	// if IMT >= 0 && IMT <= 16 {
	// 	fmt.Println("У вас сильный дефицит массы тела")
	// }
	// if IMT >= 16 && IMT <= 18.5 {
	// 	fmt.Println("У вас дефицит массы тела")
	// }
	// if IMT >= 18.5 && IMT <= 25 {
	// 	fmt.Println("У нормальный вес")
	// }
	// if IMT >= 25 && IMT <= 30 {
	// 	fmt.Println("У вас избыточная масса")
	// }
	// if IMT >= 30 && IMT <= 35 {
	// 	fmt.Println("У вас 1-я степень ожирения")
	// }

}

func outputResult(IMT float64) {
	result := fmt.Sprintf("Ваш индекса массы тела: %.2f", IMT)
	fmt.Println(result)
	switch {
	case IMT < 16:
		fmt.Println("У вас сильный дефицит массы тела")
	case IMT < 18.5:
		fmt.Println("У вас дефицит массы тела")
	case IMT < 25:
		fmt.Println("У вас нормальный вес")
	case IMT < 30:
		fmt.Println("У вас избыточная масса")
	case IMT < 35:
		fmt.Println("У вас 1-я степень ожирения")
	default:
		fmt.Println("ХВАТИТ ЖРАТЬ!")

	}
}

func calculateIMT(height float64, weight float64) (float64, error) {
	if height <= 0 || weight <= 0 {
		return 0, errors.New("No_params_error")
	}
	const IMTPower = 2
	IMT := weight / math.Pow(height/100, IMTPower)
	return IMT, nil
}

func getUserInput() (float64, float64) {
	var height float64
	var weight float64
	fmt.Print("Введите свой рост в сантиметрах: ")
	fmt.Scan(&height)
	fmt.Print("Введите свой вес: ")
	fmt.Scan(&weight)
	return height, weight
}

func checkRepeatCalculation() bool {
	var userChoise string
	fmt.Println("Сделать еще расчет (y/n)?")
	fmt.Scan(&userChoise)
	if userChoise == "y" || userChoise == "Y" {
		return true
	}
	return false

}
