package main

//go:generate protoc -I ../../api --go_out=plugins=grpc:../../api ../../api/pogo.proto

import (
	"log"
	"net"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	pb "github.com/mt-inside/pogo/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port string = ":50001"
)

func main() {
	initConfig()

	sock, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterPogoServer(srv, &PogoServer{})
	pb.RegisterTasksServer(srv, &TasksServer{})
	// Turn on reflection so that clients can dynamically query our services
	reflection.Register(srv)
	log.Printf("serving on %v", port)
	if err := srv.Serve(sock); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func initConfig() {
	/* Config sources, in viper presidence order */

	/* Defaults */
	viper.SetDefault("pomodoro_time", 25)

	/* Config file */
	viper.SetConfigName("config") /* e.g. config.yaml, config.json */
	viper.AddConfigPath("$HOME/.pogo/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		//Also goes off for no config file
		//log.Fatalf("Fatal error in config file: %s \n", err)
	}
	/* ...with auto-reload */
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		viper.Debug()
	})

	/* Environment */
	viper.SetEnvPrefix("pogo")
	viper.AutomaticEnv() /* e.g. POGO_POMODORO_TIME */

	/* Command-line args */
	/* Call pflag direct because we're not using Cobra */
	pflag.Int("pomodoro_time", 25, "Length of a pomodoro, in minutes") /* i.e. --pomodoro_time */
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
}
