package auth

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/twitter/proto"
	"html/template"
	"log"
	"net/http"
	"time"
)

const (
	address = "localhost:50051"
)

var user = &pb.Username{}

//创建服务端：grpc.NewServer()；注册服务：pb.RegisterHelloServiceServer()；启动服务端：s.Serve(lis)。
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Working")
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := pb.NewTwitterActionClient(conn)
	defer conn.Close()

	if r.Method == "GET" {
		t,_ := template.ParseFiles("show/login1.html")
		t.Execute(w, nil)
	} else {

		redirectAddress := ""
		r.ParseForm()
		username := r.Form["username"][0]
		password := r.Form["password"][0]
		method := r.Form["lr"][0]
		user.Name = username
		info := &pb.User{UserName:username, PassWord:password}
		ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
		if method == "login" {
			isRight,err := c.LoginCheck(ctx, info)
			if err != nil {fmt.Println(err)}
			//fmt.Println(isRight.IsTrue)

			if isRight.IsTrue == true {
				redirectAddress = "personalPage"
			} else {
				redirectAddress = "wrongPassword"
			}
			http.Redirect(w, r, redirectAddress, http.StatusFound)

		}else {
			isRight,_ := c.RegisterCheck(ctx, info)
			if isRight.IsTrue == true {
				redirectAddress = "registerSuccess"
			} else {
				redirectAddress = "registerFail"
			}
			defer c.LoginCheck(ctx, info)
			http.Redirect(w, r, redirectAddress, http.StatusFound)
		}
		defer cancel()
	}
}

func WrongPassword (w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t,_ := template.ParseFiles("show/wrongPassword.html")
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "login", http.StatusFound)
	}
}

func RegisterSuccess (w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t,_ := template.ParseFiles("show/registerSuccess.html")
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "login", http.StatusFound)
	}
}

func RegisterFail (w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t,_ := template.ParseFiles("show/registerFail.html")
		t.Execute(w, nil)
	} else {
		http.Redirect(w, r, "login", http.StatusFound)
	}
}


func PersonalPage(w http.ResponseWriter, r *http.Request){
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTwitterActionClient(conn)

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	//username,_ := c.GetName(ctx, &pb.Ack{})
	username := user
	//fmt.Println(username.Name)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("show/personalPage.html")
		//ack := &pb.Username{}
		//pagecontent,_ := c.GetTwitterPage(ctx, ack)
		pagecontent,_ := c.GetTwitterPage(ctx, username)
		err := t.Execute(w, pagecontent)
		if err != nil {fmt.Println(err)}

	} else {
		r.ParseForm()
		logout := r.Form.Get("logout")
		if logout == "logout" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		method := r.Form.Get("pg")
		switch method {
		case "Send Twitte":
			content := r.Form.Get("twitte")
			temp := &pb.TIn{Username:username.Name, Content:content}
			c.SendTwitte(ctx, temp)
		case "follow":
			follow := r.Form.Get("follow")
			temp := &pb.FollowUnfollow{Username:username.Name, Other:follow}
			c.FollowUser(ctx, temp)
		case "unfollow":
			unfollow := r.Form.Get("unfollow")
			temp := &pb.FollowUnfollow{Username:username.Name, Other:unfollow}
			c.UnfollowUser(ctx, temp)
		}
		http.Redirect(w, r, "/personalPage", http.StatusFound)
	}
	defer cancel()
}