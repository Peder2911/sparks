
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

   async handle_login(username, password){
      let response = await fetch("/api/auth/token",{
         method : "POST",
         body: JSON.stringify({username: username, password: password})
      })
      let data = await response.json()
      this.on_login(data.access_token)
   }
}
