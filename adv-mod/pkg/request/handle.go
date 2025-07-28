package request

import (
	"adv-mod/pkg/response"
	"net/http"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	body, err := Decode[T](r.Body)
	if err != nil {
		response.Json(*w, err.Error(), 402)
		return nil, err
	}

	err = IsValid(body)
	if err != nil {
		response.Json(*w, err.Error(), 402)
		return nil, err
	}
	return &body, nil

}

// if payload.Email == "" {
// 	response.Json(w, "Email required", 402)
// 	return
// }
// // Compile требует ввод регулярных выражений
// _, err = mail.ParseAddress(payload.Email)
// if err != nil {
// 	response.Json(w, "Wrong email", 402)
// 	return
// }

// if payload.Password == "" {
// 	response.Json(w, "Password required", 402)
// 	return
// }

// ВСЁ, ЧТО ВЫШЕ. ДАЕТ ВОЗМОЖНОСТЬ ВАЛИДАЦИИ ПОЛУЧЕННЫХ ДАННЫХ,
// НО МОЖНО ПРОВЕСТИ ВАЛИДАЦИЮ С ПОМОЩЬЮ БИБЛИОТЕКИ - github.com/go-playground/validator/v10
