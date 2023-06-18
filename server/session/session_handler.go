package session

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/peder2911/sparks/server/client"
	"github.com/peder2911/sparks/server/gameserver"
)

var upgrader = websocket.Upgrader{
   ReadBufferSize: 1024,
   WriteBufferSize: 1024,
}

func NewSessionHandler(ctx context.Context, gameserver gameserver.GameServer) (func(w http.ResponseWriter, r *http.Request)) {
   return func(w http.ResponseWriter, r *http.Request){
         var err error
         //username := r.Header.Get("X-Username")
         con, err := upgrader.Upgrade(w, r, nil)
         if err != nil {
            log.Println(fmt.Sprintf("Failed to upgrade a request to websocket: %s", err))
            return
         }

         handshake := client.LoginHandshake{Callback: make(chan *client.Client)}
         gameserver.Logins <- handshake
         session_client := <- handshake.Callback
         close(handshake.Callback)
         session_client.ServeWs(ctx, con)
      }
}
