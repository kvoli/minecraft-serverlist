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
var Votes map[int32]int32
var IDS int32

type Server struct {
	ServerID int32 `json:"serverId"`
	Address  string `json:"address"`
	Desc     string `json:"description"`
}

type Vote struct {
	ServerID int32 `json:"serverId"` 
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
	Votes[server.ServerID] = 0

	json.NewEncoder(w).Encode(server)
}

// vote for ur fav server , probs need a mutex here if u concurrent cuz race condition
func voteHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var vote Vote
	json.Unmarshal(reqBody, &vote)

	if v, ok := Votes[vote.ServerID] ; ok{
		Votes[vote.ServerID] = v + vote.Delta
		json.NewEncoder(w).Encode(Votes[vote.ServerID])
	}
}

func hiscoreHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Votes)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/server", serverHandler).Methods("GET")
	router.HandleFunc("/server", createServerHandler).Methods("POST")
	router.HandleFunc("/server/{id}", singleServerHandler).Methods("GET")
	router.HandleFunc("/vote", voteHandler).Methods("POST")
	router.HandleFunc("/hiscore", hiscoreHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}


func makeInit(){
	IDS = 3
	Servers = []Server{
		Server{ServerID: 1, Desc: "Christian, Family Friendly, Drug Free, No Swearing Minecraft Server", Address: "christiansdoitbest.faith"},
		Server{ServerID: 2, Desc: "Anarchy Server - fk sht up", Address: "anarchy.gg"},
	}
	Votes = make(map[int32]int32)
	Votes[1] = 420
	Votes[2] = 0
}


func main() {
	fmt.Println("Minecraft Server List API is live")
	makeInit()
	handleRequests()
}
