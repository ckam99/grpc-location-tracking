package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"grpc-location/docs"
	v1 "grpc-location/proto/gen/go"
	"log"
	"net"
	"net/http"
)

type server struct {
	v1.UnimplementedGreeterServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	return &v1.HelloReply{Message: in.Name + " world"}, nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8880")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	srv := &server{}
	v1.RegisterGreeterServer(s, srv)
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8880")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.NewClient(
		"0.0.0.0:8880",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	muxJsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	rmux := runtime.NewServeMux(muxJsonOptions)
	// Register Greeter
	err = v1.RegisterGreeterHandler(context.Background(), rmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	//gwServer := &http.Server{
	//	Addr:    ":8890",
	//	Handler: rmux,
	//}
	//log.Println("Serving gRPC-Gateway on http://0.0.0.0:8890")
	//log.Fatalln(gwServer.ListenAndServe())

	server := http.NewServeMux()
	fs := http.FileServer(http.Dir("./docs/swagger"))
	server.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	server.Handle("/", rmux)
	server.Handle("/docs/swagger", docs.SwaggerHandler)
	listener, err := net.Listen("tcp", ":8899")
	if err != nil {
		log.Fatalf("cannot create http network listener:%w", err)
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8899")
	if err = http.Serve(listener, server); err != nil {
		log.Fatalf("cannot start HTTP server: %w", err)
	}
	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8899")
}

//
//func (s *server) ServeHttpGateway(address string) error {
//	muxJsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
//		MarshalOptions: protojson.MarshalOptions{
//			UseProtoNames:   true,
//			EmitUnpopulated: true,
//		},
//		UnmarshalOptions: protojson.UnmarshalOptions{
//			DiscardUnknown: true,
//		},
//	})
//	rmux := runtime.NewServeMux(muxJsonOptions)
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	if err := v1.RegisterGreeterHandlerServer(ctx, rmux, s); err != nil {
//		return fmt.Errorf("cannot register book handler server: %s", err)
//	}
//
//	server := http.NewServeMux()
//	fs := http.FileServer(http.Dir("./docs/swagger"))
//	server.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
//	server.Handle("/", rmux)
//
//	listener, err := net.Listen("tcp", address)
//	if err != nil {
//		return fmt.Errorf("cannot create http network listener:%w", err)
//	}
//	if err = http.Serve(listener, server); err != nil {
//		return fmt.Errorf("cannot start HTTP server: %w", err)
//	}
//
//	//
//	//// http server
//	//mux := NewHTTPServer()
//	//hdle := HttpLogger(rmux)
//	//mux.Handle("/", hdle)
//	//if err := mux.Serve(address); err != nil {
//	//	return fmt.Errorf("grpc http gateway :%w", err)
//	//}
//	return nil
//}
//
//type ResponseHanlder struct {
//	http.ResponseWriter
//	StatusCode int
//	Body       []byte
//}
//
//func (r *ResponseHanlder) WriteHeader(statusCode int) {
//	r.StatusCode = statusCode
//	r.ResponseWriter.WriteHeader(statusCode)
//}
//
//func (r *ResponseHanlder) Write(b []byte) (int, error) {
//	r.Body = b
//	return r.ResponseWriter.Write(b)
//}
//
//func HttpLogger(hdle http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		duration := time.Since(time.Now())
//
//		rec := &ResponseHanlder{
//			ResponseWriter: w,
//			StatusCode:     http.StatusOK,
//		}
//		hdle.ServeHTTP(rec, r)
//
//		if rec.StatusCode != http.StatusOK {
//			log.Println(duration, rec.Body)
//		}
//
//	})
//}
