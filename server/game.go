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

type Game struct {
   Entities [][5]int
   FreeIndices []int
   Cancel chan int
}

func NewGame() Game {
   game := Game{}
   game.Cancel = make(chan int)
   game.Entities = [][5]int{} 
   game.FreeIndices = []int{} 
   return game
}

func (g *Game) Create(x int, y int) int {
   var index int
   if len(g.FreeIndices) > 0 {
      index, g.FreeIndices = g.FreeIndices[0], g.FreeIndices[1:]
   }
   new_entity := [5]int{1,x,y,1,1}
   if index+1 > len(g.Entities){
      g.Entities = append(g.Entities, new_entity)
   } else {
      g.Entities[index] = new_entity
   }
   log.Println(fmt.Sprintf("Created new entity: %v", g.Entities))
   return index
}

func (g *Game) Destroy(index int) int {
   g.Entities[index] = [5]int{}
   g.FreeIndices = append(g.FreeIndices, index)
   return index
}
func (g *Game) Login() int {
   return g.Create(0,0)
}

func (g *Game) Logout(index int) {
   g.Destroy(index)
}

func (g *Game) Move() {
   for i,entity := range g.Entities {
      log.Println(fmt.Sprintf("Moving %v", i))
      if entity[0]!=0{
         entity[1] = entity[1] + entity[3]
         entity[2] = entity[2] + entity[4]
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
         case <- tick:
            log.Println(fmt.Sprintf("Systems %v", g.Entities))
            g.Systems()
         case <- g.Cancel:
            g.Cleanup()
            return
      }
   }
}
