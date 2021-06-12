package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//HTTP statuses: https://golang.org/pkg/net/http/#pkg-constants

// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// User struct
type User struct {
	Id      int       `json:"id"`
	NewId   uuid.UUID `json:"newId"`
	Name    string    `json:"name"`
	Surname string    `json:"surname"`
	Phone   string    `json:"phone"`
	Email   string    `json:"email"`
	Address *Address  `json:"address"`
	Pets    *[]Pet    `json:"pets"`
	PetIds  []int     `json:"petIds"`
}

// Address struct
type Address struct {
	Id         int    `json:"id"`
	OwnerId    int    `json:"ownerId"`
	NewId      int    `json:"newId"`
	NewOwnerId int    `json:"newOwnerId"`
	Country    string `json:"country"`
	City       string `json:"city"`
	Street     string `json:"street"`
	House      string `json:"house"`
	Building   string `json:"building"`
	Flat       string `json:"flat"`
}

// Pet struct
type Pet struct {
	Id         int       `json:"id"`
	NewId      uuid.UUID `json:"newId"`
	OwnerId    int       `json:"ownerId"`
	NewOwnerId uuid.UUID `json:"newOwnerId"`
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	Sex        string    `json:"sex"`
	Age        int       `json:"age"` //@TODO считать самому на основе birthday
	Breed      string    `json:"breed"`
	Birthday   string    `json:"birthday"`
}

type Password struct {
	UserId   int       `json:"userId"`
	UserUuid uuid.UUID `json:"userUuid"`
	Password string    `json:"password"`
	//Hash
}

var users []User
var pets []Pet
var passwords []Password

type Service struct {
	Id            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PriceInfo     string    `json:"priceInfo"`
	IsVisible     bool      `json:"isVisible"`     //TODO: do not return false values
	OrderPosition int       `json:"orderPosition"` //TODO: return in corresponding order
}

var services = []Service{
	{
		Id:            uuid.Must(uuid.Parse("fffc5eea-5d55-4fd4-8305-33a06d99743b")),
		Name:          "Выгул",
		Description:   "Выгульщик погуляет с собакой, если вы на работе, в отъезде или заболели",
		PriceInfo:     "от 399₽/выгул",
		IsVisible:     true,
		OrderPosition: 0,
	},
	{
		Id:            uuid.Must(uuid.Parse("bf2f98e9-a924-472b-811c-6a180cad54ad")),
		Name:          "Догситтинг",
		Description:   "Ситтер поживёт с вашим питомцем, пока вы в отъезде",
		PriceInfo:     "от 699₽/сутки",
		IsVisible:     true,
		OrderPosition: 1,
	},
	{
		Id:            uuid.Must(uuid.Parse("23da6ee9-f35a-424a-828a-97f2b23c6a22")),
		Name:          "Уход за питомцем",
		Description:   "Поможем с ежедневным уходом за питомцем",
		PriceInfo:     "от 250₽/услуга",
		IsVisible:     true,
		OrderPosition: 2,
	},
	{
		Id:            uuid.Must(uuid.Parse("75451f2c-04d2-4773-9bdd-f7b9b411fb20")),
		Name:          "Дневная няня",
		Description:   "Няня посидит с питомцем у вас дома, пока вас нет",
		PriceInfo:     "от 399₽/визит",
		IsVisible:     true,
		OrderPosition: 3,
	},
}

func getServiceName(id uuid.UUID) string {
	var name = ""
	for _, item := range services {
		if item.Id == id {
			name = item.Name
			break
		}

	}
	return name
}

//Пока что излишне
/*
type UsersOrders struct {
	UserId  uuid.UUID `json:"userId"`
	OrderId uuid.UUID `json:"orderId"`
}

var usersOrders []UsersOrders
*/

//@TODO: заказ на несколько питомцев
type Order struct {
	Id           uuid.UUID   `json:"id"`
	UserId       uuid.UUID   `json:"userId"`
	PetId        uuid.UUID   `json:"petId"`
	ExecutorId   uuid.UUID   `json:"executorId"`
	ServiceId    uuid.UUID   `json:"serviceId"`
	ServiceName  string      `json:"type"`
	Status       OrderStatus `json:"status"`
	CreatedAt    string      `json:"createdAt"` //More on Time pkg: https://golang.org/pkg/time/
	UpdatedAt    string      `json:"updatedAt"`
	ScheduledFor string      `json:"scheduledFor"`
}

//Add mock orders
var orders = []Order{
	{
		Id:           uuid.Must(uuid.Parse("1a8c5856-38f0-4f7f-a689-52469bddbf9e")),
		UserId:       uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")),
		PetId:        uuid.Must(uuid.Parse("367c7a83-3f7e-4870-ace0-aabdec0fe18e")),
		ExecutorId:   uuid.Must(uuid.Parse("162f9de7-3b3f-49a1-81eb-3d65f4048b84")),
		ServiceId:    uuid.Must(uuid.Parse("bf2f98e9-a924-472b-811c-6a180cad54ad")),
		ServiceName:  getServiceName(uuid.Must(uuid.Parse("bf2f98e9-a924-472b-811c-6a180cad54ad"))),
		Status:       Pending,
		CreatedAt:    "2021-06-11T21:59:31+03:00",
		UpdatedAt:    "2021-06-11T21:59:31+03:00",
		ScheduledFor: "2019-07-15T21:59:31+03:00",
	},
	{
		Id:           uuid.Must(uuid.Parse("4a933b9e-a2e7-47a7-8bbb-e19bf162e23f")),
		UserId:       uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")),
		PetId:        uuid.Must(uuid.Parse("367c7a83-3f7e-4870-ace0-aabdec0fe18e")),
		ExecutorId:   uuid.Must(uuid.Parse("162f9de7-3b3f-49a1-81eb-3d65f4048b84")),
		ServiceId:    uuid.Must(uuid.Parse("fffc5eea-5d55-4fd4-8305-33a06d99743b")),
		ServiceName:  getServiceName(uuid.Must(uuid.Parse("fffc5eea-5d55-4fd4-8305-33a06d99743b"))),
		Status:       SearchingExecutor,
		CreatedAt:    "2021-06-11T21:59:31+03:00",
		UpdatedAt:    "2021-06-11T21:59:31+03:00",
		ScheduledFor: "2025-07-15T21:59:31+03:00",
	},
	{
		Id:           uuid.Must(uuid.Parse("02122c6a-b517-4998-8e28-39089006188f")),
		UserId:       uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")),
		PetId:        uuid.Must(uuid.Parse("367c7a83-3f7e-4870-ace0-aabdec0fe18e")),
		ExecutorId:   uuid.Must(uuid.Parse("162f9de7-3b3f-49a1-81eb-3d65f4048b84")),
		ServiceId:    uuid.Must(uuid.Parse("75451f2c-04d2-4773-9bdd-f7b9b411fb20")),
		ServiceName:  getServiceName(uuid.Must(uuid.Parse("75451f2c-04d2-4773-9bdd-f7b9b411fb20"))),
		Status:       Completed,
		CreatedAt:    "2021-06-11T21:59:31+03:00",
		UpdatedAt:    "2021-06-11T21:59:31+03:00",
		ScheduledFor: "2020-02-26T21:59:31+03:00",
	},
	{
		Id:           uuid.Must(uuid.Parse("4f04bd63-01f7-4671-ac74-dcd83ec3ed43")),
		UserId:       uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")),
		PetId:        uuid.Must(uuid.Parse("c3d25bae-bd8a-4a61-83ea-59084ec332a8")),
		ExecutorId:   uuid.Must(uuid.Parse("162f9de7-3b3f-49a1-81eb-3d65f4048b84")),
		ServiceId:    uuid.Must(uuid.Parse("75451f2c-04d2-4773-9bdd-f7b9b411fb20")),
		ServiceName:  getServiceName(uuid.Must(uuid.Parse("75451f2c-04d2-4773-9bdd-f7b9b411fb20"))),
		Status:       Cancelled,
		CreatedAt:    "2021-06-11T21:59:31+03:00",
		UpdatedAt:    "2021-06-11T21:59:31+03:00",
		ScheduledFor: "2021-01-11T21:59:31+03:00",
	},
}

//Enum implementation: https://levelup.gitconnected.com/implementing-enums-in-golang-9537c433d6e2
type OrderStatus int

// Declare related constants
const (
	Pending OrderStatus = iota + 1 // 0 is reserved for default encoding output
	SearchingExecutor
	Scheduled
	InProgress
	Completed
	Cancelled
)

// String - Creating common behavior - give the type a String function
func (os OrderStatus) String() string {
	return [...]string{"Pending", "SearchingExecutor", "Scheduled", "InProgress", "Completed", "Cancelled"}[os-1]
}

// EnumIndex - Creating common behavior - give the type a EnumIndex function
func (os OrderStatus) EnumIndex() int {
	return int(os)
}

/*
func main() {
	var weekday = Sunday
	fmt.Println(weekday)             // Print : Sunday
	fmt.Println(weekday.String())    // Print : Sunday
	fmt.Println(weekday.EnumIndex()) // Print : 1
}
*/

//var addresses []Address

/*
var books []Book

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	//w.Header().Set("Access-Control-Allow-Credentials", "true") // Required for cookies, authorization headers with HTTPS
	json.NewEncoder(w).Encode(books)
}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Add new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}
*/

// Get all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	json.NewEncoder(w).Encode(users)
}

// Get single user
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params
	for _, user := range users {
		if strconv.Itoa(user.Id) == params["id"] {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

// Get single user
func getUserNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params
	for _, user := range users {
		if user.NewId == uuid.Must(uuid.Parse(params["id"])) {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

// Add new user
/*
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var user User

	_ = json.NewDecoder(r.Body).Decode(&user)
	user.Id = rand.Intn(100000000) // Mock ID - not safe //@TODO: Change for GUID
	users = append(users, user)
	json.NewEncoder(w).Encode(user.Id)
}
*/

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	type Tmp struct {
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var tmp Tmp
	_ = json.NewDecoder(r.Body).Decode(&tmp)

	var user User
	user.Name = tmp.Name
	user.Phone = tmp.Phone
	user.Email = tmp.Email
	user.Id = rand.Intn(100000000) // Mock ID - not safe //@TODO: Change for GUID
	user.NewId = uuid.New()

	users = append(users, user)
	passwords = append(passwords, Password{UserId: user.Id, Password: tmp.Password})
	prettyPrintSomething()

	json.NewEncoder(w).Encode(user.Id)
}

//@TODO: проверка на наличие мобильного телефона в базе
func createUserNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	type Tmp struct {
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var tmp Tmp
	_ = json.NewDecoder(r.Body).Decode(&tmp)

	var user User
	user.Name = tmp.Name
	user.Phone = tmp.Phone
	user.Email = tmp.Email
	user.Id = rand.Intn(100000000) // Mock ID - not safe //@TODO: Change for GUID
	user.NewId = uuid.New()

	users = append(users, user)
	passwords = append(passwords, Password{UserId: user.Id, Password: tmp.Password})
	prettyPrintSomething()

	json.NewEncoder(w).Encode(user.NewId)
}

// Update user
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // Required for cookies, authorization headers with HTTPS

	params := mux.Vars(r)
	for index, item := range users {
		if strconv.Itoa(item.Id) == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)

			var err error
			user.Id, err = strconv.Atoi(params["id"])
			if err == nil {
				users = append(users, user)
				json.NewEncoder(w).Encode(user)
			}
			return
		}
	}
}

func updateUserNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, OPTIONS, DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	//w.Header().Set("Access-Control-Allow-Credentials", "true") // Required for cookies, authorization headers with HTTPS

	params := mux.Vars(r)
	for index, item := range users {
		if item.NewId == uuid.Must(uuid.Parse(params["id"])) {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)

			user.Id = 666 //убрать потом!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!1
			user.NewId = uuid.Must(uuid.Parse(params["id"]))
			users = append(users, user)
			json.NewEncoder(w).Encode(user)
			return
		}
	}
}

// Delete user
//Do not actually delete, set flag to deleted=true
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r)
	for index, item := range users {
		if strconv.Itoa(item.Id) == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	//json.NewEncoder(w).Encode(users)
	http.StatusText(http.StatusOK)

}

// Get all pets
func getPets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	json.NewEncoder(w).Encode(pets)
}

// Get all user's pets

/*
func getUserPets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params

	//потом надо животных отдельно просто хранить, а в структуре юзера - указатель
	for _, user := range users {
		if strconv.Itoa(user.Id) == params["userId"] {
			json.NewEncoder(w).Encode(user.Pets)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}
*/

func getUserPets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params

	var _pets []Pet

	for _, pet := range pets {
		if strconv.Itoa(pet.OwnerId) == params["userId"] {
			_pets = append(_pets, pet)
		}
	}
	json.NewEncoder(w).Encode(_pets)
}

func getUserPetsNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params

	var _pets []Pet

	for _, pet := range pets {
		if pet.NewOwnerId == uuid.Must(uuid.Parse(params["userId"])) {
			_pets = append(_pets, pet)
		}
	}
	json.NewEncoder(w).Encode(_pets)
}

/*
func getPet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params

	for _, user := range users {
		for _, pet := range *user.Pets {
			if strconv.Itoa(pet.Id) == params["id"] {
				json.NewEncoder(w).Encode(pet)
				return
			}
		}
	}
	json.NewEncoder(w).Encode(&User{})
}
*/

func getPet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params

	for _, pet := range pets {
		if strconv.Itoa(pet.Id) == params["id"] {
			json.NewEncoder(w).Encode(pet)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

func updatePet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // Required for cookies, authorization headers with HTTPS

	params := mux.Vars(r)
	for index, pet := range pets {
		if strconv.Itoa(pet.Id) == params["id"] {
			pets = append(pets[:index], pets[index+1:]...)
			var pet Pet
			_ = json.NewDecoder(r.Body).Decode(&pet)

			var err error
			pet.Id, err = strconv.Atoi(params["id"])
			if err == nil {
				pets = append(pets, pet)
				json.NewEncoder(w).Encode(pet)
			}
			return
		}
	}
}

// Add new pet
func createPet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var pet Pet
	_ = json.NewDecoder(r.Body).Decode(&pet)
	pet.Id = rand.Intn(100000000) // Mock ID - not safe //@TODO: Change for GUID
	pet.NewId = uuid.New()

	pets = append(pets, pet)
	json.NewEncoder(w).Encode(pet)
}

// Delete pet
//Do not actually delete, set flag to deleted=true
func deletePet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r)
	for index, item := range pets {
		if strconv.Itoa(item.Id) == params["id"] {
			pets = append(pets[:index], pets[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(pets) //@подумать надо ли это возвращать
}

/*
func authorizeUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	params := mux.Vars(r) // Gets params

	for _, password := range passwords {
		if strconv.Itoa(password.UserId) == params["login"] {
			if password.Password == params["password"] {
				json.NewEncoder(w).Encode(password.UserId)
				return
			} else {
				break
			}

		}
	}
	json.NewEncoder(w).Encode(nil)
}
*/

func authorizeUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	params := mux.Vars(r) // Gets params

	for _, user := range users {
		if user.Phone == params["login"] {
			for _, password := range passwords {
				if password.UserId == user.Id {
					if password.Password == params["password"] {
						json.NewEncoder(w).Encode(user.Id)
						return
					} else {
						break
					}

				}
			}
		}
	}
	http.Error(w, "Authorization unsuccessful", http.StatusForbidden)
	// or using the default message error
	//http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)

	//json.NewEncoder(w).Encode(nil)
}

func authorizeUserNew(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	params := mux.Vars(r) // Gets params

	for _, user := range users {
		if user.Phone == params["login"] {
			for _, password := range passwords {
				if password.UserUuid == user.NewId {
					if password.Password == params["password"] { //@TODO: hash!
						json.NewEncoder(w).Encode(user.NewId)
						return
					} else {
						break
					}

				}
			}
		}
	}
	http.Error(w, "Authorization unsuccessful", http.StatusForbidden)
}

// Get all services
func getServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	json.NewEncoder(w).Encode(services)
}

// Get single service
func getService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params

	var err error
	var id uuid.UUID
	id, err = uuid.Parse(params["id"])

	if err != nil {
		http.Error(w, "Wrong parameter format", http.StatusBadRequest)
		return
	} else {
		for _, service := range services {
			if service.Id == id {
				json.NewEncoder(w).Encode(service)
				return
			}
		}
		json.NewEncoder(w).Encode(&Service{})
	}
}

// Get all orders
//TODO: offset
func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	json.NewEncoder(w).Encode(orders)
}

// Get single order
func getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work

	params := mux.Vars(r)

	var err error
	var id uuid.UUID
	id, err = uuid.Parse(params["id"])

	if err != nil {
		http.Error(w, "Wrong parameter format", http.StatusBadRequest)
		return
	}

	for _, order := range orders {
		if order.Id == id {
			json.NewEncoder(w).Encode(order)
			return
		}
	}
	json.NewEncoder(w).Encode(&Order{}) //TODO: return {"success": true, "data": null, "errors": null}
}

// Get user's orders
//TODO: offset

func getUserOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params

	var err error
	var userId uuid.UUID
	userId, err = uuid.Parse(params["userId"])

	if err != nil {
		http.Error(w, "Wrong parameter format", http.StatusBadRequest)
		return
	}
	var thisUserOrders []Order
	for _, order := range orders {
		if order.UserId == userId {
			thisUserOrders = append(thisUserOrders, order)
		}
	}
	json.NewEncoder(w).Encode(thisUserOrders)
}

//returns sorted slice (by ScheduledFor)
func getUserOrdersSorted(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r) // Gets params

	var err error
	var userId uuid.UUID
	userId, err = uuid.Parse(params["userId"])

	if err != nil {
		http.Error(w, "Wrong parameter format", http.StatusBadRequest)
		return
	}
	var thisUserOrders []Order
	for _, order := range orders {
		if order.UserId == userId {
			thisUserOrders = append(thisUserOrders, order)
		}
	}

	sort.Slice(thisUserOrders, func(i, j int) bool {
		return thisUserOrders[i].ScheduledFor > thisUserOrders[j].ScheduledFor
	})

	json.NewEncoder(w).Encode(thisUserOrders)
}

//TODO: проверять инпуты "userId", "petId", "serviceId", "scheduledFor" на валидность
func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var order Order

	_ = json.NewDecoder(r.Body).Decode(&order)

	order.Id = uuid.New()
	order.CreatedAt = time.Now().Format(time.RFC3339)
	order.UpdatedAt = order.CreatedAt
	order.Status = Pending
	order.ServiceName = getServiceName(order.ServiceId)

	orders = append(orders, order)

	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "PUT, POST, GET, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // Required for cookies, authorization headers with HTTPS

	params := mux.Vars(r)

	var err error
	var id uuid.UUID
	id, err = uuid.Parse(params["id"])

	if err != nil {
		http.Error(w, "Wrong parameter format", http.StatusBadRequest)
		return
	}

	for index, item := range orders {
		if item.Id == id {
			orders = append(orders[:index], orders[index+1:]...)

			var tmp Order
			_ = json.NewDecoder(r.Body).Decode(&tmp)

			var updatedOrder Order = item
			updatedOrder.UpdatedAt = time.Now().Format(time.RFC3339)
			//Leave the status AS IS unless it was explicitly set in request body
			if tmp.Status > 0 {
				updatedOrder.Status = tmp.Status
			}
			//TODO: check if input is int
			if tmp.ScheduledFor != "" {
				_, err = time.Parse(time.RFC3339, tmp.ScheduledFor)
				if err == nil {
					updatedOrder.ScheduledFor = tmp.ScheduledFor
				}
			}
			orders = append(orders, updatedOrder)
			json.NewEncoder(w).Encode(updatedOrder)
			return
		}
	}
}

//Do not actually delete, set flag to deleted=true
func deleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	params := mux.Vars(r)

	var err error
	var id uuid.UUID
	id, err = uuid.Parse(params["id"])

	if err != nil {
		http.Error(w, "Wrong parameter format", http.StatusBadRequest)
		return
	}
	for index, item := range orders {
		if item.Id == id {
			orders = append(orders[:index], orders[index+1:]...)
			break
		}
	}
	http.StatusText(http.StatusOK)
}

//Helper function to print json to console
func prettyPrintSomething() {
	prettyJSON, err := json.MarshalIndent(passwords, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", string(prettyJSON))
}

//Helper function to

// Main function
func main() {
	// Init router
	r := mux.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Hardcoded data - @todo: add database
	//books = append(books, Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	//books = append(books, Book{ID: "2", Isbn: "454555", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})

	// Hardcoded data
	//var address_user0 Address
	//var address_user1 Address
	var pets_user0 []Pet
	var pets_user1 []Pet
	var pets_user2 []Pet

	address_user0 := Address{Id: 0, OwnerId: 0, Country: "Россия", City: "Москва", Street: "Петушиная", House: "69", Flat: "420"}
	address_user2 := Address{Id: 1, OwnerId: 2, Country: "Russia", City: "Moscow", Street: "Tverskaya", House: "420", Flat: "69"}

	pets_user0 = append(pets_user0, Pet{Id: 0, OwnerId: 0, NewId: uuid.Must(uuid.Parse("367c7a83-3f7e-4870-ace0-aabdec0fe18e")), NewOwnerId: uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")), Type: "cat", Name: "Владимир", Sex: "male", Age: 1})
	pets_user0 = append(pets_user0, Pet{Id: 1, OwnerId: 0, NewId: uuid.Must(uuid.Parse("c3d25bae-bd8a-4a61-83ea-59084ec332a8")), NewOwnerId: uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")), Type: "dog", Name: "Джереми", Sex: "male", Age: 4, Breed: "Алабай", Birthday: "2010-04-10"})
	pets_user1 = append(pets_user1, Pet{Id: 2, OwnerId: 1, NewId: uuid.Must(uuid.Parse("7e58091e-3736-49ab-b502-7cb3314bbe21")), NewOwnerId: uuid.Must(uuid.Parse("4b69a783-65d4-4bba-adc7-8935f22c1fc6")), Type: "crocodile", Name: "Антон", Sex: "male", Age: 35})
	pets_user2 = append(pets_user2, Pet{Id: 3, OwnerId: 2, NewId: uuid.Must(uuid.Parse("6f8cb018-1afd-412b-86c0-a83d3b3c47bd")), NewOwnerId: uuid.Must(uuid.Parse("52a226aa-9ee1-4ba1-a053-c67eeff55366")), Type: "indricotherium", Name: "Musinit", Sex: "male", Age: 10})

	pets = append(pets, pets_user0[0], pets_user0[1], pets_user1[0], pets_user2[0])

	users = append(users, User{Id: 0, NewId: uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")), Name: "Ицхак", Surname: "Пинтосевич", Phone: "79000000000", Email: "test@mail.ru", Address: &address_user0, Pets: &pets_user0})
	users = append(users, User{Id: 1, NewId: uuid.Must(uuid.Parse("4b69a783-65d4-4bba-adc7-8935f22c1fc6")), Name: "Александр", Surname: "Тестовый", Phone: "79111111111", Pets: &pets_user1})
	users = append(users, User{Id: 2, NewId: uuid.Must(uuid.Parse("52a226aa-9ee1-4ba1-a053-c67eeff55366")), Name: "Oleg", Surname: "Musin", Phone: "79222222222", Address: &address_user2, Pets: &pets_user2})

	passwords = append(
		passwords,
		Password{UserId: 0, UserUuid: uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")), Password: "zero"},
		Password{UserId: 1, UserUuid: uuid.Must(uuid.Parse("4b69a783-65d4-4bba-adc7-8935f22c1fc6")), Password: "one"},
		Password{UserId: 2, UserUuid: uuid.Must(uuid.Parse("52a226aa-9ee1-4ba1-a053-c67eeff55366")), Password: "two"})

	/*
		usersOrders = append(
			usersOrders,
			UsersOrders{UserId: uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")), OrderId: uuid.Must(uuid.Parse("fffc5eea-5d55-4fd4-8305-33a06d99743b"))},
			UsersOrders{UserId: uuid.Must(uuid.Parse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")), OrderId: uuid.Must(uuid.Parse("bf2f98e9-a924-472b-811c-6a180cad54ad"))})
	*/

	// Route handles & endpoints
	/*
		r.HandleFunc("/books", getBooks).Methods("GET")
		r.HandleFunc("/books/{id}", getBook).Methods("GET")
		r.HandleFunc("/books", createBook).Methods("POST")
		r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
		r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	*/

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users/new/{id}", getUserNew).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/new", createUserNew).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/new/{id}", updateUserNew).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	r.HandleFunc("/pets/user/{userId}", getUserPets).Methods("GET")
	r.HandleFunc("/pets/user/new/{userId}", getUserPetsNew).Methods("GET")
	r.HandleFunc("/pets/{id}", getPet).Methods("GET")
	r.HandleFunc("/pets", getPets).Methods("GET")
	r.HandleFunc("/pets/{id}", updatePet).Methods("PUT")
	r.HandleFunc("/pets", createPet).Methods("POST")
	r.HandleFunc("/pets/{id}", deletePet).Methods("DELETE")

	r.HandleFunc("/services", getServices).Methods("GET")
	r.HandleFunc("/services/{id}", getService).Methods("GET")

	r.HandleFunc("/orders/user/{userId}", getUserOrders).Methods("GET")
	r.HandleFunc("/orders/user/sorted/{userId}", getUserOrdersSorted).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders/{id}", updateOrder).Methods("PUT")
	r.HandleFunc("/orders/{id}", deleteOrder).Methods("DELETE")

	r.HandleFunc("/auth/", authorizeUser).
		Queries("login", "{login}", "password", "{password}").Methods("GET")
	r.HandleFunc("/auth/new", authorizeUserNew).
		Queries("login", "{login}", "password", "{password}").Methods("GET")

	//Print port info
	//fmt.Printf("sobaken-vigulyaken starting on port: %s...", port)

	// Print Json with indents, the pretty way:
	/*
		prettyJSON, err := json.MarshalIndent(users, "", "    ")
		if err != nil {
			log.Fatal("Failed to generate json", err)
		}
		fmt.Printf("%s\n", string(prettyJSON))
	*/

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}

// Request sample
// {
// 	"isbn":"4545454",
// 	"title":"Book Three",
// 	"author":{"firstname":"Harry","lastname":"White"}
// }

//Request sample: POST/users
/*
{
    "name": "Misha",
    "surname": "Smolin",
    "phone": "+79162712542",
    "email": "kokus@mail.ru",
    "address": null,
    "pets": null
}
*/
