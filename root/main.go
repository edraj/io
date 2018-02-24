package main

import (
	"fmt"
	"time"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Config options
type Config struct {
	assetsPath     string
	certsPath      string
	templatesPath  string
	servingAddress string
	mongoAddress   string
	dataPath       string
	shutdownWait   time.Duration
}

var (
	// Config optoins TODO read from config file and/or params
	config = Config{
		assetsPath:     "./assets/",
		certsPath:      "../../workspace/certs/",
		templatesPath:  "./templates/",
		dataPath:       "./data", // Mongo, files, indexes ...
		servingAddress: "127.0.0.1:5533",
		shutdownWait:   15 * time.Second,
		mongoAddress:   "127.0.0.1:27017",
	}

	rootGrpc = new(RootGRPC)
	//entryMan   = &EntryMan{}
	grpcServer *grpc.Server
)

// RootGRPC ...
type RootGRPC struct {
}

// Issue ...
func (r *RootGRPC) Issue(ctx context.Context, actor *Actor) (response *Response, err error) {
	return
}

// Revoke ...
func (r *RootGRPC) Revoke(ctx context.Context, actor *Actor) (response *Response, err error) {
	return
}

// List ...
func (r *RootGRPC) List(*Query, Root_ListServer) (err error) {
	return

}

// IssueAnother ...
func (r *RootGRPC) IssueAnother(ctx context.Context, actor *Actor) (certificate *Certificate, err error) {
	return
}

// Ping ...
func (r *RootGRPC) Ping(ctx context.Context, empty *Empty) (response *Response, err error) {
	return
}

// Get details on other Actors
func (r *RootGRPC) Get(ctx context.Context, actor *Actor) (response *Actor, err error) {
	return

}

func main() {
	fmt.Println("Hello there")
}
