
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
         self.entities.forEach((e,i) => {
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

   initialize(token){
      console.log("Starting the game!!")
      this.socket = new WebSocket(`ws://${document.location.host}/api/api/session?token=${token}`)
      this.controller = new Controller(this.socket)
      this.socket.onmessage = this.receiver()
      document.addEventListener("keydown", this.controller.keydown())
      document.addEventListener("keyup", this.controller.keyup())
      this.looper()()
   }

   el(){
      let div = document.createElement("div")
      div.appendChild(this.canvas)
      div.appendChild(this.stop_button)
      return div 
   }
}
