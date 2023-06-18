
package client 

import (
   "github.com/peder2911/sparks/server/ecs"
   "github.com/peder2911/sparks/server/protocol"
   "fmt"
)

func InputToDelta(input IdentifiedClientMessage) (ecs.Delta, error) {
   var err error
   var message protocol.ClientMessage = input.Message
   var delta ecs.Delta = ecs.Delta{input.Id,0,0}
   if message.Action == protocol.Press {
      if message.Key == protocol.Up {
         delta[1] = ecs.YMoving
         delta[2] = -1
      } else if message.Key == protocol.Down {
         delta[1] = ecs.YMoving
         delta[2] = 1
      } else if message.Key == protocol.Left {
         delta[1] = ecs.XMoving
         delta[2] = -1 
      } else if message.Key == protocol.Right {
         delta[1] = ecs.XMoving
         delta[2] = 1 
      } else {
         err = fmt.Errorf("Failed to parse input into delta.") 
         return delta, err
      }
   } else if message.Action == protocol.Unpress {
      if message.Key == protocol.Up {
         delta[1] = ecs.YMoving
         delta[2] = 0
      } else if message.Key == protocol.Down {
         delta[1] = ecs.YMoving
         delta[2] = 0
      } else if message.Key == protocol.Left {
         delta[1] = ecs.XMoving
         delta[2] = 0 
      } else if message.Key == protocol.Right {
         delta[1] = ecs.XMoving
         delta[2] = 0 
      } else {
         err = fmt.Errorf("Failed to parse input into delta.") 
         return delta, err
      }
   } else {
      err = fmt.Errorf("Failed to parse input into delta.") 
      return delta, err
   }


   return delta, err
}
