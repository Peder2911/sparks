/*
The SPRAWL game server API.
Not much here yet!
*/
package main

import (
	"encoding/json"
	"net/http"
        "fmt"
        "log"
	"github.com/gorilla/websocket"
)

type MeResponse struct {
   Name string
}

var upgrader = websocket.Upgrader{
   ReadBufferSize: 1024,
   WriteBufferSize: 1024,
}

func NewApiHandler(base_path string, game Game) (http.Handler) {
   mux := http.NewServeMux()
   mux.HandleFunc(base_path + "/me", func(w http.ResponseWriter, r *http.Request){
      username := r.Header.Get("X-Username")
      json.NewEncoder(w).Encode(MeResponse{username})
   })
   mux.HandleFunc(base_path + "/session", func(w http.ResponseWriter, r *http.Request){
      var err error
      username := r.Header.Get("X-Username")
      conn, err := upgrader.Upgrade(w, r, nil)
      if err != nil {
         log.Println(fmt.Sprintf("Failed to upgrade a request to websocket: %s", err))
         return
      }
      for {
         player_index := game.Login()
         player := game.Entities[player_index]
         defer game.Logout(player_index)
         message_type, p, err := conn.ReadMessage()
         if err != nil {
            log.Println(fmt.Sprintf("Error while reading from websocket: %s", err))
            break
         }
         if message_type != websocket.TextMessage {
            log.Println(fmt.Sprintf("Received binary message (not supported)"))
            break
         }      
         err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello, %s! You said %v", username, string(p))))
         err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("You are located at %v, %v", player[1],player[2])))
         if err != nil {
            log.Println(fmt.Sprintf("Error while writing: %s", err))
            break
         }
      }
   })
   return mux
}
