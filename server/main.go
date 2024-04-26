package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"service/backup/databases/proto"
	"service/backup/databases/server/config"
	"service/backup/databases/server/controllers"
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

type FileUpoad struct {
	proto.FileUploadServer
}

func (s *FileUpoad) UploadFile(stream proto.FileUpload_UploadFileServer) error {
	// Temporary buffer to hold file chunks
	fileName := ""
	var file []byte

	// Loop to receive file chunks
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			// End of file, save the file here
			log.Println("File received, saving...")
			// You can save 'file' to disk or wherever you want
			return stream.SendAndClose(&proto.UploadStatus{Success: true})
		}
		if err != nil {
			return err
		}
		if fileName == "" {
			fileName = chunk.DbName
		}

		file = append(file, chunk.ZipFile...)
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

	_ = os.Mkdir("upload", 0777)
}

func runGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterFileUploadServer(s, &FileUpoad{})
	log.Println("Server started at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	// Buat channel untuk menangkap sinyal SIGINT (Ctrl+C) dan SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Jalankan gRPC server secara asinkron
	go runGRPCServer()

	// Jalankan Fiber server secara asinkron
	go runFiberServer()

	// Tunggu sinyal SIGINT atau SIGTERM
	<-stop
	log.Println("Shutting down servers...")
}
