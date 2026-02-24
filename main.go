package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof" // Подключаем pprof (будет доступен на /debug/pprof/)
	"os"
	"time"
)

func main() {
	// Основные эндпоинты приложения
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/heavy", heavyCPUHandler)
	http.HandleFunc("/syscall", syscallHandler)

	// Запускаем сервер (pprof автоматически доступен по /debug/pprof/)
	port := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", port)
	log.Printf("pprof доступен на http://localhost%s/debug/pprof/", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

// 1. Простой эндпоинт, отдающий привет
func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Привет, мир! 🚀\n")
}

// 2. Эндпоинт с тяжелой вычислительной нагрузкой
func heavyCPUHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Начало тяжелых вычислений...")

	// Загружаем CPU вычислением чисел Фибоначчи (рекурсивно - очень тяжело)
	result := fibonacci(45) // 45 - уже достаточно большое число для нагрузки

	log.Println("Тяжелые вычисления завершены")
	_, _ = fmt.Fprintf(w, "Результат вычислений: %d\n", result)
}

// Рекурсивное вычисление Фибоначчи (специально неоптимальное)
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

// 3. Эндпоинт, который делает syscall и ждет 2 секунды
func syscallHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Начало системного вызова...")

	// Делаем реальный системный вызов - читаем файл /dev/null
	// (это именно syscall, а не просто time.Sleep)
	file, err := os.Open("/dev/null")
	if err != nil {
		http.Error(w, "Не удалось открыть /dev/null", http.StatusInternalServerError)
		return
	}

	// Читаем 1 байт из /dev/null (это системный вызов read)
	buf := make([]byte, 1)
	_, err = file.Read(buf)
	_ = file.Close()

	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка чтения", http.StatusInternalServerError)
		return
	}

	// Теперь делаем паузу, чтобы заблокировать горутину
	// Используем time.Sleep - это тоже в итоге системный вызов nanosleep
	time.Sleep(2 * time.Second)

	log.Println("Системный вызов завершен")
	_, _ = io.WriteString(w, "Системный вызов выполнен, подождали 2 секунды\n")
}
