
class LoginPage {
   constructor(on_login){
      this.on_login = on_login
   }

   el(){
      let form = document.createElement("form")

      let username = document.createElement("input")
      username.type = "text"

      let password = document.createElement("input")
      password.type = "text"

      let submit = document.createElement("input")
      submit.onclick = (e) => {
         e.preventDefault()
         this.handle_login(username.value, password.value)
         username.value = ""
         password.value = ""
      }
      submit.type = "submit"

      new Array(username,password,submit).forEach(el => form.appendChild(el)) 
      return form
   }

   handle_login(username, password){
      console.log(`Trying to log in with ${username} ${password}`)
      this.on_login("mocktoken")
   }
}
