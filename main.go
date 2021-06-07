package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

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
	Id      int      `json:"id"`
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Phone   string   `json:"phone"`
	Email   string   `json:"email"`
	Address *Address `json:"address"`
	Pets    *[]Pet   `json:"pets"`
	PetIds  []int    `json:"petIds"`
}

// Address struct
type Address struct {
	Id       int    `json:"id"`
	OwnerId  int    `json:"ownerId"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Street   string `json:"street"`
	House    string `json:"house"`
	Building string `json:"building"`
	Flat     string `json:"flat"`
}

// Pet struct
type Pet struct {
	Id      int    `json:"id"`
	OwnerId int    `json:"ownerId"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Sex     string `json:"sex"`
	Age     int    `json:"age"`
}

// Init books var as a slice Book struct
var books []Book

// Init users var as a slice User struct
var users []User

// Init users var as a slice Pet struct
var pets []Pet

//var addresses []Address
//var pets []Pet

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
	for _, item := range users {
		if strconv.Itoa(item.Id) == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

// Add new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // Required for CORS support to work
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.Id = rand.Intn(100000000) // Mock ID - not safe //@TODO: Change for GUID
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
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

// Delete user
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
	json.NewEncoder(w).Encode(users)
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
	pets = append(pets, pet)
	json.NewEncoder(w).Encode(pet)
}

// Delete pet
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

// Main function
func main() {
	// Init router
	r := mux.NewRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Hardcoded data - @todo: add database
	books = append(books, Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "454555", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})

	// Hardcoded data
	//var address_user0 Address
	//var address_user1 Address
	var pets_user0 []Pet
	var pets_user1 []Pet
	var pets_user2 []Pet

	address_user0 := Address{Id: 0, OwnerId: 0, Country: "Россия", City: "Москва", Street: "Петушиная", House: "69", Flat: "420"}
	address_user2 := Address{Id: 1, OwnerId: 2, Country: "Russia", City: "Moscow", Street: "Tverskaya", House: "420", Flat: "69"}

	pets_user0 = append(pets_user0, Pet{Id: 0, OwnerId: 0, Type: "cat", Name: "Владимир", Sex: "male", Age: 1})
	pets_user0 = append(pets_user0, Pet{Id: 1, OwnerId: 0, Type: "dog", Name: "Джереми", Sex: "male", Age: 4})
	pets_user1 = append(pets_user1, Pet{Id: 2, OwnerId: 1, Type: "crocodile", Name: "Антон", Sex: "male", Age: 35})
	pets_user2 = append(pets_user2, Pet{Id: 3, OwnerId: 2, Type: "indricotherium", Name: "Musinit", Sex: "male", Age: 10})

	pets = append(pets, pets_user0[0], pets_user0[1], pets_user1[0], pets_user2[0])

	users = append(users, User{Id: 0, Name: "Ицхак", Surname: "Пинтосевич", Phone: "+79123456789", Email: "test@mail.ru", Address: &address_user0, Pets: &pets_user0})
	users = append(users, User{Id: 1, Name: "Александр", Surname: "Тестовый", Phone: "+79150554477", Pets: &pets_user1})
	users = append(users, User{Id: 2, Name: "Oleg", Surname: "Musin", Phone: "+79150554477", Address: &address_user2, Pets: &pets_user2})

	// Route handles & endpoints
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	r.HandleFunc("/pets/user/{userId}", getUserPets).Methods("GET")
	r.HandleFunc("/pets/{id}", getPet).Methods("GET")
	r.HandleFunc("/pets", getPets).Methods("GET")
	r.HandleFunc("/pets/{id}", updatePet).Methods("PUT")
	r.HandleFunc("/pets", createPet).Methods("POST")
	r.HandleFunc("/pets/{id}", deletePet).Methods("DELETE")

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
