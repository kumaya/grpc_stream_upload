package main

import (
	"io"
    "log"
    "net"
    "time"
    // "github.com/pkg/errors"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    pb "github.com/kumaya/grpc_stream_upload/protos"
)

const (
	port = ":8003"
)

type server struct{}


func (s *server) SendFile(stream pb.GRPCStreamUploadService_SendFileServer) (err error) {
	log.Printf("File receive start at: %v", time.Now())
	var strm *pb.FileChunk
	for {
		strm, err = stream.Recv()
		log.Printf("===================================")
		log.Printf("%v", strm)
		if err != nil {
			if err == io.EOF {
				goto END
			}
			log.Fatalf("failed while reading chunks: %v", err)
			// err = errors.Wrapf(err, "failed while reading chunks")
			return
		}
	}
	
	// GOTO statement begins
	END:
	log.Printf("File received complete at: %v", time.Now())

	err = stream.SendAndClose(&pb.FileUploadAck{
					Message: "File received successfully, YooHoo!!",
					Code: pb.FileUploadStatusCode_Ok,
		})
	if err != nil {
		log.Fatalf("failed to send status code: %v", err)
		// err = errors.Wrapf(err, "failed to send status code")
		return
	}

	return
}


func main() {
	log.Printf("Server started. Listening on port %s", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		// errors.Wrapf(err, "failed to listen")
	}
	s := grpc.NewServer()
	pb.RegisterGRPCStreamUploadServiceServer(s, &server{})

	reflection.Register(s)

	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
		// errors.Wrapf(err, "failed to serve")
	}
	return
}
