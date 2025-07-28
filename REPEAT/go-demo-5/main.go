package main

import (
	"demo/passwords/account"
	"demo/passwords/encrypter"
	"demo/passwords/files"
	"demo/passwords/output"
	"fmt"
	"strings"

	"github.com/joho/godotenv"

	"github.com/fatih/color"
)

var menu = map[string]func(*account.VaultWithDb){
	"1": createAccount,
	"2": findAccountByUrl,
	"3": findAccountByLogin,
	"4": deleteAccount,
}

var userInputVariants = []string{
	"1. Создать аккаунт",
	"2. Найти аккаунт по URL",
	"3. Найти аккаунт по LOGIN",
	"4. Удалить аккаунт",
	"5. Выход",
	"Выберите вариант",
}

// Замыкание - процесс, когда функция возвращает другую функцию.
// Эта функция считает, сколько раз мы вызвали меню (как пример).
// func menuCounter() func() {
// 	i := 0
// 	return func() {
// 		i++
// 		fmt.Println(i)
// 	}
// }

func main() {
	fmt.Println("___Менеджер паролей___")
	err := godotenv.Load()
	if err != nil {
		output.PrintError("Не удалось найти env-файл")
	}
	// vault := account.NewVault(cloud.NewCloudDb("https://google.com"))
	vault := account.NewVault(files.NewJsonDb("data.vault"), *encrypter.NewEncrypter())

	// infoEnv := os.Getenv("VAR")
	// fmt.Println(infoEnv)

	// for _, env := range os.Environ() {
	// 	parts := strings.SplitN(env, "=", 2)
	// 	fmt.Println(parts[0])
	// }

	// counter := menuCounter()

Menu:
	for {
		// counter()
		userInput := promptData(userInputVariants...)
		menuFunc := menu[userInput]
		if menuFunc == nil {
			break Menu
		}
		menuFunc(vault)
		// 	switch userInput {
		// 	case "1":
		// 		createAccount(vault)
		// 	case "2":
		// 		findAccount(vault)
		// 	case "3":
		// 		deleteAccount(vault)
		// 	default:
		// 		break Menu
		// 	}
		// }

	}
}

// func getMenu() int {
// 	var userInput int
// 	fmt.Println("Выберите вариант: ")
// 	fmt.Println("1. Создать аккаунт")
// 	fmt.Println("2. Найти аккаунт")
// 	fmt.Println("3. Удалить аккаунт")
// 	fmt.Println("4. Выход")
// 	fmt.Scanln(&userInput)
// 	return userInput
// }

func deleteAccount(vault *account.VaultWithDb) {
	url := promptData("Введите URL для поиска")
	isDeleted := vault.DeleteAccountByUrl(url)
	if isDeleted {
		color.Green("Удалено!")
	} else {
		output.PrintError("Не найдено")

	}
}

func findAccountByUrl(vault *account.VaultWithDb) {
	url := promptData("Введите URL для поиска")
	accounts := vault.FindAccounts(url, func(acc account.Account, str string) bool { // объявление анонимной функции на месте, если она используется один раз
		return strings.Contains(acc.Url, str)
	})
	for _, account := range accounts {
		account.Output()
	}
	if len(accounts) == 0 {
		output.PrintError("Аккаунт не найден")
	}
}

func findAccountByLogin(vault *account.VaultWithDb) {
	login := promptData("Введите LOGIN для поиска")
	accounts := vault.FindAccounts(login, func(acc account.Account, str string) bool { // объявление анонимной функции на месте, если она используется один раз
		return strings.Contains(acc.Login, str)
	})
	for _, account := range accounts {
		account.Output()
	}
	if len(accounts) == 0 {
		output.PrintError("Аккаунт не найден")
	}
}

// func checkUrl(acc account.Account, str string) bool {
// 	return strings.Contains(acc.Url, str)
// }

func createAccount(vault *account.VaultWithDb) { //
	// files.ReadFile()
	// files.WriteFile([]byte("Привет! Я файл"), "file.txt")
	login := promptData("Введите логин")
	password := promptData("Введите пароль")
	url := promptData("Введите URL")
	myAccount, err := account.NewAccount(login, password, url)
	if err != nil {
		fmt.Println("INVALID_URL")
		return
	}
	vault.AddAccount(*myAccount)
}

// Функция считывает ввод юзера и возвращает его

func promptData(prompt ...string) string {
	for i, line := range prompt {
		if i == len(prompt)-1 {
			fmt.Printf("%v: ", line)
		} else {
			fmt.Println(line)
		}
	}
	var result string
	fmt.Scanln(&result)
	return result
}

// func promptData(prompt string) string {
// 	fmt.Println(prompt + ": ")
// 	var result string
// 	fmt.Scanln(&result)
// 	return result
// }

// func outputPassword(acc *account) {
// 	// fmt.Println(acc)
// 	fmt.Println(acc.login, acc.password, acc.url)
// }

// fmt.Println(generatePassword(12))

// Работа с рунами:

// str := []rune("Привет!)")
// for _, ch := range string(str) {
// 	fmt.Println(ch, string(ch))
// }
