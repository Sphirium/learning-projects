package main

import "fmt"

func main()  {
	myMap := map[string]string{
		"VKONTAKTE": "https://vk.com",
	}
	fmt.Println(myMap)
	fmt.Println(myMap["VKONTAKTE"])
	myMap["VKONTAKTE"] = "https://vkontakteee.com"
	myMap["OLKA-POLKA"] = "SHVEYA-NUMBER-ONE.RU"
	fmt.Println(myMap)
	myMap["Balamut"] = "Balamut.RU"	
	delete(myMap, "Balamut")
	fmt.Println(myMap)

	// Итерация по мапе:

	for key, value := range myMap {
		fmt.Println(key, value)
	}

}








