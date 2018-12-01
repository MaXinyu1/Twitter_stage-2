package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sort"
	"strings"
	"sync"
	"time"
	"context"
	pb "google.golang.org/grpc/examples/twitter/proto"
)

const (
	userName = "root"
	password = ""
	ip = "127.0.0.1"
	database_port = "3306"
	dbName = "twitter"

	server_port = ":50051"
)

var DB = &sql.DB{}

type TwitterActionServer struct{
	savedname *pb.Username
	saveduser *pb.User
	m      sync.RWMutex
}

type T_in struct {
    username string
	content string
}

type Twitte struct {
	message T_in
	time string
}

type Twitlist []Twitte

// Sort Function needed these three Function
func (I Twitlist) Len() int {
	return len(I)
}
func (I Twitlist) Less(i, j int) bool {
	return I[i].time > I[j].time
}
func (I Twitlist) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//database
func DBstart() {
	path := strings.Join([]string{userName, ":", password, "@tcp(",ip, ":", database_port, ")/", dbName, "?charset=utf8"}, "")
	DB,_ = sql.Open("mysql", path);
	_,err := sql.Open("mysql", path);
	if err != nil {
		fmt.Printf("connect mysql failed! [%s]", err)
		return
	}

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
		return
	}
}

func showFollow(username string) []string{
	DBstart()
	rows,_ := DB.Query("select toU from follow where fromU = ?", username)
	DB.Close()
	var following []string
	for rows.Next(){
		var who string
		rows.Scan(&who)
		following = append(following, who)
	}
	return following
}

func showUnfollow(username string) []string{
	DBstart()
	rows,_ := DB.Query("select username from user where username not in (select toU from follow where fromU = ?)", username)
	DB.Close()
	var unfollowing []string
	for rows.Next(){
		var unfollow string
		rows.Scan(&unfollow)
		//fmt.Println(unfollow)
		unfollowing = append(unfollowing, unfollow)
	}
	return unfollowing
}

func showPost(follows []string) Twitlist{ //good
	twitlist := Twitlist{}
	DBstart()
	for i := range follows {
		name := follows[i]
		rows,_ := DB.Query("select * from twitte where username = ?", name)
		for rows.Next(){
			var username, content, time string
			rows.Scan(&username, &content, &time)
			t := Twitte{
				time:time,
				message : T_in{username:username,content:content},
			}
			twitlist = append(twitlist, t)
			/*
			twitlist = &pb.Twitlist{
				Twitlists: []*pb.Twitte{
					t,
				},
			}*/
		}
	}
	DB.Close()
	return twitlist
}

func deletes(Following []string, username string) []string{
	var res []string
	for _,n := range Following{
		if n != username {
			res = append(res, n)
		}
	}
	return res
}

func getContent (Posts Twitlist) []string {
	var res []string
	var temp string
	for _,n := range Posts{
		temp = n.message.username + ":     >" + n.message.content
		//fmt.Println(temp)
		res = append(res, temp)
	}
	return res
}

//need method: GetTwitterPage(context.Context, *Ack) (*TwitterPage, error)
//have method: GetTwitterPage(ctx context.Context, a *pb.Ack) *pb.TwitterPage
func (s *TwitterActionServer) GetTwitterPage(ctx context.Context, a *pb.Username) (*pb.TwitterPage, error) { //u is unused
	DBstart()

	//Read and
	s.m.RLock()
	defer s.m.RUnlock()

	username := a.Name
	Following := showFollow(username)
	UnFollowed := showUnfollow(username)
	Posts := showPost(Following)


	sort.Sort(Posts)

	//transfer Posts to string form
	newPosts := getContent(Posts)

	// Remove the user itself from following list (just not shown in screen but in memory)
	Following = deletes(Following, username)

	DB.Close()
	return &pb.TwitterPage{Username: username, Following: Following, UnFollowed: UnFollowed, Posts: newPosts}, nil

}

//login
// need method: LoginCheck(context.Context, *User) (*IsTrue, error)
// have method: LoginCheck(ctx context.Context, user *pb.User) *pb.IsTrue
func (s *TwitterActionServer) LoginCheck(ctx context.Context, user *pb.User) (*pb.IsTrue, error) {
	DBstart()

	s.m.RLock()
	defer s.m.RUnlock()

	Username := user.UserName
	Password := user.PassWord
	rows, err := DB.Query("SELECT password FROM user where username = ?", Username);
	DB.Close()
	if err != nil {
		//fmt.Println(err)
		return &pb.IsTrue{IsTrue:false}, err
	} else {
		for rows.Next() {
			password := ""
			rows.Scan(&password)
			if Password == password {

				//store the data into the server, doing session jobs
				t1 := &pb.Username{Name:Username}
				t2 := &pb.User{UserName:Username,PassWord:Password}
				s = &TwitterActionServer{savedname:t1, saveduser:t2}
				return &pb.IsTrue{IsTrue:true}, nil
			}
		}
	}
	return &pb.IsTrue{IsTrue:false}, nil
}

func (s *TwitterActionServer) RegisterCheck(ctx context.Context, user *pb.User) (*pb.IsTrue, error) {
	DBstart()

	s.m.RLock()
	defer s.m.RUnlock()

	username := user.UserName
	password := user.PassWord
	rows,err := DB.Query("select username from user")

	if err != nil {
		fmt.Println(err)
		return &pb.IsTrue{IsTrue:false}, err
	} else {
		name := ""
		for rows.Next() {
			rows.Scan(&name)
			if name == username {
				return &pb.IsTrue{IsTrue:false}, nil
			}
		}
	}

	//if the user doesn't exit then add it to the database
	_,err = DB.Exec("insert into user (username, password) values (?, ?)", username, password)
	//add him/herself to friend list
	_,err = DB.Exec("insert into follow (fromU, toU) values (?, ?)", username, username)
	//if err != nil {fmt.Println(err)} else {fmt.Println("insert succcess !")}
	return &pb.IsTrue{IsTrue:true}, nil
}

//action
func (s *TwitterActionServer) SendTwitte(ctx context.Context, data *pb.TIn) (*pb.IsTrue, error){
	DBstart()

	//write lock
	s.m.Lock()
	defer s.m.Unlock()

	username := data.Username
	content := data.Content
	t1 := time.Now().Year()
	t2 := time.Now().Month()
	t3 := time.Now().Day()
	t4 := time.Now().Hour()
	t5 := time.Now().Minute()
	t6 := time.Now().Second()
	t7 := time.Now().Nanosecond()
	t := time.Date(t1, t2, t3, t4, t5, t6, t7, time.Local)
	fmt.Println(t)
	_,err := DB.Exec("insert into twitte (username, content, time) values (?, ?, ?)", username, content, t)
	DB.Close();
	if err == nil {
		return &pb.IsTrue {IsTrue:true}, nil
	}
	return &pb.IsTrue {IsTrue:false}, nil
}

func (s *TwitterActionServer) FollowUser(ctx context.Context, data *pb.FollowUnfollow) (*pb.IsTrue, error){
	DBstart()

	username := data.Username
	follow := data.Other
	_,err := DB.Query("insert into follow (fromU, toU) values (?, ?)", username, follow)
	DB.Close()
	if err == nil {
		return &pb.IsTrue {IsTrue:true}, nil
	}
	return &pb.IsTrue {IsTrue:false}, nil
}

func (s *TwitterActionServer) UnfollowUser(ctx context.Context, data *pb.FollowUnfollow) (*pb.IsTrue, error){
	DBstart()

	username := data.Username
	unfollowname := data.Other
	_,err := DB.Query("delete from follow where fromU = ? and toU = ?", username, unfollowname)
	DB.Close()
	if err == nil {
		return &pb.IsTrue {IsTrue:true}, nil
	}
	return &pb.IsTrue {IsTrue:false}, nil
}

func (s *TwitterActionServer) GetName (ctx context.Context, ack *pb.Ack) (*pb.Username, error){
	return s.savedname, nil
}

func main() {
	lis, err := net.Listen("tcp", server_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()//create a new server
	fmt.Println("create server success")
	pb.RegisterTwitterActionServer(s, &TwitterActionServer{}) // register the server
	fmt.Println("Register server success")
	reflection.Register(s) // Register reflection service on gRPC server.
	if err := s.Serve(lis); err != nil { //start the server
		log.Fatalf("failed to serve: %v", err)
	}
}

