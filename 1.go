package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type User struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Функция для загрузки данных из файла
func loadUsers(filename string) ([]User, error) {
	// Проверяем, существует ли файл
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Если файл не существует, возвращаем пустой список
		return []User{}, nil
	}

	// Читаем содержимое файла
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Десериализуем JSON в список пользователей
	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Функция для сохранения данных в файл
func saveUsers(filename string, users []User) error {
	// Сериализуем список пользователей в JSON
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	// Записываем данные в файл
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Функция для проверки наличия пользователя по никнейму
func findUserByNickname(users []User, nickname string) *User {
	for _, user := range users {
		if user.Nickname == nickname {
			return &user
		}
	}
	return nil
}

// Структура для JSON-ответа
type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}
type ResponseAuthorisation struct {
	Status int `json:"status"`
}

const filename = "челы.aaa"

func main() {
	// Обработчик для корневого маршрута "/"
	fmt.Println(1)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовок Content-Type
		w.Header().Set("Content-Type", "application/json")
		// Проверяем метод запроса
		switch r.Method {
		case http.MethodGet:
			// Ответ на GET-запрос

			resp := Response{
				Message: "Добро пожаловать на сервер!",
				Status:  200,
			}
			json.NewEncoder(w).Encode(resp)

		case http.MethodPost:
			// Ответ на POST-запрос
			var requestBody map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&requestBody)
			if err != nil {
				http.Error(w, "Неверный формат данных", http.StatusBadRequest)
				return
			}
			if requestBody["type"] == "authorization" { // 1
				users, _ := loadUsers("челы.aaa")
				user := findUserByNickname(users, string(fmt.Sprintf("%v", requestBody["nickname"])))
				status := 0
				if user != nil {
					if requestBody["email"] == user.Email && user.Password == requestBody["password"] {
						status = 1
					} else if requestBody["email"] == user.Email {
						status = 2
					} else {
						status = 3
					}
				} else {
					status = 3
				}
				resp := ResponseAuthorisation{
					Status: status,
				}
				json.NewEncoder(w).Encode(resp)
			} else if requestBody["type"] == "create-account" { // 1
				users, _ := loadUsers(filename)
				user := findUserByNickname(users, string(fmt.Sprintf("%v", requestBody["nickname"])))
				status := 0
				if user != nil {
					status = 3
				} else {
					status = 1
					users = []User{{
						Nickname: fmt.Sprintf("%v", requestBody["nickname"]),
						Email:    fmt.Sprintf("%v", requestBody["email"]),
						Password: fmt.Sprintf("%v", requestBody["password"])}}
					saveUsers(filename, users)
				}
				resp := ResponseAuthorisation{
					Status: status,
				}
				json.NewEncoder(w).Encode(resp)
			} else if requestBody["type"] == "create-task" { // 3

			} else if requestBody["type"] == "add-task" { // 4

			} else if requestBody["type"] == "remove-task" { // 7

			} else if requestBody["type"] == "appointment-task" { // 6

			} else if requestBody["type"] == "connect" { // 2

			} else if requestBody["type"] == "update-task" { // 5

			}
			// Отправляем ответ с данными из запроса
			resp := Response{
				Message: fmt.Sprintf("Данные получены: %v", requestBody),
				Status:  200,
			}
			json.NewEncoder(w).Encode(resp)

		default:
			// Ошибка для не поддерживаемых методов
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})
	fmt.Println(2)

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
