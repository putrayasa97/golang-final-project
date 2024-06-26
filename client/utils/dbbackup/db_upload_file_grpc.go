package dbbackup

import (
	"context"
	"fmt"
	"io"
	"os"
	"service/backup/databases/client/model"
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
		mErr = fmt.Sprintf("Error uploadFileGrpc failed to dial server, Error : %v\n", err.Error())
		return mErr, err
	}
	defer conn.Close()

	client := proto.NewFileUploadClient(conn)

	file, err := os.Open(pathNameFileZip)
	if err != nil {
		mErr = fmt.Sprintf("Error uploadFileGrpc to open file with db %s to zip , Error : %s\n", fileName.NameDatabaseFile, err.Error())
		return mErr, err
	}
	defer file.Close()

	stream, err := client.UploadFile(context.Background())
	if err != nil {
		mErr = fmt.Sprintf("Error uploadFileGrpc to upload file with db %s, Error : %s\n", fileName.NameDatabaseFile, err.Error())
		return mErr, err
	}

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			mErr = fmt.Sprintf("Error uploadFileGrpc to read file with db %s, Error : %s\n", fileName.NameDatabaseFile, err.Error())
			return mErr, err
		}
		if n == 0 {
			break
		}

		if err := stream.Send(&proto.FileChunk{NameFile: fileName.NameFileZip, NameDb: fileName.NameDatabaseFile, ZipFile: buf[:n]}); err != nil {
			mErr = fmt.Sprintf("Error uploadFileGrpc to send chunk with db %s, Error : %s\n", fileName.NameDatabaseFile, err.Error())
			return mErr, err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		mErr = fmt.Sprintf("Error uploadFileGrpc to receive response with db %s, Error : %s\n", fileName.NameDatabaseFile, err.Error())
		return mErr, err
	}

	return mErr, nil
}
