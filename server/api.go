/*
The SPRAWL game server API.
*/
package main

import (
	"encoding/json"
	"net/http"
	"github.com/peder2911/sparks/server/client"
	//"github.com/peder2911/sparks/server/ecs"
	"github.com/peder2911/sparks/server/gameserver"
	"github.com/peder2911/sparks/server/session"
        "context"
)

type MeResponse struct {
   Name string
}

type LoginHandshake struct {
   callback chan client.Client
}

func NewApiHandler(ctx context.Context, base_path string, gameserver gameserver.GameServer) (http.Handler) {
   mux := http.NewServeMux()
   mux.HandleFunc(base_path + "/me", func(w http.ResponseWriter, r *http.Request){
      username := r.Header.Get("X-Username")
      json.NewEncoder(w).Encode(MeResponse{username})
   })
   mux.HandleFunc(base_path + "/session", session.NewSessionHandler(ctx, gameserver))
   return mux
}
