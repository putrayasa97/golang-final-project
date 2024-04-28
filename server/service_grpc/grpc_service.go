package servicegrpc

import (
	"io"
	"os"
	"service/backup/databases/proto"
	"service/backup/databases/server/model"
	"service/backup/databases/server/utils"
)

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
