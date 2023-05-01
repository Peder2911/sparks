package main

import (
   "log"
   "fmt"
)

// Entity indices
const (
   Status  int = 0
   X           = 1
   Y           = 2
   Xvel        = 3
   Yvel        = 4
)

type Delta [3]int

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

type Entity [5] int

type Ecs struct {
   entities []Entity
   //Delta chan Delta
   deltas []Delta

   Index Index
}

func NewEcs() Ecs{
   ecs := Ecs{}
   ecs.entities = []Entity{}
   ecs.deltas = nil
   //ecs.Delta = make(chan Delta)
   ecs.Index = NewIndex()
   return ecs
}

func (e *Ecs) Grow(){
   e.entities = append(e.entities, Entity{})
}

func (e *Ecs) Create(active int, x int, y int, xvel int, yvel int) int{
   i,should_grow := e.Index.Get()

   if should_grow {
      e.Grow()
   }

   for j,v := range [5]int{active, x, y, xvel, yvel} {
      e.entities[i][j] = v
   }

   return i 
}

func (e *Ecs) Destroy(i int){
   e.entities[i][0] = 0
   //c.Active[index] = false
   e.Index.Recycle(i)
}

func (e *Ecs) Get(i int, j int) int {
   return e.entities[i][j]
}

func (e *Ecs) Set(i int, j int, v int) {
   e.entities[i][j] = v
   e.deltas = append(e.deltas, Delta{i,j,v})
}

func (e *Ecs) Iterate() []Delta {
   log.Println(fmt.Sprintf("Iterating with %v entities", len(e.entities)))
   e.Move()
   deltas := e.deltas
   e.deltas = nil 
   return deltas
}

func (e *Ecs) Move(){
   for i := 0 ; i < e.Index.Current ; i++ {
      if e.Get(i,Status) != 0 {
         if xvel := e.Get(i,Xvel); xvel > 0 {
            e.Set(i, X, e.Get(i, X) + e.Get(i, Xvel))
         }
         if yvel := e.Get(i,Yvel); yvel > 0 {
            e.Set(i, Y, e.Get(i, Y) + e.Get(i, Yvel))
         }
      }
   }
}

func (e *Ecs) Snapshot() []Delta {
   log.Println("Preparing snapshot")
   snapshot := make([]Delta, len(e.entities) * 5)

   s := 0
   for i := range e.entities {
      for j := range e.entities[0] {
         log.Println(fmt.Sprintf("Sending %v,%v to snapshot",i,j))
         snapshot[s] = Delta{i,j,e.Get(i,j)}
         s++
      }
   }
   return snapshot
}
