package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// EntryClient ...
type EntryClient struct {
	host      string
	port      int
	certsPath string
	conn      *grpc.ClientConn
	service   OwnerClient
}

// Conn ...
func Conn(certsPath string, host string, port int) (c *grpc.ClientConn) {
	var err error
	grpc.EnableTracing = true
	certificate, err := tls.LoadX509KeyPair(
		path.Join(certsPath, "kefah.crt"),
		path.Join(certsPath, "kefah.key"),
	)

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(path.Join(certsPath, "edrajRootCA.crt"))
	if err != nil {
		log.Fatalf("failed to read ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		log.Fatal("failed to append certs")
	}

	transportCreds := credentials.NewTLS(&tls.Config{
		ServerName:   host,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	// TODO : additionally consider stats: grpc.WithStatsHandler(th)
	c, err = grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithUnaryInterceptor(clientUnaryInterceptor),
		grpc.WithStreamInterceptor(clientStreamInterceptor))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return c
}

func check(response *Response, err error) {

	if err != nil {
		log.Printf("call Failed: %v", err)
	} else {
		log.Printf("Response: %v", response)
		/*
			switch r := response.(type) {
			case *Response:
				log.Printf("Response: %v", r)
			case *Receipt:
				log.Printf("Receipt: %v", r)
			}*/
	}
}

func printReturnedMeta(meta ...metadata.MD) {
	for i, one := range meta {
		for key, value := range one {
			fmt.Printf("[%d] %s => %s\n", i, key, value)
		}
	}
}

func clientUnaryInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	headers := metadata.New(
		map[string]string{
			"edraj-signature-bin": "mysig",
			"edraj-pubkey-bin":    "mykey",
			"edraj-id":            "myid",
			"edraj-timestamp":     "mytime"})

	//log.Println(headers)

	// https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-metadata.md
	// this is the critical step that includes your headers
	ctx = metadata.NewOutgoingContext(ctx, headers)
	var header, trailer metadata.MD

	opts = append(opts, grpc.Header(&header))
	opts = append(opts, grpc.Trailer(&trailer))

	err := invoker(ctx, method, req, reply, cc, opts...) // <==
	printReturnedMeta(header, trailer)
	log.Printf("invoke remote method=%s duration=%s error=%v", method, time.Since(start), err)
	return err
}

func clientStreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (client grpc.ClientStream, err error) {
	start := time.Now()
	client, err = streamer(ctx, desc, cc, method, opts...)
	log.Printf("invoke remote stream method=%s duration=%s error=%v", method, time.Since(start), err)
	return
}

// TODO clientStreamInterceptor

func main() {
	conn := Conn("../../../workspace/certs/", "localhost", 50050)
	defer conn.Close()
	//health := NewHealthClient(conn)
	//r, err := health.Check(context.Background(), &HealthCheckRequest{Service: "Hi there"})
	//if err != nil {
	//	panic(err)
	//}

	//log.Println("Got this reponse: ", r)

	/*
		// Set up a connection to the server.
		owner := NewOwnerClient(conn)

		ctx := context.Background()
		one := Content{Id: "one", Pathname: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}
		two := Content{Id: "two", Pathname: "/home", Shortname: "Ali", Tags: []string{"Aee", "Bee", "Cee"}}
		check(owner.Delete(ctx, &Entry{Type: EntryType_CONTENT, Id: one.Id}))

		check(owner.Delete(ctx, &Entry{Type: EntryType_CONTENT, Id: two.Id}))

		check(owner.Create(ctx, &Entry{Type: EntryType_CONTENT, Content: &one}))
		check(owner.Create(ctx, &Entry{Type: EntryType_CONTENT, Content: &two}))

		check(owner.Query(ctx, &Filter{EntryType: EntryType_CONTENT}))
		//check(owner.Get(ctx, &Entry{Type: EntryType_CONTENT, Id: one.Id}))

		check(owner.Delete(ctx, &Entry{Type: EntryType_CONTENT, Id: one.Id}))
		check(owner.Delete(ctx, &Entry{Type: EntryType_CONTENT, Id: two.Id}))
	*/
	/*
		stream, err := owner.Notifications(ctx, &Filter{})
		if err != nil {
			log.Println("Error on streaming", err)
			return
		}
		for {
			notification, err := stream.Recv()
			if err == io.EOF {
				log.Println("Notifications stream ends here")
				break
			}
			if err != nil {
				log.Fatalf("%v.Notifications(_) = _, %v", client, err)
			}
			log.Println(notification)
		}*/
}
