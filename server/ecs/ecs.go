/*
Sparks ECS

Sparks uses an ECS to hold and manipulate the game state.
*/
package ecs

import "math/rand"

const EntitySize = 11
const WorldSizeX = 400
const WorldSizeY = 400

// Entity indices
const (
   Status  int = iota
   X
   Y
   Xvel
   Yvel
   XMoving
   YMoving
   Speed
   Inertia
   Size
   Value
)

type Index struct {
   Current int
   Free []int
}

func (i *Index) Get() (int, bool){
   var next int
   var is_new bool
   if len(i.Free) > 0 {
      next, i.Free = i.Free[0], i.Free[1:]
      is_new = false
   } else {
      next = i.Current
      i.Current ++
      is_new = true
   }
   return next, is_new 
}

func (i *Index) Recycle(index int){
   i.Free = append(i.Free, index)
}

func NewIndex() Index{
   index := Index{}
   index.Current = 0
   index.Free = []int{}
   return index
}

type Entity [EntitySize]int

func DefaultEntity() Entity {
   entity := Entity{}
   entity[Status] = 1
   entity[X] = rand.Intn(WorldSizeX) 
   entity[Y] = rand.Intn(WorldSizeY) 
   entity[Xvel] = 0
   entity[Yvel] = 0
   entity[XMoving] = 0
   entity[YMoving] = 0
   entity[Speed] = 4
   entity[Inertia] = 1
   entity[Size] = 1
   entity[Value] = 0
   return entity
}

type Ecs struct {
   entities []Entity
   deltas []Delta
   Index Index
}

func NewEcs() Ecs{
   ecs := Ecs{}
   ecs.entities = []Entity{}
   ecs.deltas = nil
   ecs.Index = NewIndex()
   return ecs
}

func (e *Ecs) Grow(){
   e.entities = append(e.entities, Entity{})
}

func (e *Ecs) Create(entity Entity) int{
   i,should_grow := e.Index.Get()

   if should_grow {
      e.Grow()
   }

   for j,v := range entity {
      e.Set(i,j,v)
   }

   return i 
}

func (e *Ecs) Destroy(i int){
   e.Set(i, Status, 0)
   e.Index.Recycle(i)
}

func (e *Ecs) Get(i int, j int) int {
   return e.entities[i][j]
}

func (e *Ecs) Set(i int, j int, v int) {
   e.entities[i][j] = v
   e.deltas = append(e.deltas, Delta{i,j,v})
}

func (e *Ecs) Iterate(overrides []Delta) []Delta {

   for _,override := range overrides {
      e.Set(override[0],override[1],override[2])
   }

   e.physics()
   e.control()

   deltas := e.deltas
   e.deltas = nil 
   return deltas
}

func (e *Ecs) physics(){
   for i := 0 ; i < e.Index.Current ; i++ {
      if e.Get(i,Status) != 0 {
         xvel, yvel, inertia := e.Get(i, Xvel), e.Get(i, Yvel), e.Get(i, Inertia)
         var xinertia, yinertia int
 
         // Calculate and apply inertia for x and y
         if xvel < 0 {
            xinertia = inertia * -1
         } else if xvel > 0 {
            xinertia = inertia
         } else {
            xinertia = 0
         }
         if xinertia != 0 {
            e.Set(i,Xvel,xvel - xinertia)
         }

         if yvel < 0 {
            yinertia = inertia * -1
         } else if yvel > 0 {
            yinertia = inertia
         } else {
            yinertia = 0
         }
         if yinertia != 0 {
            e.Set(i,Yvel,yvel - yinertia)
         }

         // Apply speed to position
         // TODO colision detection?
         var new_position int
         if xvel := e.Get(i,Xvel); xvel != 0 {
            new_position = e.Get(i,X) + xvel
            if new_position < WorldSizeX && new_position > 0 {
               e.Set(i, X, new_position) 
            } else {
               e.Set(i, Xvel, -e.Get(i, Xvel))
            }
         }
         if yvel := e.Get(i,Yvel); yvel != 0 {
            new_position = e.Get(i,Y) + yvel
            if new_position < WorldSizeY && new_position > 0 {
               e.Set(i, Y, new_position) 
            } else {
               e.Set(i, Yvel, -e.Get(i, Yvel))
            }
         }
      }
   }
}

func (e *Ecs) control(){
   for i := 0 ; i < e.Index.Current ; i++ {
      if e.Get(i,Status) != 0 {
         if xmoving := e.Get(i,XMoving); xmoving != 0 {
            e.Set(i, Xvel, e.Get(i, Speed) * xmoving)
         }
         if ymoving := e.Get(i,YMoving); ymoving != 0 {
            e.Set(i, Yvel, e.Get(i, Speed) * ymoving)
         }
      }
   }
}

func (e *Ecs) Snapshot() []Delta {
   snapshot := make([]Delta, len(e.entities) * EntitySize)

   s := 0
   for i := range e.entities {
      for j := range e.entities[0] {
         snapshot[s] = Delta{i,j,e.Get(i,j)}
         s++
      }
   }
   return snapshot
}
