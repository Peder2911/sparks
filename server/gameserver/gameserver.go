/*
The Sparks game object

This object is responsible for holding the game state, as well as sending that state to players.

The game state is streamed in a delta-fashion, meaning that only changes to the
state are sent to clients, instead of the whole state. This saves a lot of
network traffic. Clients are "bootstrapped" when joining, receiving a one-time delta
of the whole game state.

*/

package gameserver

import (
   "time"
   "log"
   "fmt"
   "github.com/peder2911/sparks/server/ecs"
   "github.com/peder2911/sparks/server/client"
   "context"
)

type GameServer struct {
   ecs ecs.Ecs
   Clients map[int] client.Client
   Logins chan client.LoginHandshake
   Inputs chan client.IdentifiedClientMessage
   Unregister chan int
   Cancel context.CancelFunc
   Context context.Context
   //current_inputs []protocol.ClientMessage
   override_buffer []ecs.Delta
}

func NewGameServer(ctx context.Context) GameServer {
   cancelcontext, cancel  := context.WithCancel(ctx)
   game := GameServer{}
   game.ecs = ecs.NewEcs()
   game.Clients = make(map[int] client.Client)
   game.Logins = make(chan client.LoginHandshake)
   game.Inputs = make(chan client.IdentifiedClientMessage)
   game.Unregister = make(chan int) 
   game.Cancel = cancel 
   game.Context = cancelcontext
   game.override_buffer = nil
   return game
}

func (g *GameServer) Login() client.Client {
   new_client := client.NewClient(g.ecs.Create(ecs.DefaultEntity()), make(chan ecs.Delta), g.Inputs, g.Unregister)
   g.Clients[new_client.Id] = new_client 
   // Get new player up to speed
   go g.BootstrapClient(new_client.Id)
   return new_client
}

func (g *GameServer) Logout(index int) {
   g.ecs.Destroy(index)
   delete(g.Clients, index)
}

func (g *GameServer) Cleanup() {
}

func (g *GameServer) BootstrapClient(id int){
   client := g.Clients[id]
   for _,delta := range g.ecs.Snapshot() {
      client.Deltas <- delta
   }
}

func (g *GameServer) flush_overrides() []ecs.Delta {
   overrides := g.override_buffer
   g.override_buffer = nil
   return overrides 
}

func (g *GameServer) Iterate() []ecs.Delta {
   return g.ecs.Iterate(g.flush_overrides())
}

func (g *GameServer) Loop() {
   tick := time.Tick(10 * time.Millisecond)
   for {
      select {
         case handshake := <- g.Logins:
            new_client := g.Login()
            log.Println(fmt.Sprintf("Got a new login with id %v", new_client.Id))
            handshake.Callback <- &new_client
         case logmeout := <- g.Unregister:
            log.Println(fmt.Sprintf("Logging out %v", logmeout))
            g.ecs.Destroy(logmeout)
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
         case <- g.Context.Done():
            g.Cleanup()
            return
      }
   }
}
