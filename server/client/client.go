// The client package contains code for handling client sessions.
package client

import (
	"context"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/peder2911/sparks/server/ecs"
	"github.com/peder2911/sparks/server/protocol"
)

type Client struct {
   Id int
   Deltas chan ecs.Delta
   Inputs chan IdentifiedClientMessage
   Unregister chan int
   delta_buffer []ecs.Delta
}

type IdentifiedClientMessage struct {
   Id int
   Message protocol.ClientMessage
}

type LoginHandshake struct {
   Callback chan *Client
}

func NewClient(id int, deltas chan ecs.Delta, inputs chan IdentifiedClientMessage, unregister chan int) Client {
   return Client {
      Id: id,
      Deltas: deltas,
      Inputs: inputs,
      Unregister: unregister,
      delta_buffer: make([]ecs.Delta, 128),
   }
}

func (c *Client) flush_deltas() []ecs.Delta {
   deltas := c.delta_buffer
   c.delta_buffer = nil
   return deltas
}

// ServeWs
//
// Send and receive data with the clients websocket.
func (c *Client) ServeWs(ctx context.Context, con *websocket.Conn) {
   cctx, cancel := context.WithCancel(ctx)
   go c.receive(cctx, cancel, con)
   go c.send(cctx, con)
}

// send
//
// Send data to the clients websocket. This goroutine also handles cleanup if
// the context is cancelled.
func (c *Client) send(ctx context.Context, con *websocket.Conn) {
   for {
      select {
         case delta := <- c.Deltas:
            message := protocol.ServerMessage{Delta: delta}
            con.WriteJSON(message)
         case <- ctx.Done():
            log.Println(fmt.Sprintf("Cleaning up client %v", c.Id))
            goodbye := protocol.GoodbyeMessage{Message:"goodbye!"}
            con.WriteJSON(goodbye)
            c.Unregister <- c.Id
            return
      }
   }
}

// receive
//
// Receive data from the clients websocket. The client will end its session if
// the input cannot be deserialized as a ClientMessage.
func (c *Client) receive(ctx context.Context, cancel context.CancelFunc, con *websocket.Conn) {
   for {
      var message protocol.ClientMessage
      err := con.ReadJSON(&message)
      if err != nil {
         log.Println(fmt.Sprintf("Received weird input! Ending session: %s",err))
         cancel()
         return
      }
      c.Inputs <- IdentifiedClientMessage{Id: c.Id, Message: message}
   }
}
