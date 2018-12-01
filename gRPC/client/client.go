package main

import (
	"fmt"
	"gRPC/client/auth"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/login",auth.Login)
	http.HandleFunc("/wrongPassword", auth.WrongPassword)
	http.HandleFunc("/registerSuccess", auth.RegisterSuccess)
	http.HandleFunc("/registerFail", auth.RegisterFail)
	http.HandleFunc("/personalPage", auth.PersonalPage)

	fmt.Println("Load success")
	er := http.ListenAndServe(":9090",nil)

	if er != nil {
		log.Fatal("ListenAndServer: ", er)
	}
	fmt.Println("working...")
}