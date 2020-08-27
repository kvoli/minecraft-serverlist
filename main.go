package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
)

// global var ... lets say we have an actual DB here......
var Servers []Server
var Votes map[string]int32
var IDS int32

type Server struct {
	ServerID int32 `json:"serverId"`
	Address  string `json:"address"`
	Desc     string `json:"description"`
}

type Vote struct {
	ServerID string `json:"serverId"` 
	Delta		int32   `json:"delta"`
}


// return all servers
func serverHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Servers)
}

// return single servers
func singleServerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	for _, s := range Servers {
		if fmt.Sprint(s.ServerID) == key {
			json.NewEncoder(w).Encode(s)
		}
	}
}

// create a new server (POST)
func createServerHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var server Server
	json.Unmarshal(reqBody, &server)
	
	server.ServerID = IDS
	IDS += 1
	Servers = append(Servers, server)

	json.NewEncoder(w).Encode(server)
}

// vote for ur fav server , probs need a mutex here if u concurrent cuz race condition
func voteHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var vote Vote
	json.Unmarshal(reqBody, &vote)

	Votes[vote.ServerID]+=1

	json.NewEncoder(w).Encode(Votes[vote.ServerID])
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/server", serverHandler).Methods("GET")
	router.HandleFunc("/server", createServerHandler).Methods("POST")
	router.HandleFunc("/server/{id}", singleServerHandler).Methods("GET")
	router.HandleFunc("/vote", voteHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("Minecraft Server List API")
	IDS = 3
	Servers = []Server{
		Server{ServerID: 1, Desc: "Christian, Family Friendly, Drug Free, No Swearing Minecraft Server", Address: "christiansdoitbest.faith"},
		Server{ServerID: 2, Desc: "Anarchy Server - fk sht up", Address: "anarchy.gg"},
	}
	handleRequests()
}
