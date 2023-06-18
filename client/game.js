
class Game {
   constructor(stop_game){
      console.log("setting up a brand new game")
      this.canvas = document.createElement("canvas")
      this.canvas.width = 720
      this.canvas.height = 480
      this.ctx = this.canvas.getContext("2d")
      this.entities = new(Array) 

      this.stop_button = document.createElement("button")
      this.stop_button.innerHTML = "Stop"
      this.stop_button.onclick = stop_game
   }

   stop(){
      this.socket.send(JSON.stringify({message:"logout"}))
      this.socket.close()
      console.log("Stopping the game")
   }

   receiver(){
      let self = this
      return function(m){
         let data = JSON.parse(m.data)
         let delta = data.delta
         let entity = self.entities[delta[0]]
         if (self.entities[delta[0]] === undefined) {
            self.entities[delta[0]] = new(Array)
         }
         console.log(delta)
         self.entities[delta[0]][delta[1]] = delta[2]
      }
   }
   
   clear(){
      this.ctx.fillStyle = "white"
      this.ctx.fillRect(0,0,this.canvas.width,this.canvas.height)
   }

   blitter(){
      let self = this
      return function(){
         self.clear()
         self.ctx.fillStyle = "black"
         self.entities.forEach(e => {
            if(e[0]){
               self.ctx.fillRect(e[1],e[2],10,10)
            }
         })
      }
   }


   looper(){
      let self = this
      let blit = self.blitter()
      let loop = function(){
         blit()
         window.requestAnimationFrame(loop)
      }
      return loop
   }

   translate_key(k){
      switch (k) {
         case "ArrowUp":
            return "up"
            break
         case "ArrowDown":
            return "down"
            break
         case "ArrowLeft":
            return "left"
            break
         case "ArrowRight":
            return "right"
            break
         default:
            return undefined
      }
   }

   keypress_handler(direction, socket){
      let self = this
      return function(e){
         let command = self.translate_key(e.key)
         if (command !== undefined) {
            socket.send(JSON.stringify({action: direction, key: self.translate_key(e.key)}))
         }
      }
   }

   initialize(token){
      console.log("Starting the game!!")
      this.socket = new WebSocket(`ws://${document.location.host}/api/api/session?token=${token}`)
      this.socket.onmessage = this.receiver()
      document.addEventListener("keydown", this.keypress_handler("press", this.socket))
      document.addEventListener("keyup", this.keypress_handler("release", this.socket))
      this.looper()()
   }

   el(){
      let div = document.createElement("div")
      div.appendChild(this.canvas)
      div.appendChild(this.stop_button)
      return div 
   }
}
