package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"path"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// OwnerGRPC implements EntryServiceServer and delegates the calls to EntryMan
type OwnerGRPC struct {
}

// FilesGRPC ...
type FilesGRPC struct {
}

// InteractionsGRPC ...
type InteractionsGRPC struct {
}

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
			grpclog.Info("peer certs: ", info.State.PeerCertificates)
			if len(info.State.PeerCertificates) > 0 {
				grpclog.Info("peer cert cn: ", info.State.PeerCertificates[0].Subject.CommonName)
			}
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
	grpclog.Infof("invoke unary method=%s duration=%s error=%v", info.FullMethod, time.Since(start), err)
	return resp, err
}

// TODO serverStreamInterceptor

// Create ...
func (es *OwnerGRPC) Create(ctx context.Context, request *Entry) (*Response, error) {
	return entryMan.create(request)
}

// Update ...
func (es *OwnerGRPC) Update(ctx context.Context, request *Entry) (*Response, error) {
	return entryMan.update(request)
}

// Query ...
func (es *OwnerGRPC) Query(ctx context.Context, filter *Filter) (*Response, error) {
	return entryMan.query(filter)
}

// Delete ...
func (es *OwnerGRPC) Delete(ctx context.Context, request *Entry) (*Response, error) {
	return entryMan.delete(request)
}

func logSleep(ctx context.Context, d time.Duration) {
	if tr, ok := trace.FromContext(ctx); ok {
		tr.LazyPrintf("sleeping for %s", d)
	}
}

// Notifications ...
func (es *OwnerGRPC) Notifications(request *Filter, stream Owner_NotificationsServer) (err error) {
	// TODO establish per-call (user/call notification channel)
	// TODO handle cancelation
	ctx := stream.Context()
	for i := 0; i < 10; i++ {
		d := time.Duration(rand.Intn(3)) * time.Second
		logSleep(ctx, d)
		select {
		case <-time.After(d):
			err := stream.Send(&Notification{
				What:      fmt.Sprintf("result %d for [%s] from backend %d", i, request, d),
				Timestamp: uint64(i),
			})
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	/*
		for j := 0; j < 10; j++ {

			if err := stream.Send(&Notification{}); err != nil {
				return err
			}
		}

	*/
	return nil

}

// DownloadServers ...
func (ag *FilesGRPC) DownloadServers(ask *ChunkAsk, stream Files_DownloadServersServer) (err error) {
	return
}

// List ...
/*func (ag *FilesGRPC) List(*ListFilter, Files_ListServer) (err error) {
	return
}*/

// Download ...
func (ag *FilesGRPC) Download(Files_DownloadServer) (err error) {
	return
}

// Upload ...
func (ag *FilesGRPC) Upload(Files_UploadServer) (err error) {
	return
}

// SendMessage ...
func (ui *InteractionsGRPC) SendMessage(ctx context.Context, message *Message) (receipt *Response, err error) {
	return
}

// Notifications ...
func (ui *InteractionsGRPC) Notifications(filter *Filter, stream Interactions_NotificationsServer) (err error) {
	return
}

// Query ...
func (ui *InteractionsGRPC) Query(ctx context.Context, filter *Filter) (response *Response, err error) {
	return
}

// React ...
func (ui *InteractionsGRPC) React(ctx context.Context, reactionRequest *ReactionRequest) (receipt *Response, err error) {
	return
}

// Share ...
func (ui *InteractionsGRPC) Share(ctx context.Context, shareRequest *ShareRequest) (receipt *Response, err error) {
	return
}

// MakeComment ...
func (ui *InteractionsGRPC) MakeComment(ctx context.Context, commentRequest *CommentRequest) (receipt *Response, err error) {
	return
}

// MissedCalls ...
func (ui *InteractionsGRPC) MissedCalls(stream Interactions_MissedCallsServer) (err error) {
	return
}

// Resend ...
func (ui *InteractionsGRPC) Resend(context.Context, *Empty) (response *Response, err error) {
	return
}

// Run ...
func runGRPC() {
	listen, err := net.Listen("tcp", "localhost:50051")
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
		ClientAuth:   tls.VerifyClientCertIfGiven, //tls.RequireAndVerifyClientCert,
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

	RegisterOwnerServer(grpcServer, ownerGrpc)
	RegisterFilesServer(grpcServer, filesGrpc)
	//RegisterAdminServer(grpcServer, adminGrpc)
	RegisterInteractionsServer(grpcServer, interactionsGrpc)

	grpclog.Info("Starting the gRPC server")
	grpcServer.Serve(listen)
}