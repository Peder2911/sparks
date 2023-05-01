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

      callback := make(chan chan Delta)
      game.Logins <- callback
      client_channel := <- callback 

      log.Println(fmt.Sprintf("Serving traffic to %s", username))
      
      for {
         select {
            case entity := <- client_channel:
               ws, err := conn.NextWriter(websocket.TextMessage)
               if err != nil {
                  break
               }
               json.NewEncoder(ws).Encode(entity)
               if err = ws.Close(); err != nil {
                  break
               }
         }
      }
   })
   return mux
}
