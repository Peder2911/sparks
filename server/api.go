/*
The SPRAWL game server API.
*/
package main

import (
	"encoding/json"
	"net/http"
	"github.com/peder2911/sparks/server/client"
	//"github.com/peder2911/sparks/server/ecs"
	"github.com/peder2911/sparks/server/game"
	"github.com/peder2911/sparks/server/session"
)

type MeResponse struct {
   Name string
}

type LoginHandshake struct {
   callback chan client.Client
}

func NewApiHandler(base_path string, game game.Game) (http.Handler) {
   mux := http.NewServeMux()
   mux.HandleFunc(base_path + "/me", func(w http.ResponseWriter, r *http.Request){
      username := r.Header.Get("X-Username")
      json.NewEncoder(w).Encode(MeResponse{username})
   })
   mux.HandleFunc(base_path + "/session", session.NewSessionHandler(game))
   return mux
}
