package main

import (
	"demo/weather/geo"
	"demo/weather/weather"
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Новый проект")
	city := flag.String("city", "", "Город пользователя")
	format := flag.Int("format", 2, "Формат выводы погоды")
	// format := flag.Int("age", 18, "Формат вывода погоды")

	flag.Parse()

	fmt.Println(*city)
	geoData, err := geo.GetMyLocation(*city)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(*geoData)

	fmt.Printf("Город: %s\n", geoData.City)
	weatherData, err := weather.GetWeather(*geoData, *format)
	fmt.Println(weatherData)
	// reader := strings.NewReader("Привет! Я поток данных")
	// blockBytes := make([]byte, 4)
	// for {
	// 	_, err := reader.Read(blockBytes)
	// 	fmt.Printf("%q\n", blockBytes)
	// 	if err == io.EOF {
	// 		break
	// 	}
	// }

}
