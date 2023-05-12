/*
The Sparks game object

This object is responsible for holding the game state, as well as sending that state to players.

The game state is streamed in a delta-fashion, meaning that only changes to the
state are sent to clients, instead of the whole state. This saves a lot of
network traffic. Clients are "bootstrapped" when joining, receiving a one-time delta
of the whole game state.

*/

package game 

import (
   "time"
   "log"
   "fmt"
   "github.com/peder2911/sparks/server/ecs"
   "github.com/peder2911/sparks/server/client"
)


type Game struct {
   ecs ecs.Ecs
   Clients map[int] client.Client
   Logins chan client.LoginHandshake
   Inputs chan client.ClientInput
   Unregister chan int
   Cancel chan int
   //current_inputs []protocol.ClientMessage
   override_buffer []ecs.Delta
}

func NewGame() Game {
   game := Game{}
   game.ecs = ecs.NewEcs()
   game.Clients = make(map[int] client.Client)
   game.Logins = make(chan client.LoginHandshake)
   game.Inputs = make(chan client.ClientInput)
   game.Unregister = make(chan int) 
   game.Cancel = make(chan int) 
   game.override_buffer = nil
   return game
}

func (g *Game) Login() client.Client {
   new_client := client.NewClient(g.ecs.Create(ecs.DefaultEntity()), make(chan ecs.Delta), g.Inputs, g.Unregister)
   g.Clients[new_client.Id] = new_client 
   // Get new player up to speed
   go g.BootstrapClient(new_client.Id)
   return new_client
}

func (g *Game) Logout(index int) {
   g.ecs.Destroy(index)
   delete(g.Clients, index)
}

func (g *Game) Cleanup() {
}

func (g *Game) BootstrapClient(id int){
   client := g.Clients[id]
   for _,delta := range g.ecs.Snapshot() {
      client.Deltas <- delta
   }
}

func (g *Game) flush_overrides() []ecs.Delta {
   overrides := g.override_buffer
   g.override_buffer = nil
   return overrides 
}

func (g *Game) Iterate() []ecs.Delta {
   return g.ecs.Iterate(g.flush_overrides())
}

func (g *Game) Loop() {
   tick := time.Tick(1000 * time.Millisecond)
   for {
      select {
         case handshake := <- g.Logins:
            new_client := g.Login()
            log.Println(fmt.Sprintf("Got a new login with id %v", new_client.Id))
            handshake.Callback <- &new_client
         case logmeout := <- g.Unregister:
            delete(g.Clients,logmeout)
         case input := <- g.Inputs:
            input_delta, err := client.InputToDelta(input)
            if err == nil {
               g.override_buffer = append(g.override_buffer, input_delta)
            }
         case <- tick:
            deltas := g.Iterate()
            for _,client := range g.Clients {
               for _,delta := range deltas {
                  client.Deltas <- delta
               }
            }
         case <- g.Cancel:
            g.Cleanup()
            return
      }
   }
}
