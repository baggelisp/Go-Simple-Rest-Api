package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

type user struct {
	id       string `json:"id"`
	name     string `json:"name"`
	username string `json:"username"`
	email    string `json:"email"`
}

type allUser []user

var users = allUser{
	{
		id:       "1",
		name:     "Donna Nash",
		email:    "donnanash@rockyard.com",
		username: "Carver",
	},
	{
		id:       "2",
		name:     "Joann Albert",
		email:    "joannalbert@rockyard.com",
		username: "Felecia",
	},
	{
		id:       "3",
		name:     "Hamilton Guerrero",
		email:    "hamiltonguerrero@rockyard.com",
		username: "Rosario",
	},
	{
		id:       "4",
		name:     "Beverley Oliver",
		email:    "beverleyoliver@rockyard.com",
		username: "Dianne",
	},
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the user title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newUser)
	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}

func getOneUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	for _, singleUser := range users {
		if singleUser.id == userID {
			json.NewEncoder(w).Encode(singleUser)
		}
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", users)
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	var updatedUser user

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the user title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedUser)

	for i, singleUser := range users {
		if singleUser.id == userID {
			singleUser.name = updatedUser.name
			singleUser.username = updatedUser.username
			singleUser.email = updatedUser.email
			users = append(users[:i], singleUser)
			json.NewEncoder(w).Encode(singleUser)
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	for i, singleUser := range users {
		if singleUser.id == userID {
			users = append(users[:i], users[i+1:]...)
			fmt.Fprintf(w, "The user with ID %v has been deleted successfully", userID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getOneUser).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
