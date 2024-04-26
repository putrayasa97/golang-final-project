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
	"service/backup/databases/server/model"
	"service/backup/databases/server/utils"
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

type FileUpload struct {
	proto.FileUploadServer
}

func (s *FileUpload) UploadFile(stream proto.FileUpload_UploadFileServer) error {
	firstChunk, err := stream.Recv()
	if err != nil {
		return err
	}
	nameFile := firstChunk.GetNameFile()
	nameBD := firstChunk.GetNameDb()

	fileName := "upload/" + nameFile

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(firstChunk.GetZipFile()); err != nil {
		return err
	}

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if _, err := file.Write(chunk.GetZipFile()); err != nil {
			return err
		}
	}

	_, errCreateBckpDatabse := utils.InsertBckpDatabase(model.BckpDatabase{
		DatabaseName: nameBD,
		FileName:     nameFile,
		FilePath:     fileName,
	})

	if errCreateBckpDatabse != nil {
		return errCreateBckpDatabse
	}

	if err := stream.SendAndClose(&proto.UploadStatus{Success: true}); err != nil {
		return err
	}

	return nil
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
	proto.RegisterFileUploadServer(s, &FileUpload{})
	log.Println("Server started at :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go runGRPCServer()

	go runFiberServer()

	<-stop
	log.Println("Shutting down servers...")
}
