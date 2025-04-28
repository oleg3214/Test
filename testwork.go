package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	// 1. Путь к chromedriver
	const chromeDriverPath = "C:/goProject/Test/chromedriver.exe"

	// 2. Запуск ChromeDriver сервиса
	service, err := selenium.NewChromeDriverService(chromeDriverPath, 4444)
	if err != nil {
		log.Fatal("Ошибка запуска ChromeDriver:", err)
	}
	defer service.Stop()
	userDataDir := filepath.Join(os.TempDir(), "chrome_profile")
	fmt.Println(os.TempDir())
	// defer os.RemoveAll(userDataDir)
	// 3. Настройка параметров Chrome
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{
		Path: "C:/Program Files/Google/Chrome/Application/chrome.exe",
		Args: []string{
			fmt.Sprintf("--user-data-dir=%s", userDataDir),
			"--start-maximized",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--lang=en-US", // Устанавливаем английский язык для стабильности
		},
	})

	// 4. Подключение к веб-драйверу
	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Ошибка подключения к WebDriver:", err)
	}
	defer driver.Quit()
	driver.ExecuteScript(`
    window.navigator.chrome = {
        runtime: {},
    };`, nil)

	// Удаляем следы автоматизации
	driver.ExecuteScript(`
    Object.defineProperty(navigator, 'plugins', {
        get: () => [1, 2, 3],
    });`, nil)
	// 5. Открытие страницы входа
	err = driver.Get("https://www.threads.net/login")
	if err != nil {
		log.Fatal("Ошибка загрузки страницы:", err)
	}

	// 6. Ожидание полной загрузки страницы
	err = driver.WaitWithTimeout(func(wd selenium.WebDriver) (bool, error) {
		_, err := wd.FindElement(selenium.ByCSSSelector, "input[autocomplete='username']")
		return err == nil, nil
	}, 20*time.Second)
	if err != nil {
		log.Fatal("Поля ввода не найдены:", err)
	}

	// 7. Ввод логина и пароля
	username := "oleg.ivanchenko" // Замените на реальный логин
	password := "Bart009912"      // Замените на реальный пароль

	// Находим поле ввода логина
	usernameField, err := driver.FindElement(selenium.ByCSSSelector, "input[autocomplete='username']")
	if err != nil {
		log.Fatal("Не найдено поле логина:", err)
	}
	time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
	usernameField.SendKeys(username)

	// Находим поле ввода пароля
	passwordField, err := driver.FindElement(selenium.ByCSSSelector, "input[autocomplete='current-password']")
	if err != nil {
		log.Fatal("Не найдено поле пароля:", err)
	}
	time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
	passwordField.SendKeys(password)

	// 8. Нажимаем кнопку входа
	loginButton, err := driver.FindElement(selenium.ByXPATH, "//div[normalize-space()='Log in']")
	if err != nil {
		log.Fatal("Не найдена кнопка входа:", err)
	}
	time.Sleep(time.Duration(1+rand.Intn(3)) * time.Second)
	loginButton.Click()

	fmt.Println("Успешный вход в систему!")

	// 10. Ожидание ввода пользователя перед закрытием
	fmt.Println("Нажмите Enter чтобы закрыть браузер...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
