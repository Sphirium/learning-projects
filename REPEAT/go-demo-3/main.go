package main

import "fmt"

/*
Создать приложение, которое сначала выдает меню:
1. Посмотреть закладки
2. Добавить закладку
3. Удалить закладку
4. Выход
При 1 - Выводит закладки
При 2 - два поля ввода названия и адресе, а после добавление
При 3 - Ввод названия и удаление по нему
При 4 - Завершение
*/
var bookmarks = make(map[string]string)

func main() {
Menu:
	for {
		userInput := getMenu()
		switch userInput {
		case 1:
			reviewBookmarks(bookmarks)
		case 2:
			addBookmark(bookmarks)
		case 3:
			deleteBookmark(bookmarks)
		case 4:
			break Menu
		}
	}
}

func getMenu() int {
	var userInput int
	fmt.Println("Выберите вариант: ")
	fmt.Println("1. Посмотреть закладки")
	fmt.Println("2. Добавить закладку")
	fmt.Println("3. Удалить закладку")
	fmt.Println("4. Выход")
	fmt.Scan(&userInput)
	return userInput
}

func reviewBookmarks(bookmarks map[string]string) {
	// bookmarks = map[string]string {
	// 	"Google": "google.com",
	// 	"Yandex": "yandex.ru",
	// 	"Facebook": "facebook.com",
	// }
	if len(bookmarks) == 0 {
		fmt.Println("Пока нет закладок")
	}
	for key, value := range bookmarks {
		fmt.Println(key, ": ", value)
	}
}

func addBookmark(bookmarks map[string]string) {
	var newBookmarkKey string
	var newBookmarkValue string
	fmt.Println("Введите название: ")
	fmt.Scan(&newBookmarkKey)
	fmt.Println("Введите ссылку: ")
	fmt.Scan(&newBookmarkValue)
	bookmarks[newBookmarkKey] = newBookmarkValue

}

func deleteBookmark(bookmarks map[string]string) {
	var bookmarkKeyToDelete string
	fmt.Println("Введите название: ")
	fmt.Scan(&bookmarkKeyToDelete)
	delete(bookmarks, bookmarkKeyToDelete)

}
