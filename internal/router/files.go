package router

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/db"
)

type ResultFile struct {
	Id  uint   `json:"id"`
	Url string `json:"url"`
}

func HandleListAllFiles(c *fiber.Ctx) error {
	var files []*db.File
	db.Db.Find(&files)

	var result []ResultFile
	for _, file := range files {
		fileUrl := fmt.Sprintf("%s/%s", os.Getenv("APP_URL"), file.Name)

		result = append(result, ResultFile{
			Id:  file.ID,
			Url: fileUrl,
		})
	}

	return c.JSON(result)
}

func HandleFindFileByName(c *fiber.Ctx) error {
	fileName := c.Params("name")

	var file *db.File
	db.Db.Find(&file, "name = ?", fileName)

	fileUrl := fmt.Sprintf("%s/%s", os.Getenv("APP_URL"), file.Name)

	result := &ResultFile{
		Id:  file.ID,
		Url: fileUrl,
	}

	return c.JSON(result)
}

func HandleUploadFile(c *fiber.Ctx) error {
	var err error

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	fileName := file.Filename

	err = c.SaveFile(file, fmt.Sprintf("./static/%s", fileName))
	if err != nil {
		return err
	}

	newFile := db.File{
		Name: fileName,
	}

	tx := db.Db.Create(&newFile)
	if tx.Error != nil {
		return tx.Error
	}

	return c.JSON(newFile)
}
