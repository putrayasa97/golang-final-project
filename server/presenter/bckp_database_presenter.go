package presenter

import (
	"fmt"
	"service/backup/databases/server/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BckpDatabase struct {
	ID           uint   `json:"id"`
	DatabaseName string `json:"database_name"`
	FileName     string `json:"file_name"`
	TimeStamp    string `json:"time_stamp"`
}

func AddBckpDatabseSuccessResponse(data *model.BckpDatabase) *fiber.Map {
	loc, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		fmt.Println("Error loading location:", err)
	}

	bckpDatabase := BckpDatabase{
		ID:           data.ID,
		DatabaseName: data.DatabaseName,
		FileName:     data.FileName,
		TimeStamp:    data.UpdatedAt.In(loc).Format("2006-01-02 15:04:05 MST"),
	}
	return &fiber.Map{
		"data":    bckpDatabase,
		"message": "Success Import",
	}
}
