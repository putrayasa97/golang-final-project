package dbbackup

import (
	"context"
	"fmt"
	"io"
	"os"
	"service/backup/databases/client/model"
	"service/backup/databases/client/utils/logger"
	"service/backup/databases/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func uploadFileGrpc(pathFile *model.PathFile, fileName model.NameFile) (string, error) {
	pathNameFileZip := fmt.Sprintf("%s/%s", pathFile.PathFileZip, fileName.NameFileZip)
	mErr := ""

	conn, err := grpc.NewClient(":50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		mErr = fmt.Sprintf("Error failed to dial server , Error : %v\n", err.Error())
		return mErr, err
	}
	defer conn.Close()

	client := proto.NewFileUploadClient(conn)

	file, err := os.Open(pathNameFileZip)
	if err != nil {
		mErr = fmt.Sprintf("Error to open file %s to zip , Error : %s\n", fileName.NameFileZip, err.Error())
		return mErr, err
	}
	defer file.Close()

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		mErr = fmt.Sprintf("Error to upload file: %s  , Error : %s\n", fileName.NameFileZip, err.Error())
		return mErr, err
	}

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			mErr = fmt.Sprintf("Error to read file:  %s  , Error : %s\n", fileName.NameFileZip, err.Error())
			return mErr, err
		}
		if n == 0 {
			break
		}

		if err := stream.Send(&proto.FileChunk{NameFile: fileName.NameFileZip, ZipFile: buf[:n]}); err != nil {
			mErr = fmt.Sprintf("Error to send chunk  %s  , Error : %s\n", fileName.NameFileZip, err.Error())
			return mErr, err
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		mErr = fmt.Sprintf("Error to receive response:  %s  , Error : %s\n", fileName.NameFileZip, err.Error())
		return mErr, err
	}

	mErr = fmt.Sprintf("File uploaded successfully: %t", response.Success)
	logger.Info(mErr)

	return mErr, nil
}
