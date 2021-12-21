package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

var router = mux.NewRouter()



func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/", indexPage)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("logout", logout).Methods("POST")
	router.HandleFunc("example", examplePage)
	router.HandleFunc("/signup", signup).Methods("POST", "GET")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)

	fmt.Println("Welcome to Goodman Chat Pole")





	//Start the server on localhost port 8000 and log any possible error
	log.Println("http server started on :8080")
	err:=http.ListenAndServe(":8000",nil)
	if err!=nil {
		log.Fatal("ListenAndServe",err)
	}


}

func signup(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		tmpl, _ := template.ParseFiles("signup.html", "index.html", "base.html")
		u := &User{}
		tmpl.ExecuteTemplate(writer, "base", u)

	case "POST":
		f := request.FormValue("fName")
		l := request.FormValue("lName")
		em := request.FormValue("email")
		un := request.FormValue("userName")
		pass := request.FormValue("password")

		U := &User{Fname: f, Lname: l, Email: em, Username: un, Password: pass}
		print(U.Username)
		//Append second line
		file, err := os.OpenFile("data.txt",os.O_APPEND|os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		scanner:=bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(),U.Username){
				http.Redirect(writer, request, "http://localhost:8000/", 302)
			}else {
				hash, err := bcrypt.GenerateFromPassword([]byte(U.Password), bcrypt.MinCost)
				if err != nil {
					log.Println(err)
				}
				str:=U.Username+" "+string(hash)+" "+"\n"
				print(str)
				//w:=bufio.NewWriter(file)
				if _, err := file.WriteString(str); err != nil {
					log.Fatal(err)
				}
				break
			}

		}
		// First User is going to be signed up
		if scanner.Scan()==false {
			hash, err := bcrypt.GenerateFromPassword([]byte(U.Password), bcrypt.MinCost)
			if err != nil {
				log.Println(err)
			}
			str:=U.Username+" "+string(hash)+"\n"
			print(str)
			//w:=bufio.NewWriter(file)
			if _, err := file.WriteString(str); err != nil {
				log.Fatal(err)
			}
		}

		setSession(U, writer)
		//cmnd := exec.Command("room/chatRoom")
		//cmnd.Run() // and wait
		//cmnd.Start()
		//username := getUserName(request)
		http.Redirect(writer, request, "http://localhost:8080/", 302)

		//tmpl, _ := template.ParseFiles("room/index.html")
		//if username != "" {
		//	err := tmpl.ExecuteTemplate(writer, "base", &User{Username: username, Email: U.Email})
		//	if err != nil {
		//		http.Error(writer, err.Error(), http.StatusInternalServerError)
		//
		//	}
		//}

	}
}

func examplePage(writer http.ResponseWriter, request *http.Request) {
	tmpl, _ := template.ParseFiles("base.html", "index.html", "main.html")
	username := getUserName(request)
	if username != "" {
		err := tmpl.ExecuteTemplate(writer, "base", &User{Username: username})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)

		}
	}
}

func logout(writer http.ResponseWriter, request *http.Request) {
	clearSession(writer)
	http.Redirect(writer, request, "/", 302)

}

func login(writer http.ResponseWriter, request *http.Request) {
	name := request.FormValue("uname")
	pass := request.FormValue("password")
	//print("This is pass:",name)
	//print("\n")
	//print(pass)
	file, err := os.OpenFile("data.txt",os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	scanner:=bufio.NewScanner(file)

	//redirect := "http://localhost:8000"
	if name != "" && pass != "" {

		for scanner.Scan() {
			//print(scanner.Text())
			//print("\n")
			split:=strings.Split(scanner.Text()," ")
			print(split[0],split[1])
			print("\n")
			flagPass1:=false
			flagPass:=bcrypt.CompareHashAndPassword([]byte(split[1]), []byte(pass))
			if flagPass==nil{
				flagPass1 =true
			}
			flagUname:= strings.Contains(split[0],name)
			print(flagUname,flagPass1)
			if  flagUname && flagPass1{
				setSession(&User{Username: name, Password: pass}, writer)
				//cmnd := exec.Command("room/chatRoom.exe")
				//cmnd.Run() // and wait
				////cmnd.Start()
				//log.Println("log")
				redirect := "http://localhost:8080/"
				print("This is username",name,pass)
				http.Redirect(writer, request, redirect, 302)
				//file.Close()
				//return
			}
			//http.Redirect(writer, request,"http://localhost:8000" , 302)

		}
	}

}

func indexPage(writer http.ResponseWriter, request *http.Request) {
	u := &User{}
	tmpl, _ := template.ParseFiles("base.html", "index.html", "main.html")
	err := tmpl.ExecuteTemplate(writer, "base", u)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}


