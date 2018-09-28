package main

import (
	"io"
	"log"
	"os"
	"time"
	// "github.com/pkg/errors"
	pb "github.com/kumaya/grpc_stream_upload/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:8003"
	chunksize   = 200
	defaultFile = "./main.go"
)

func UploadFile(client pb.GRPCStreamUploadServiceClient, ctx context.Context, f string) (err error) {
	log.Printf("File send start at: %v", time.Now())
	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
		status  *pb.FileUploadAck
	)

	file, err = os.Open(f)
	if err != nil {
		log.Fatalf("failed to open file: %s", f)
		// err = errors.Wrapf(err, "failed to open file %s", f)
		return
	}
	defer file.Close()

	stream, err := client.SendFile(ctx)
	if err != nil {
		log.Fatalf("failed to create upload stream for %s", f)
		// err = errors.Wrapf(err, "failed to create upload stream for %s", f)
		return
	}
	defer stream.CloseSend()

	buf = make([]byte, chunksize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}
			log.Fatalf("errored while copying from file to buffer: %v", err)
			// err = errors.Wrapf(err, "errored while copying from file to buffer")
			return
		}

		err = stream.Send(&pb.FileChunk{
			Content: buf[:n],
		})
		if err != nil {
			log.Fatalf("failed to send chunk: %v", err)
			// err = errors.Wrapf(err, "failed to send chunk")
			return
		}
	}

	log.Printf("File send complete at: %v", time.Now())

	status, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to receive response: %v", err)
		// err = errors.Wrapf(err, "failed to receive response")
		return
	}

	if status.Code != pb.FileUploadStatusCode_Ok {
		log.Fatalf("upload failed: %s", status.Message)
		// err = errors.Errorf("upload failed: %s", status.Message)
		return
	} else {
		log.Printf("Message: %s, Code: %s", status.Message, status.Code)
	}
	return
}

func main() {
	// Setup connection to the server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to start grpc connection: %v", err)
		// errors.Wrapf(err, "failed to start grpc connection with %s", address)
	}
	defer conn.Close()

	c := pb.NewGRPCStreamUploadServiceClient(conn)

	file := defaultFile
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	// log.Printf("%s", os.Args)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = UploadFile(c, ctx, file)
	if err != nil {
		log.Fatalf("could not call upload: %v", err)
		// err = errors.Wrapf(err, "could not call uplaod")
	}
	return
}
