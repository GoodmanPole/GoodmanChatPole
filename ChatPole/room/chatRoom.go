package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)
//connected clients
var  clients =make(map[*websocket.Conn]bool)

//broadcast channel
var broadcast = make(chan Message)

//configure an upgrader
var upgrader = websocket.Upgrader{}

//Message Object
type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Message string`json:"message"`
}


func main() {
	fmt.Println("Welcome to Goodman Chat Pole")

	//Create a simple file server
	//create a static file server and tie that to the "/" route so that when a user accesses the site they will be able to view index.html and any assets.
	fs:=http.FileServer(http.Dir("../room"))
	http.Handle("/",fs)

	//Configure Websocket route
	http.HandleFunc("/ws",handleConnections)


	//Start listening for incoming chat message
	go handleMessage()


	//Start the server on localhost port 8000 and log any possible error
	log.Println("http server started on :8080")
	err:=http.ListenAndServe(":8080",nil)
	if err!=nil {
		log.Fatal("ListenAndServe",err)
	}






}

func handleMessage() {
	for  {
		// Grab the next message from the broadcast channel
		msg:= <-broadcast
		// Send The message out to every online client

		for client :=range clients{
			err:=client.WriteJSON(msg)
			if err!=nil {
				log.Printf("error: %v",err)
				client.Close()
				delete(clients,client)

			}

		}

	}
}

func handleConnections(writer http.ResponseWriter, request *http.Request) {
	//Upgrade initial GET request to a websocket
	ws,err:=upgrader.Upgrade(writer,request,nil)
	if err!=nil {
		log.Fatal(err)

	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	//Register our new client
	clients[ws]=true

	for  {
		var msg Message
		// Read in a new json object and map it to a Message object
		err:=ws.ReadJSON(&msg)
		if err!=nil {
			log.Printf("error: %v",err)
			delete(clients,ws)
			break

		}

		// Send the new message to broadcast channel
		broadcast <- msg


	}
}
