

class Controller {
   constructor(socket){
      this.socket = socket
   }

   _translate_key(k){
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

   handle(event, direction){
      let command = this._translate_key(event.key)
      if (command !== undefined) {
         this.socket.send(JSON.stringify({action: direction, key: command}))
      }
   }

   keydown(){
      self = this
      return function(event){
         self.handle(event, "press")
      }
   }

   keyup(){
      self = this
      return function(event){
         self.handle(event, "release")
      }
   }
}
