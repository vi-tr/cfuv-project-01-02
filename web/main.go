package main

import (
    "fmt"
    "net/http"

    //"github.com/golang-jwt/jwt/v5"
)
import _ "embed"

//go:embed index.html
var page string

func main() {
    //jwt.New(jwt.SigningMethodES256, nil)

    http.HandleFunc("/admin", func(w http.ResponseWriter, _ *http.Request) {
        fmt.Fprint(w, page)
    })
    http.ListenAndServe(":8082", nil)
}

/*
var clients = make(map[*websocket.Conn]bool)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Error upgrading to WebSocket:", err)
        return
    }
    defer conn.Close()

    clients[conn] = true

    // Handle messages from this client
    for {
        // Read message from client
        _, msg, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading WebSocket message:", err)
            delete(clients, conn)
            break
        }

        // Process the message (e.g., add user to the list)
        // ...

        // Broadcast updated user list to all clients
        userList := getUsers() // Your function to get the user list
        for client := range clients {
            err := client.WriteJSON(userList)
            if err != nil {
                log.Println("Error sending WebSocket message:", err)
            }
        }
    }
}
*/
