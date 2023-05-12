package main

const (
   KeyDown int = iota
   KeyUp
)

const (
   Up   int = iota
   Down
   Left
   Right
)

type Button int
type Pressed int

type Input struct {
   player int
   button Button
   pressed Pressed
}
