all:
	protoc --gogofaster_out=plugins=grpc:. ./protos/filestream.proto
	cd server; go build
	cd client; go build

clean:
	cd server; go clean
	cd client; go clean