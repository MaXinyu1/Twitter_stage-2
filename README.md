# Twitter_stage-2
2018 Fall Distributed System (stage 2 - gRPC)

## Description
This project implement a twitter with simple functions.  
Built for 2018 fall distributed system course.

It is divided into 3 stages:
- [ ] Build simple web application with database
- [x] Split off backend into a seperate service (using gPRC)
- [ ] Bind the service with a distributed system
This is the second stage of the project. Different from the stage 1, stage 2 put all the data processing functions into the server. So that client don't have to process any data locally.

## Main Features
1. Creating an account, with username and password
2. Logging in as a given user, given username and password
3. Users can follow other users (reversible).
4. Users can create posts that are associated with their identity.
5. Users can view some feed composed only of content generated by users they follow.
6. Recommend users who they can follow

## Instructions To Run
**1. Install thrid-party packages**   
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;*go get github.com/gorilla/securecookie*

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;*go get github.com/go-sql-driver/mysql/*

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;*go get github.com/golang/protobuf*

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;*go install github.com/golang/protobuf/protoc-gen-go/*

**2. Clone the project into "/your/path"**  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;*git clone ...*  

**3. Go into the src directory and run it**  
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;*cd /your/path/awesomeProject*  

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;*go run web.go*

## Project Structure
```bash
├── README.md
├── github.com
│   └── go-sql-driver
├── google.golang.org
├── proto
│   ├── *action.proto
│   └── *action.pb.go
└── gRPC
    ├── server       // go run twitter_server.go
    │   └── *twitter_server.go
    ├── client       // go run twitter_server.go
    │   ├── *client.go
    │   └── auth
    │       └── *auth.go
    └── show         // html and css templates
        ├── *.html
        ├── css
        └── js
```

## Team
- Xinyu Ma (xm546)
