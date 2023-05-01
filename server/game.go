
package main

import (
   "time"
   "log"
)

/*
Entity schema
materialized (flag),
x,
y,
velocity x,
velocity y,
*/


type Game struct {
   ecs Ecs
   Clients map[int] chan Delta 
   Logins chan chan chan Delta 
   Cancel chan int
}

func NewGame() Game {
   game := Game{}
   game.Cancel = make(chan int)
   game.ecs = NewEcs()
   game.Clients = make(map[int]chan Delta)
   game.Logins = make(chan chan chan Delta)
   return game
}

func (g *Game) Login() chan Delta {
   client_id := g.ecs.Create(1,0,0,1,1)
   g.Clients[client_id] = make(chan Delta)
   // Get new player up to speed
   go g.BootstrapClient(client_id)
   return g.Clients[client_id]
}

func (g *Game) Logout(index int) {
   g.ecs.Destroy(index)
   delete(g.Clients, index)
}

func (g *Game) Cleanup() {
}

func (g *Game) BootstrapClient(client_id int){
   for _,delta := range g.ecs.Snapshot() {
      g.Clients[client_id] <- delta
   }
}

func (g *Game) UpdateClient(client_id int, deltas []Delta){
   for _, delta := range deltas {
      g.Clients[client_id] <- delta
   }
}

func (g *Game) Loop() {
   tick := time.Tick(1000 * time.Millisecond)
   for {
      select {
         case new_login := <- g.Logins:
            log.Println("Got a new login")
            new_login <- g.Login()
         case <- tick:
            deltas := g.ecs.Iterate()
            for c := range g.Clients {
               go g.UpdateClient(c, deltas)
            }
         case <- g.Cancel:
            g.Cleanup()
            return
      }
   }
}
