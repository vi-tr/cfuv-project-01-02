#include <cstdio>
#include <httplib.h>
#include <nlohmann/json.hpp>
using json = nlohmann::json;

int main() {
    using namespace httplib;
    Server svr;

    if (!svr.is_valid()) {
	printf("Some error occured with the server");
	return -1;
    }

    svr.Get("/oauth", [](const Request &req, Response &res) {
	res.set_content("tessting","text/plain");
    });

    svr.set_error_handler([](const Request &req, Response &res) {
	res.set_content("[]", "text/html");
    });

    svr.set_logger([](const Request &req, const Response &res) {
	printf("Something has occured");
    });

    svr.listen("localhost", 8081);

    return 0;
}

/*
type User struct {
	Name        string
	UserID      int64
	AccessToken string
	Access      []int
}

type UserData struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// Глобальная переменная для проверки, что пользователь дал доступ
var authenticate struct {
	is_done bool
	code    string
}

func main() {
	// Список пользователей изначально пуст
	users := make(map[int64]*User)

	go startServer()

	var authURL string = "https://github.com/login/oauth/authorize?client_id=" + CLIENT_ID
	for !authenticate.is_done {
		fmt.Println("Чтобы зайти перейдите по ссылке:\n" + authURL)
		fmt.Println("и нажмите Enter")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	accessToken := getAccessToken(authenticate.code)
	userData := getUserData(accessToken)

	if _, ok := users[userData.Id]; !ok {
		// Добавляем пользователя с дефолтными правами
		users[userData.Id] = &User{
			Name:        userData.Name,
			UserID:      userData.Id,
			AccessToken: accessToken,
			Access:      []int{13},
		}
	}
	user := users[userData.Id]
	fmt.Println("Добро пожаловать,", user.Name)

	// Авторизация
	fmt.Print("В какую зону хотите попасть? ")
	var area int
	fmt.Scan(&area)

	if !slices.Contains(user.Access, area) {
		fmt.Print("Нет доступа в эту зону")
		return
	}
	fmt.Print("Доступ получен")
}

// Запуск сервера
func startServer() {
	http.HandleFunc("/oauth", handleOauth) // Вызов функции при запросе на /oauth
	http.ListenAndServe(":8080", nil)      // Запуск сервера
}

// Обработчик запроса
func handleOauth(w http.ResponseWriter, r *http.Request) {
	var responseHtml = "<html><body><h1>Вы НЕ аутентифицированы!</h1></body></html>"

	code := r.URL.Query().Get("code") // Достаем временный код из запроса
	if code != "" {
		authenticate.is_done = true
		authenticate.code = code
		responseHtml = "<html><body><h1>Вы аутентифицированы!</h1></body></html>"
	}

	fmt.Fprint(w, responseHtml) // Ответ на запрос
}

// Меняем временный код на токен доступа
func getAccessToken(code string) string {
	// Создаём http-клиент с дефолтными настройками
	client := http.Client{}
	requestURL := "https://github.com/login/oauth/access_token"

	// Добавляем данные в виде Формы
	form := url.Values{}
	form.Add("client_id", CLIENT_ID)
	form.Add("client_secret", CLIENT_SECRET)
	form.Add("code", code)

	// Готовим и отправляем запрос
	request, _ := http.NewRequest("POST", requestURL, strings.NewReader(form.Encode()))
	request.Header.Set("Accept", "application/json") // просим прислать ответ в формате json
	response, _ := client.Do(request)
	defer response.Body.Close()

	// Достаём данные из тела ответа
	var responsejson struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(response.Body).Decode(&responsejson)
	return responsejson.AccessToken
}

// Получаем информацию о пользователе
func getUserData(AccessToken string) UserData {
	// Создаём http-клиент с дефолтными настройками
	client := http.Client{}
	requestURL := "https://api.github.com/user"

	// Готовим и отправляем запрос
	request, _ := http.NewRequest("GET", requestURL, nil)
	request.Header.Set("Authorization", "Bearer "+AccessToken)
	response, _ := client.Do(request)
	defer response.Body.Close()

	var data UserData
	json.NewDecoder(response.Body).Decode(&data)
	return data
}
*/
