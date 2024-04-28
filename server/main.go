package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"service/backup/databases/proto"
	"service/backup/databases/server/config"
	"service/backup/databases/server/controllers"
	servicegrpc "service/backup/databases/server/service_grpc"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func InitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Warn("Cannot load env file, using system env")
	}
}

func runFiberServer() {
	InitEnv()
	config.OpenDB()

	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024, //50mb
	})

	// controllers.RouteCars(app)
	controllers.RouteBckpDatabase(app)

	err := app.Listen(":3000")
	if err != nil {
		logrus.Fatal(
			"Error on running fiber, ",
			err.Error())
	}

	_ = os.MkdirAll("upload/", 0777)
}

func runGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterFileUploadServer(s, &servicegrpc.FileUpload{})
	log.Println("Server started at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	_ = os.Mkdir("upload", 0777)

	go runGRPCServer()

	go runFiberServer()

	<-stop
	log.Println("Shutting down servers...")
}
