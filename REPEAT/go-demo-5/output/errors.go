package output

import (
	"github.com/fatih/color"
)

func PrintError(value any) {
	// Реализация метода Type Switch
	// switch t := value.(type) {
	// case string:
	// 	color.Red(t)
	// case int:
	// 	color.Red("Код ошибки: %d", t)
	// case error:
	// 	color.Red(t.Error())
	// default:
	// 	color.Red("Неизвестный тип ошибки")
	// }

	// Альтернативный метод получения типа элемента
	intValue, ok := value.(int)
	if ok {
		color.Red("Код ошибки: %d", intValue)
		return
	}
	strValue, ok := value.(string)
	if ok {
		color.Red("Код ошибки: %d", strValue)
		return
	}
	errorValue, ok := value.(error)
	if ok {
		color.Red("Код ошибки: %d", errorValue)
		return
	}
	color.Red("Неизвестный тип ошибки")

}
