package main

import (
   "time"
   "fmt"
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

type Components [5]int


type Game struct {
   Entities []Components
   FreeIndices []int
   Broadcast []int
   Cancel chan int
   Clients map[int] chan Entity
   Logins chan chan chan Entity
}

type Entity struct {
   Index int `json:"index"`
   Components Components `json:"components"`
}

type GameConnection struct {
}

func NewGame() Game {
   game := Game{}
   game.Cancel = make(chan int)
   game.Entities = []Components{}
   game.FreeIndices = []int{}
   game.Clients = make(map[int]chan Entity)
   game.Logins = make(chan chan chan Entity)
   return game
}

func (g *Game) Create(x int, y int) int {
   var index int
   new_entity := Components{1,x,y,1,1}
   
   if len(g.FreeIndices) > 0 {
      log.Println("Using a free index for entity")
      index, g.FreeIndices = g.FreeIndices[0], g.FreeIndices[1:]
      g.Entities[index] = new_entity
   } else {
      g.Entities = append(g.Entities, new_entity)
      index = len(g.Entities) - 1
   }
   return index
}

func (g *Game) Destroy(index int) int {
   g.Entities[index] = Components{}
   g.FreeIndices = append(g.FreeIndices, index)
   return index
}
func (g *Game) Login() chan Entity {
   player_index := g.Create(0,0)
   log.Println(fmt.Sprintf("New login created player with ID %v", player_index))
   g.Clients[player_index] = make(chan Entity)
   go g.SendEntityToClient(player_index, player_index)
   for i := range g.Entities {
      go g.SendEntityToClient(player_index, i)
   }
   return g.Clients[player_index]
}

func (g *Game) SendEntityToClient(client_index int, entity_index int) {
   log.Println(fmt.Sprintf("Sending entity %v to client %v", entity_index, client_index))
   g.Clients[client_index] <- Entity{entity_index, g.Entities[entity_index]}
}

func (g *Game) Logout(index int) {
   g.Destroy(index)
   delete(g.Clients, index)
}

func (g *Game) Move() {
   for i,entity := range g.Entities {
      log.Println(fmt.Sprintf("Moving %v", i))
      if entity[0]!=0 && (entity[3]>0|| entity[4]>0) {
         entity[1] = entity[1] + entity[3]
         entity[2] = entity[2] + entity[4]
         g.Entities[i] = entity 
         g.Broadcast = append(g.Broadcast, i)
      }
   }
}

func (g *Game) Systems() {
   g.Move()
}

func (g *Game) Cleanup() {
}

func (g *Game) Loop() {
   tick := time.Tick(1000 * time.Millisecond)
   for {
      select {
         case new_login := <- g.Logins:
            log.Println("Sending client channel to client")
            new_login <- g.Login()
         case <- tick:
            log.Println(fmt.Sprintf("Systems %v", len(g.Entities)))
            g.Systems()
            for _,entity_index := range g.Broadcast {
               for client_index := range g.Clients {
                  g.SendEntityToClient(client_index, entity_index)
               }
            }
            g.Broadcast = nil
         case <- g.Cancel:
            g.Cleanup()
            return
      }
   }
}
