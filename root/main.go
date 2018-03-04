package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strings"
	"time"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	consullib "github.com/hashicorp/consul/api"
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
type RootGRPC struct{}

// Check Health
func (r *RootGRPC) Check(ctx context.Context, request *HealthCheckRequest) (response *HealthCheckResponse, err error) {
	grpclog.Info("Got Healthcheck request: ", request.Service)
	response = &HealthCheckResponse{Status: HealthCheckResponse_SERVING}
	err = status.Errorf(codes.OK, "We are good")
	return response, err
}

// Issue ...
func (r *RootGRPC) Issue(ctx context.Context, actor *Actor) (response *Response, err error) { return }

// Revoke ...
func (r *RootGRPC) Revoke(ctx context.Context, actor *Actor) (response *Response, err error) { return }

// List ...
func (r *RootGRPC) List(*Filter, Root_ListServer) (err error) { return }

// IssueAnother ...
func (r *RootGRPC) IssueAnother(ctx context.Context, actor *Actor) (certificate *Certificate, err error) {
	return
}

// Ping ...
func (r *RootGRPC) Ping(ctx context.Context, empty *Empty) (response *Response, err error) { return }

// Get details on other Actors
// Get(ctx context.Context, actor *Actor) (response *Actor, err error) { return }

//Issue(context.Context, *Actor) (*Response, error)
// Revoke (RootAdmin only) Equivelent to De-activate
//Revoke(context.Context, *Actor) (*Response, error)
// List actors (RootAdmin only)
//List(*Filter, Root_ListServer) error
// IssueAnother (Actor) Produce CRT from CSR for the same commonName.
// i.e. multiple devices / server cluster deployments.
//IssueAnother(context.Context, *Actor) (*Certificate, error)
// Ping update the ip of the domain (server). IP contains the TTL,
// and a flag to say whether to delete existing ips or add to them

// Beacon ...
func (r *RootGRPC) Beacon(ctx context.Context, _ *Empty) (response *Response, err error) { return }

// Query details on other Actors
// Return details on another actor. for domains (servers) the ip(s) are returned.
func (r *RootGRPC) Query(ctx context.Context, actor *Actor) (ractor *Actor, err error) { return }

// MissedCall informs Root of server A's attempt to reach server B. Server B will be notified of all AwayCalls the moment they show beacon of life.
func (r *RootGRPC) MissedCall(ctx context.Context, domain *Domain) (response *Response, err error) {
	return
}

// GoingOffline ...
func (r *RootGRPC) GoingOffline(ctx context.Context, _ *Empty) (response *Response, err error) { return }

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	start := time.Now()
	//newStream := grpc_middleware.WrapServerStream(stream)
	//newStream.WrappedContext = context.WithValue(ctx, "user_id", "john@example.com")
	err = handler(srv, stream)
	grpclog.Infof("invoke stream method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return
}

//grpc.StreamServerInterceptor()
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range headers {
			if strings.HasPrefix(k, "edraj-") {
				grpclog.Infof("Ctx %v: %v", k, v)
			}
		}
	}

	var username string
	pr, ok := peer.FromContext(ctx)
	if ok {
		switch info := pr.AuthInfo.(type) {
		case credentials.TLSInfo:
			if len(info.State.VerifiedChains) > 0 && len(info.State.VerifiedChains[0]) > 0 {
				username = info.State.VerifiedChains[0][0].Subject.CommonName
			}
			//grpclog.Info("peer certs: ", info.State.PeerCertificates)
			//if len(info.State.PeerCertificates) > 0 {
			//grpclog.Info("peer cert cn: ", info.State.PeerCertificates[0].Subject.CommonName)
			//}
			//default:
			//return nil, status.Error(codes.Unauthenticated, "Unknown AuthInfo type")
		}
	}

	grpclog.Info("Username: ", username)

	grpc.SendHeader(ctx, metadata.New(map[string]string{"edraj-header": "my-value"}))
	grpc.SetTrailer(ctx, metadata.New(map[string]string{"edraj-trailer": "my-value"}))
	//ctx = context.WithValue(ctx, "user_id", "john@example.com")
	start := time.Now()
	resp, err := handler(ctx, req)
	grpclog.Infof("Unary=%s took=%s error=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

// GrpcService ...
func GrpcService() {
	listen, err := net.Listen("tcp", "localhost:50050")
	if err != nil {
		grpclog.Fatal(err)
	}

	certificate, err := tls.LoadX509KeyPair(path.Join(config.certsPath, "localhost.crt"), path.Join(config.certsPath, "localhost.key"))

	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(path.Join(config.certsPath, "edrajRootCA.crt"))
	if err != nil {
		grpclog.Fatalf("failed to read client ca cert: %s", err)
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		grpclog.Fatal("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert, //tls.NoClientCert, //tls.RequireAnyClientCert, //tls.VerifyClientCertIfGiven, //tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}

	// TODO : Additionally consider grpc.StatsHandler(th)
	grpcServer = grpc.NewServer(
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.MaxConcurrentStreams(64),
		//grpc.InTapHandle(NewTap.Handler),
	)

	RegisterRootServer(grpcServer, rootGrpc)
	RegisterHealthServer(grpcServer, rootGrpc)

	grpclog.Info("Starting the gRPC server")
	grpcServer.Serve(listen)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	grpc.EnableTracing = true
	grpclog.SetLogger(log.New(os.Stdout, "edraj: ", log.LstdFlags))

	consulConfig := consullib.DefaultConfig()
	consulConfig.Address = "127.0.0.1:8500"
	consul, err := consullib.NewClient(consulConfig)
	check(err)
	kv := consul.KV()
	d := &consullib.KVPair{Key: "sites/1/domain", Value: []byte("mydomain.com")}
	wm, err := kv.Put(d, nil)
	check(err)
	grpclog.Info("Put returned ", wm)
	kvp, qm, err := kv.Get("sites/1/domain", nil)
	check(err)
	grpclog.Info("Get returned ", kvp, qm)
	grpclog.Info("Value: ", string(kvp.Value))
	grpclog.Info("Hello there")

	GrpcService()
}
