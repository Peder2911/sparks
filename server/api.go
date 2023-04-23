/*
The SPRAWL game server API.
Not much here yet!
*/
package main
import (
   "net/http"
   "encoding/json"
)

type MeResponse struct {
   Name string
}

func NewApiHandler(base_path string) (http.Handler) {
   mux := http.NewServeMux()
   mux.HandleFunc(base_path + "/me", func(w http.ResponseWriter, r *http.Request){
      username := r.Header.Get("X-Username")
      json.NewEncoder(w).Encode(MeResponse{username})
   })
   return mux
}
