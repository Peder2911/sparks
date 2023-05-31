
class Application {
   constructor(root_element){
      this.el = root_element
      this.game = new Game(this.game_stopper())
      this.login_page = new LoginPage(this.game_starter())
      this.set_content(this.login_page.el())
   }

   set_content(content){
      this.el.innerHTML = ""
      this.el.appendChild(content)
   }

   game_starter(){
      var self = this
      return (token) => {
         self.game.initialize(token)
         self.set_content(self.game.el())
      }
   }

   game_stopper(){
      var self = this
      return () => {
         self.game.stop()
         self.set_content(self.login_page.el())
      }
   }
}
