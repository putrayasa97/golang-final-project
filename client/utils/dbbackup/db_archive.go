package dbbackup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"service/backup/databases/client/model"
)

func archiveDatabase(pathFile model.PathFile, fileName model.NameFile) (model.NameFile, string, error) {
	pathNameFileSql := fmt.Sprintf("%s/%s", pathFile.PathFileSql, fileName.NameFileSql)
	nameFileZip := fmt.Sprintf("%s.zip", fileName.NameFileSql)
	pathNameFileZip := fmt.Sprintf("%s/%s", pathFile.PathFileZip, nameFileZip)
	mErr := ""

	archive, err := os.Create(pathNameFileZip)
	if err != nil {
		mErr = fmt.Sprintf("Error archiveDatabase creating zip file with db %s, Error : %s\n", fileName.NameDatabaseFile, err.Error())
		return model.NameFile{
			NameFileSql: fileName.NameFileSql,
			NameFileZip: nameFileZip,
		}, mErr, err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	f, err := os.Open(pathNameFileSql)
	if err != nil {
		mErr = fmt.Sprintf("Error archiveDatabase opening sql file with db %s, Error : %s\n", fileName.NameDatabaseFile, err.Error())
		return model.NameFile{
			NameFileSql: fileName.NameFileSql,
			NameFileZip: nameFileZip,
		}, mErr, err
	}
	defer f.Close()

	w, err := zipWriter.Create(fileName.NameFileSql)
	if err != nil {
		mErr = fmt.Sprintf("Error archiveDatabase creating file with db %s in zip, Error : %s\n", fileName.NameDatabaseFile, err.Error())
		return model.NameFile{
			NameFileSql: fileName.NameFileSql,
			NameFileZip: nameFileZip,
		}, mErr, err
	}

	if _, err := io.Copy(w, f); err != nil {
		mErr = fmt.Sprintf("Error archiveDatabase copying file with db %s to zip , Error : %s\n", fileName.NameDatabaseFile, err.Error())
		return model.NameFile{
			NameFileSql: fileName.NameFileSql,
			NameFileZip: nameFileZip,
		}, mErr, err
	}
	zipWriter.Close()

	return model.NameFile{
		NameFileSql:      fileName.NameFileSql,
		NameFileZip:      nameFileZip,
		NameDatabaseFile: fileName.NameDatabaseFile,
	}, mErr, nil
}
