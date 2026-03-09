package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
)

// ===== Структури =====
// Кожна структура описує певний об'єкт енергосистеми.
// JSON-теги визначають, як поля будуть виглядати у відповіді сервера.

// Зарядна станція
type ChargingStation struct {
    ID       int    `json:"id"`        // Унікальний ідентифікатор станції
    Location string `json:"location"`  // Місце розташування
    PowerKW  int    `json:"power_kw"`  // Потужність у кіловатах
    Status   string `json:"status"`    // Статус: available, busy, offline
}

// Генератор
type Generator struct {
    ID     int    `json:"id"`        // Ідентифікатор генератора
    Power  int    `json:"power_kw"`  // Потужність генератора
    Status string `json:"status"`    // Статус: active, inactive
}

// Мережа
type Network struct {
    ID      int    `json:"id"`       // Ідентифікатор мережі
    Voltage int    `json:"voltage"`  // Напруга у вольтах
    Status  string `json:"status"`   // Статус: stable, unstable
}

// Споживач
type Consumer struct {
    ID       int    `json:"id"`        // Ідентифікатор споживача
    Name     string `json:"name"`      // Назва споживача
    DemandKW int    `json:"demand_kw"` // Потреба у кіловатах
}

// Датчик
type Sensor struct {
    ID    int    `json:"id"`     // Ідентифікатор датчика
    Type  string `json:"type"`   // Тип датчика (наприклад, temperature)
    Value string `json:"value"`  // Значення (наприклад, 25C)
}

// ===== Імітація бази даних =====
// Тут ми створюємо "базу даних" у вигляді зрізів (slice).
// В реальних умовах це були б записи у справжній БД.

var stations = []ChargingStation{
    {ID: 1, Location: "Kyiv Center", PowerKW: 50, Status: "available"},
    {ID: 2, Location: "Kyiv Airport", PowerKW: 100, Status: "busy"},
}
var generators = []Generator{{ID: 1, Power: 500, Status: "active"}}
var networks = []Network{{ID: 1, Voltage: 220, Status: "stable"}}
var consumers = []Consumer{{ID: 1, Name: "Office Building", DemandKW: 200}}
var sensors = []Sensor{{ID: 1, Type: "temperature", Value: "25C"}}

// ===== Generic Helpers =====
// Допоміжні функції, які використовуються для скорочення коду.

// Функція для відправки JSON-відповіді клієнту.
func getJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json") // Встановлюємо заголовок відповіді
    json.NewEncoder(w).Encode(data)                   // Кодуємо дані у JSON і відправляємо
}

// Узагальнена функція для видалення об'єкта за ID.
// Використовуємо дженерики (T any), щоб працювати з будь-яким типом.
func deleteByID[T any](slice []T, id int, getID func(T) int) ([]T, bool) {
    for i, item := range slice {
        if getID(item) == id { // Якщо ID співпадає
            return append(slice[:i], slice[i+1:]...), true // Видаляємо елемент
        }
    }
    return slice, false // Якщо не знайдено — повертаємо false
}

// ===== Handlers =====
// Тут ми описуємо обробники HTTP-запитів для кожного об'єкта.

// === Stations ===
func getStations(w http.ResponseWriter, r *http.Request) { getJSON(w, stations) } // GET: отримати всі станції

func createStation(w http.ResponseWriter, r *http.Request) { // POST: створити нову станцію
    var obj ChargingStation
    if err := json.NewDecoder(r.Body).Decode(&obj); err != nil { // Декодуємо JSON із запиту
        http.Error(w, "Invalid input", http.StatusBadRequest)    // Якщо помилка — повертаємо 400
        return
    }
    obj.ID = len(stations) + 1 // Присвоюємо новий ID
    stations = append(stations, obj) // Додаємо у "базу"
    w.WriteHeader(http.StatusCreated) // Відповідь 201 Created
    json.NewEncoder(w).Encode(obj)    // Повертаємо створений об'єкт
}

func deleteStation(w http.ResponseWriter, r *http.Request) { // DELETE: видалити станцію
    id, _ := strconv.Atoi(r.URL.Query().Get("id")) // Отримуємо id з параметрів URL
    var ok bool
    stations, ok = deleteByID(stations, id, func(s ChargingStation) int { return s.ID })
    if ok {
        w.WriteHeader(http.StatusNoContent) // Якщо видалено — 204 No Content
    } else {
        http.Error(w, "Station not found", http.StatusNotFound) // Якщо не знайдено — 404
    }
}

// === Generators ===
// Обробники для роботи з генераторами

// GET: отримати всі генератори
func getGenerators(w http.ResponseWriter, r *http.Request) { 
    getJSON(w, generators) // Використовуємо допоміжну функцію getJSON для повернення списку генераторів у форматі JSON
}

// POST: створити новий генератор
func createGenerator(w http.ResponseWriter, r *http.Request) {
    var obj Generator
    // Декодуємо JSON із тіла запиту у структуру Generator
    if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest) // Якщо дані некоректні — повертаємо 400 Bad Request
        return
    }
    obj.ID = len(generators) + 1 // Присвоюємо новий ID (на основі довжини зрізу)
    generators = append(generators, obj) // Додаємо генератор у "базу даних"
    w.WriteHeader(http.StatusCreated)    // Встановлюємо статус 201 Created
    json.NewEncoder(w).Encode(obj)       // Повертаємо створений генератор у відповіді
}

// DELETE: видалити генератор за ID
func deleteGenerator(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Query().Get("id")) // Отримуємо параметр id з URL і перетворюємо у число
    var ok bool
    // Використовуємо узагальнену функцію deleteByID для видалення генератора зі зрізу
    generators, ok = deleteByID(generators, id, func(g Generator) int { return g.ID })
    if ok {
        w.WriteHeader(http.StatusNoContent) // Якщо видалено — повертаємо 204 No Content
    } else {
        http.Error(w, "Generator not found", http.StatusNotFound) // Якщо не знайдено — повертаємо 404 Not Found
    }
}

// === Networks ===
// Обробники для роботи з мережами

// GET: отримати всі мережі
func getNetworks(w http.ResponseWriter, r *http.Request) { 
    getJSON(w, networks) // Використовуємо допоміжну функцію getJSON для повернення списку мереж у форматі JSON
}

// POST: створити нову мережу
func createNetwork(w http.ResponseWriter, r *http.Request) {
    var obj Network
    // Декодуємо JSON із тіла запиту у структуру Network
    if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest) // Якщо дані некоректні — повертаємо 400 Bad Request
        return
    }
    obj.ID = len(networks) + 1 // Присвоюємо новий ID
    networks = append(networks, obj) // Додаємо мережу у "базу даних"
    w.WriteHeader(http.StatusCreated) // Встановлюємо статус 201 Created
    json.NewEncoder(w).Encode(obj)    // Повертаємо створену мережу у відповіді
}

// DELETE: видалити мережу за ID
func deleteNetwork(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Query().Get("id")) // Отримуємо параметр id з URL
    var ok bool
    // Використовуємо узагальнену функцію deleteByID для видалення мережі зі зрізу
    networks, ok = deleteByID(networks, id, func(n Network) int { return n.ID })
    if ok {
        w.WriteHeader(http.StatusNoContent) // Якщо видалено — повертаємо 204 No Content
    } else {
        http.Error(w, "Network not found", http.StatusNotFound) // Якщо не знайдено — повертаємо 404 Not Found
    }
}

// === Consumers ===
// Обробники для роботи зі споживачами

// GET: отримати всіх споживачів
func getConsumers(w http.ResponseWriter, r *http.Request) { 
    getJSON(w, consumers) // Використовуємо допоміжну функцію getJSON для повернення списку споживачів у форматі JSON
}

// POST: створити нового споживача
func createConsumer(w http.ResponseWriter, r *http.Request) {
    var obj Consumer
    // Декодуємо JSON із тіла запиту у структуру Consumer
    if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest) // Якщо дані некоректні — повертаємо 400 Bad Request
        return
    }
    obj.ID = len(consumers) + 1 // Присвоюємо новий ID
    consumers = append(consumers, obj) // Додаємо споживача у "базу даних"
    w.WriteHeader(http.StatusCreated)  // Встановлюємо статус 201 Created
    json.NewEncoder(w).Encode(obj)     // Повертаємо створеного споживача у відповіді
}

// DELETE: видалити споживача за ID
func deleteConsumer(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Query().Get("id")) // Отримуємо параметр id з URL
    var ok bool
    // Використовуємо узагальнену функцію deleteByID для видалення споживача зі зрізу
    consumers, ok = deleteByID(consumers, id, func(c Consumer) int { return c.ID })
    if ok {
        w.WriteHeader(http.StatusNoContent) // Якщо видалено — повертаємо 204 No Content
    } else {
        http.Error(w, "Consumer not found", http.StatusNotFound) // Якщо не знайдено — повертаємо 404 Not Found
    }
}


// === Sensors ===
// Обробники для роботи з датчиками

// GET: отримати всі датчики
func getSensors(w http.ResponseWriter, r *http.Request) { 
    getJSON(w, sensors) // Використовуємо допоміжну функцію getJSON для повернення списку датчиків у форматі JSON
}

// POST: створити новий датчик
func createSensor(w http.ResponseWriter, r *http.Request) {
    var obj Sensor
    // Декодуємо JSON із тіла запиту у структуру Sensor
    if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest) // Якщо дані некоректні — повертаємо 400 Bad Request
        return
    }
    obj.ID = len(sensors) + 1 // Присвоюємо новий ID (на основі довжини зрізу)
    sensors = append(sensors, obj) // Додаємо датчик у "базу даних"
    w.WriteHeader(http.StatusCreated) // Встановлюємо статус 201 Created
    json.NewEncoder(w).Encode(obj)    // Повертаємо створений датчик у відповіді
}

// DELETE: видалити датчик за ID
func deleteSensor(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Query().Get("id")) // Отримуємо параметр id з URL і перетворюємо у число
    var ok bool
    // Використовуємо узагальнену функцію deleteByID для видалення датчика зі зрізу
    sensors, ok = deleteByID(sensors, id, func(s Sensor) int { return s.ID })
    if ok {
        w.WriteHeader(http.StatusNoContent) // Якщо видалено — повертаємо 204 No Content
    } else {
        http.Error(w, "Sensor not found", http.StatusNotFound) // Якщо не знайдено — повертаємо 404 Not Found
    }
}

// ===== main =====
// Головна функція програми, яка запускає HTTP-сервер і реєструє всі маршрути
func main() {
    // Реєструємо маршрути для станцій
    http.HandleFunc("/stations", getStations)         // GET: отримати всі станції
    http.HandleFunc("/stations/create", createStation) // POST: створити нову станцію
    http.HandleFunc("/stations/delete", deleteStation) // DELETE: видалити станцію

    // Реєструємо маршрути для генераторів
    http.HandleFunc("/generators", getGenerators)         // GET: отримати всі генератори
    http.HandleFunc("/generators/create", createGenerator) // POST: створити генератор
    http.HandleFunc("/generators/delete", deleteGenerator) // DELETE: видалити генератор

    // Реєструємо маршрути для мереж
    http.HandleFunc("/networks", getNetworks)         // GET: отримати всі мережі
    http.HandleFunc("/networks/create", createNetwork) // POST: створити мережу
    http.HandleFunc("/networks/delete", deleteNetwork) // DELETE: видалити мережу

    // Реєструємо маршрути для споживачів
    http.HandleFunc("/consumers", getConsumers)         // GET: отримати всіх споживачів
    http.HandleFunc("/consumers/create", createConsumer) // POST: створити споживача
    http.HandleFunc("/consumers/delete", deleteConsumer) // DELETE: видалити споживача

    // Реєструємо маршрути для датчиків
    http.HandleFunc("/sensors", getSensors)         // GET: отримати всі датчики
    http.HandleFunc("/sensors/create", createSensor) // POST: створити датчик
    http.HandleFunc("/sensors/delete", deleteSensor) // DELETE: видалити датчик

    // Swagger UI
    // Цей обробник дозволяє відкривати документацію Swagger за адресою /docs/
    // Він віддає статичні файли (index.html, openapi.yaml) з поточної директорії
    http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("."))))

    // Лог повідомлення про запуск сервера
    log.Println("Server running on http://localhost:8080")

    // Запускаємо сервер на порту 8080
    // Якщо виникне помилка — програма завершиться з логом
    log.Fatal(http.ListenAndServe(":8080", nil))
}
