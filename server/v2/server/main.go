package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"

	codes "google.golang.org/grpc/codes"

	status "google.golang.org/grpc/status"

	pickupv2 "github.com/heartziq/recycle-lah/server/v2/proto/pickup"

	emptypb "google.golang.org/protobuf/types/known/emptypb"

	pickup "github.com/heartziq/recycle-lah/server/handlers"

	_ "github.com/go-sql-driver/mysql"
)

const (
	GET_ALL_PICKUPS_QUERY = `
	SELECT 
	id, ST_X(coord) as lat, ST_Y(coord) as lng,
	address,
	created_by, attend_by,
	completed
	
	FROM your_db.pickups
	WHERE attend_by='';`
)

type server struct {
	pickupv2.UnimplementedPickupsServer
	db *sql.DB
}

func CreateGPRCProxyServer() *server {
	// create db connect
	newDB, err := sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/your_db")
	if err != nil {
		panic(err)
	}
	return &server{db: newDB}
}

func (s *server) ListPickups(_ *emptypb.Empty, stream pickupv2.Pickups_ListPickupsServer) error {

	// access db
	results, err := s.db.Query(GET_ALL_PICKUPS_QUERY)

	if err != nil {
		panic(err.Error())

	}

	responseChannel := make(chan *httpbody.HttpBody, 12)

	go func() {
		defer close(responseChannel)
		for results.Next() {
			// map this type to the record in the table
			c := pickup.Pickup{}

			err = results.Scan(
				&c.Id,
				&c.Lat, &c.Lng,
				&c.Address, &c.CreatedBy,
				&c.Collector, &c.Completed,
			)

			if err != nil {

				panic(err.Error())

			}

			b, err := json.Marshal(&c)
			if err != nil {
				status.Errorf(codes.Aborted, "error: %v", err)
				return

			}

			responseChannel <- &httpbody.HttpBody{
				ContentType: "application/json",
				Data:        b,
			}

		}
	}()
	for response := range responseChannel {
		if err := stream.Send(response); err != nil {

			return err
		} else {
			log.Println("sent!")
		}
	}

	return nil
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	pickupv2.RegisterPickupsServer(s, CreateGPRCProxyServer())
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()

	// Register Greeter
	err = pickupv2.RegisterPickupsHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
