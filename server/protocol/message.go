
package protocol

type Action string

const (
   Press Action = "press"
   Unpress      = "release"
)

type Key string

const (
   Up Key = "up"
   Down   = "down"
   Left   = "left"
   Right  = "right"
)

type ClientMessage struct {
   Action Action `json:"action"`
   Key    Key    `json:"key"`
}

type ServerMessage struct {
   Delta [3]int `json:"delta"`
}

type GoodbyeMessage struct {
   Message string `json:"message"`
}
