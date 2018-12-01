package page

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"context"
	"net/http"
	"html/template"
	pb "google.golang.org/grpc/examples/twitter/proto"
	"time"

	//"awesomeProject/cookie"
	//"sort"
)

const (
	address = "localhost:50051"
)

func PersonalPage(w http.ResponseWriter, r *http.Request){
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTwitterActionClient(conn)

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	username,_ := c.GetName(ctx, &pb.Ack{})
	if r.Method == "GET" {
		t, _ := template.ParseFiles("show/personalPage.html")
		ack := &pb.Ack{}
		pagecontent,_ := c.GetTwitterPage(ctx, ack)
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
