
class Game {
   constructor(stop_game){
      console.log("setting up a brand new game")
      this.canvas = document.createElement("canvas")
      this.canvas.width = 720
      this.canvas.height = 480
      this.ctx = this.canvas.getContext("2d")
      this.entities = []

      this.stop_button = document.createElement("button")
      this.stop_button.innerHTML = "Stop"
      this.stop_button.onclick = stop_game
   }

   stop(){
      console.log("Stopping the game")
   }

   initialize(token){
      console.log("Starting the game!!")
      this.ctx.fillRect(10,10,10,10)
      this.socket = new WebSocket(`ws://${document.location.host}/api/api/session?token=${token}`)
      this.socket.onmessage = (e) => {
         console.log(e)
      }
   }

   el(){
      let div = document.createElement("div")
      div.appendChild(this.canvas)
      div.appendChild(this.stop_button)
      return div 
   }
}
