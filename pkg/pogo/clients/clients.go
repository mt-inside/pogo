package clients

import (
	"log"

	pb "github.com/mt-inside/pogo/api"
	"google.golang.org/grpc"
)

const (
	serverAddr string = "localhost:50001"
)

var (
	conn *grpc.ClientConn
)

func init() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	c, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("Couldn't connect to pogo server: %v", err)
	} else {
		conn = c
	}
	//defer conn.Close() FIXME
}

func GetPogoClient() pb.PogoClient {
	return pb.NewPogoClient(conn)
}
func GetTasksClient() pb.TasksClient {
	return pb.NewTasksClient(conn)
}
