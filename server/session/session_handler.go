
package session

import (
   "net/http"
   "github.com/peder2911/sparks/server/game"
   "github.com/peder2911/sparks/server/client"
   "fmt"
   "log"
   "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
   ReadBufferSize: 1024,
   WriteBufferSize: 1024,
}

func NewSessionHandler(game game.Game) (func(w http.ResponseWriter, r *http.Request)) {
   return func(w http.ResponseWriter, r *http.Request){
         var err error
         //username := r.Header.Get("X-Username")
         con, err := upgrader.Upgrade(w, r, nil)
         if err != nil {
            log.Println(fmt.Sprintf("Failed to upgrade a request to websocket: %s", err))
            return
         }

         handshake := client.LoginHandshake{Callback: make(chan *client.Client)}
         game.Logins <- handshake
         session_client := <- handshake.Callback
         close(handshake.Callback)
         session_client.ServeWs(con)
      }
}
