/*
The Sparks client.

This object is responsible for handling a single client session. It is
instantiated by the Game object and is passed back to a session where it is
hooked up to the websocket.
*/
package client

import (
   "github.com/gorilla/websocket"
   "github.com/peder2911/sparks/server/ecs"
   "github.com/peder2911/sparks/server/protocol"
   "log"
   "fmt"
)

type ClientInput struct {
   Id int
   Message protocol.ClientMessage
}

type Client struct {
   Id int
   Deltas chan ecs.Delta
   Inputs chan ClientInput
   Unregister chan int 
   delta_buffer []ecs.Delta
}

type LoginHandshake struct {
   Callback chan *Client
}

func NewClient(id int, deltas chan ecs.Delta, inputs chan ClientInput, unregister chan int) Client {
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

func (c *Client) ServeWs(con *websocket.Conn) {
   go c.receive(con)
   go c.send(con)
}

func (c *Client) send(con *websocket.Conn) {
   for {
      delta := <- c.Deltas
      message := protocol.ServerMessage{Delta: delta}
      con.WriteJSON(message)
   }
}

func (c *Client) cleanup(){
   c.Unregister <- c.Id
}

func (c *Client) receive(con *websocket.Conn) {
   defer c.cleanup()
   for {
      var message protocol.ClientMessage
      err := con.ReadJSON(&message)
      if err != nil {
         log.Println(fmt.Sprintf("Received weird input!: %s",err))
      }
      c.Inputs <- ClientInput{Id: c.Id, Message: message}
   }
}
