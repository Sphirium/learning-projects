package account

import (
	"errors"
	"math/rand/v2"
	"net/url"
	"time"

	"github.com/fatih/color"
)

var symbols = []rune("0123456789ABCDEFGHIJKLMNOPQabcdefghijklmnopq_`=-+")

type Account struct {
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	Url        string    `json:"url"`
	CreatedAcc time.Time `json:"CreatedAcc"`
	UpdatedAcc time.Time `json:"UpdatedAcc"`
}

func (acc *Account) Output() {
	color.Cyan(acc.Login)
	color.HiGreen(acc.Password)
	color.Yellow(acc.Url)

}

func (acc *Account) GeneratePassword(n int) {
	res := make([]rune, n)
	for i := range res {
		res[i] = symbols[rand.IntN(len(symbols))]
	}
	acc.Password = string(res)
}

func NewAccount(login, password, urlString string) (*Account, error) {
	if login == "" { // 1. Если логина нет, то ошибка
		return nil, errors.New("Invalid_LOGIN")
	}
	_, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, errors.New("Invalid_URL")
	}
	newAcc := &Account{
		CreatedAcc: time.Now(),
		UpdatedAcc: time.Now(),
		Login:      login,
		Password:   password,
		Url:        urlString,
	}

	if password == "" {
		newAcc.GeneratePassword(12)
	}
	return newAcc, nil
}

// 2. Если нет пароля, то генерим
// func newAccount(login, password, urlString string) (*account, error) {
// 	newAcc := &account{
// 		login:    login,
// 		password: password,
// 		url:      urlString,
// 	}

// 	if login == "" { // 1. Если логина нет, то ошибка
// 		return nil, errors.New("Invalid_LOGIN")
// 	}

// 	_, err := url.ParseRequestURI(urlString)
// 	if err != nil {
// 		return nil, errors.New("Invalid_URL")
// 	}

// 	if password == "" {
// 		newAcc.generatePassword(12)
// 	}
// 	return newAcc, nil
// }
