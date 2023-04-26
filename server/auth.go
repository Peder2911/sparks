/*
SPRAWLs authentication setup is relatively simple. It checks provided
credentials against the database, and returns a signed token if valid. This
token is then passed with subsequent requests and verified by a middleware that
checks its signature for validity.

TODOs:
- Salt the passwords
- Fix casing on JSON being passed back and forth
- Return "expires at" with token.
- Refresh tokens?
*/
package main

import (
	"fmt"
	"net/http"

	"encoding/json"
        "log"
        "time"
        "regexp"
        "database/sql"

        "github.com/golang-jwt/jwt/v4"
)

type login_request struct {
   UserName string `json:"username"`
   Password string `json:"password"`
}

type login_response struct {
   AccessToken string `json:"access_token"`
}

type JwtVerifier struct {
   handler http.Handler
   secret []byte 
}

func (verifier *JwtVerifier) Verify (token_string string) (*jwt.Token, error) {
   token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error){
      if token.Header["alg"].(string) != jwt.SigningMethodHS256.Alg() {
         return nil, fmt.Errorf("Token had unexpected alg: %s, expected %s", token.Header["alg"], jwt.SigningMethodRS256.Alg())
      }
      return verifier.secret, nil
   })

   if err != nil {
      return nil, err
   }

   if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
      if claims["aud"] != "sprawl.rackness.net/api" {
         return nil, fmt.Errorf("Token had wrong audience: %s", claims["aud"])
      }
      return token, nil
   } else {
      return nil, fmt.Errorf("Issue with token: %s", err)
   }
}

func (verifier *JwtVerifier) ServeHTTP (w http.ResponseWriter, r *http.Request) {
   var err error
   regex, err := regexp.Compile("(?:Bearer )(.*)")
   if err != nil {
      panic(fmt.Sprintf("Failed to compile regexp: %s", err))
   }

   match := regex.FindSubmatch([]byte(r.Header.Get("Authorization")))
   if len(match) != 2 {
      log.Println(fmt.Sprintf("Failed to parse token from Authorization header: %s", match))
      w.WriteHeader(401)
      return
   }

   token, err := verifier.Verify(string(match[1]))
   if err != nil {
      log.Println(fmt.Sprintf("Failed to verify token: %s", err))
      w.WriteHeader(401)
      return
   }

   claims := token.Claims.(jwt.MapClaims)
   r.Header.Set("X-Username", claims["name"].(string))
   r.Header.Set("X-Userid", fmt.Sprintf("%v",claims["sub"].(float64)))
   verifier.handler.ServeHTTP(w, r)
}

func NewJwtVerifier(handler_to_wrap http.Handler, secret []byte) *JwtVerifier {
   return &JwtVerifier{handler_to_wrap, secret}
}

func NewAuthHandler(secret []byte, mountpoint string, connect func()(*sql.DB)) *http.ServeMux {
   var auth = http.NewServeMux()

   auth.HandleFunc(mountpoint + "/register", func(w http.ResponseWriter, r *http.Request){
      var err error
      login_request := login_request{}
      err = json.NewDecoder(r.Body).Decode(&login_request)
      if err != nil {
         w.WriteHeader(400)
         return
      }
      db := connect()
      defer db.Close()

      _, err = db.Exec("insert into users (name, password) values (?, ?)", login_request.UserName, login_request.Password)
      if err != nil {
         log.Println(fmt.Sprintf("Failed to insert user: %s", err))
         w.WriteHeader(400)
         return
      }
      w.WriteHeader(201)
   })

   auth.HandleFunc(mountpoint + "/token", func(w http.ResponseWriter, r *http.Request){
      var err error
      // Fetch the user
      login_request := login_request{}
      err = json.NewDecoder(r.Body).Decode(&login_request)
      if err != nil {
         w.WriteHeader(400)
         return
      }

      db := connect()
      defer db.Close()

      user_name := login_request.UserName
      var user_id int
      var user_password string
      user_row := db.QueryRow("select id, password from users where users.name=?", user_name) 
      err = user_row.Scan(&user_id, &user_password)
      if err != nil {
         w.WriteHeader(401)
         log.Println(fmt.Sprintf("Failed to fetch user: %s", err))
         return
      }
      if user_password != login_request.Password {
         log.Println(fmt.Sprintf("Password verification failed: %s != %s", user_password, login_request.Password))
         w.WriteHeader(401)
         return
      }

      timestamp := time.Now().Unix() 

      token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
         "alg": jwt.SigningMethodHS256.Alg(),
         "typ": "JWT",
         "sub": user_id,
         "name": user_name,
         "iat": timestamp,
         "aud": "sprawl.rackness.net/api",

      })
      
      token_string, err := token.SignedString(secret)

      if err != nil {
         w.WriteHeader(500)
         log.Println(fmt.Sprintf("Failed to sign JWT: %s", err))
         return 
      }

      response := login_response{
         AccessToken: token_string,
      }

      err = json.NewEncoder(w).Encode(response)
      if err != nil {
         w.WriteHeader(500)
         log.Println(fmt.Sprintf("Failed to serialize JWT: %s", err))
         return
      }
      log.Println(fmt.Sprintf("Served a token to %s", user_name))
   })
   return auth
}
