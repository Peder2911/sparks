
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

      // TODO cannot send headers when opening a new websocket!
      //
      // This means I need to make some backend changes hmmmmm
      this.socket = new WebSocket("ws://localhost:8080/api/session")
      this.socket.onmessage = (e) => {
         console.log(e)
      }
   }

   stop(){
      console.log("Stopping the game")
   }

   initialize(token){
      this.ctx.fillRect(10,10,10,10)
      console.log("Starting the game!!")
   }

   el(){
      let div = document.createElement("div")
      div.appendChild(this.canvas)
      div.appendChild(this.stop_button)
      return div 
   }
}
