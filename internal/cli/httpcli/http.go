package httpcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var BaseURL = "http://localhost:8080"

// путь к файлу с cookie
func sessionFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./.bizon_session"
	}
	return filepath.Join(home, ".bizon-cli", "session")
}

// загружаем cookie из файла, если есть
func loadSessionCookie() (string, error) {
	path := sessionFilePath()
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// сохраняем cookie (например, "session=xxxxxx")
func saveSessionCookie(cookie string) error {
	path := sessionFilePath()
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0700)

	return os.WriteFile(path, []byte(cookie), 0600)
}

// DoJSON — отправка запроса + автосохранение/подстановка cookie
func DoJSON(method, path string, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, BaseURL+path, reader)
	if err != nil {
		return nil, err
	}

	// Если тело есть — указываем тип
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 1) Подставляем cookie из файла, если есть
	if cookie, err := loadSessionCookie(); err == nil && cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	// 2) Делаем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 3) Если сервер вернул Set-Cookie — сохраняем новую сессию
	setCookie := resp.Header.Get("Set-Cookie")
	if setCookie != "" {
		// Берём только ключ=значение (обрезаем "; Path=/; HttpOnly; ..." и т.п.)
		session := strings.Split(setCookie, ";")[0]
		if strings.Contains(session, "=") {
			saveSessionCookie(session)
		}
	}

	// 4) читаем тело ответа
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Красиво форматируем JSON, если это JSON
	var pretty bytes.Buffer
	if json.Indent(&pretty, raw, "", "  ") == nil {
		return pretty.Bytes(), nil
	}

	return raw, nil
}

func printJSON(b []byte) {
	fmt.Println(string(b))
}
