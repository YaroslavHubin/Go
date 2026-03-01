package main

import (
    "encoding/json"  
    "html/template"   
    "log"             
    "net/http"        
    "os"              
    "time"            
)

// Структура описує зарядну станцію
type ChargingStation struct {
    Name       string      `json:"name"`        // Назва станції
    Location   string      `json:"location"`    // Локація
    PowerKW    int         `json:"power_kw"`    // Потужність у кВт
    Available  bool        `json:"available"`   // Доступність (true/false)
    Devices    []Device    `json:"devices"`     // Список пристроїв (конекторів)
    Indicators []Indicator `json:"indicators"`  // Список показників (статистика)
}

// Структура "Device" описує пристрій (конектор)
type Device struct {
    Type   string `json:"type"`   // Тип конектора
    Status string `json:"status"` // Статус (вільний, зайнятий тощо)
}

// Структура "Indicator" описує показник (кількість зарядок)
type Indicator struct {
    Name  string `json:"name"`  // Назва показника
    Value int    `json:"value"` // Значення
    Unit  string `json:"unit"`  // Одиниця вимірювання
}

func main() {
    // Відкриваємо файл JSON з даними
    file, err := os.Open("data.json")
    if err != nil {
        log.Fatal(err) // Якщо файл не знайдено — завершуємо програму
    }
    defer file.Close() // Закриваємо файл після завершення роботи

    // Змінна для збереження масиву станцій
    var stations []ChargingStation

    // Декодуємо JSON у структуру Go
    if err := json.NewDecoder(file).Decode(&stations); err != nil {
        log.Fatal(err) // Якщо помилка при читанні JSON — завершуємо програму
    }

    // Парсимо HTML-шаблон
    tmpl, err := template.ParseFiles("templates/station.html")
    if err != nil {
        log.Fatal(err) // Якщо шаблон не знайдено — завершуємо програму
    }

    // Реєструємо HTTP handler для головної сторінки "/"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Використовуємо цикл for у Go для формування списку назв станцій
        var stationNames []string
        for _, st := range stations {
            stationNames = append(stationNames, st.Name)
        }

        // Формуємо структуру даних для передачі у шаблон
        data := struct {
            Stations     []ChargingStation // повний список станцій
            StationNames []string          // список лише назв станцій
            Now          string            // поточний час
        }{
            Stations:     stations,
            StationNames: stationNames,
            Now:          time.Now().Format("02-01-2006 15:04:05"), // формат часу
        }

        // Виконуємо шаблон і передаємо дані у HTML
        tmpl.Execute(w, data)
    })

    // Запускаємо сервер на порту 8080
    log.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
