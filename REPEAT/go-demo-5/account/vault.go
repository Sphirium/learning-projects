package account

import (
	"demo/passwords/encrypter"
	"demo/passwords/output"
	"encoding/json"
	"strings"
	"time"

	"github.com/fatih/color"
)

type ByteReader interface {
	Read() ([]byte, error)
}

type ByteWriter interface {
	Write([]byte)
}

type Db interface {
	ByteReader
	ByteWriter
}

type Vault struct {
	Accounts   []Account `json:"accounts"`
	UpdatedAcc time.Time `json:"updatedAcc"`
}

// функция создает новое хранилище, если его нет,
// а также читает данные из файла json
type VaultWithDb struct {
	Vault
	db  Db
	enc encrypter.Encrypter
}

func NewVault(db Db, enc encrypter.Encrypter) *VaultWithDb {
	file, err := db.Read()
	if err != nil {
		return &VaultWithDb{
			Vault: Vault{
				Accounts:   []Account{}, // Пустой список аккаунтов
				UpdatedAcc: time.Now(),  // Текущая дата/время
			},
			db:  db,  // Сохраняем переданную базу
			enc: enc, // Формируем шифрование
		}
	}
	data := enc.Decrypt(file)
	var vault Vault
	err = json.Unmarshal(data, &vault)
	color.Yellow("Найдено %d аккаунтов", len(vault.Accounts))
	if err != nil {
		output.PrintError("Не удалось разобрать файл data.vault")
		return &VaultWithDb{
			Vault: Vault{
				Accounts:   []Account{}, // Пустой список
				UpdatedAcc: time.Now(),  // Текущая дата
			},
			db:  db, // Сохраняем базу
			enc: enc,
		}
	}
	return &VaultWithDb{
		Vault: vault, // Загруженные из JSON данные
		db:    db,    // Переданная база
		enc:   enc,
	}
}

func (vault *VaultWithDb) FindAccounts(str string, checker func(Account, string) bool) []Account {
	var accounts []Account
	for _, account := range vault.Accounts {
		isMatched := checker(account, str)
		if isMatched {
			accounts = append(accounts, account)
		}
	}
	return accounts
}

func (vault *VaultWithDb) DeleteAccountByUrl(url string) bool {
	var accounts []Account
	isDeleted := false
	for _, account := range vault.Accounts {
		isMatched := strings.Contains(account.Url, url)
		if !isMatched {
			accounts = append(accounts, account)
			continue
		}
		isDeleted = true
	}
	vault.Accounts = accounts
	vault.save()
	return isDeleted
}

// данный метод реализует добавление нового
// аккаунта в существующее хранилище (vault)
func (vault *VaultWithDb) AddAccount(acc Account) {
	vault.Accounts = append(vault.Accounts, acc)
	vault.save()
}

// Метод переводит данные в массив байт, чтобы их передать в файл json
func (vault *Vault) ToBytes() ([]byte, error) {
	// fmt.Println(acc)
	file, err := json.Marshal(vault)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Метод при его вызове сохраняет данные хранилища (vault)
func (vault *VaultWithDb) save() {
	vault.UpdatedAcc = time.Now()
	data, err := vault.Vault.ToBytes()
	encData := vault.enc.Encrypt(data)
	if err != nil {
		output.PrintError("Не удалось преобразовать")
	}
	vault.db.Write(encData)
}
