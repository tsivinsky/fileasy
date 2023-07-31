package router

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/db"
)

type ResultFile struct {
	Id  uint   `json:"id"`
	Url string `json:"url"`
}

func HandleListAllFiles(c *fiber.Ctx) error {
	userId, err := GetUserIdFromRequest(c)
	if err != nil {
		return err
	}

	var files []*db.File
	if tx := db.Db.Preload("User").Find(&files, "user_id = ?", userId); tx.Error != nil {
		return tx.Error
	}

	return c.JSON(files)
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

	userId, err := GetUserIdFromRequest(c)
	if err != nil {
		return err
	}

	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	fileName := file.Filename
	fileExt := path.Ext(fileName)

	formFileName := c.FormValue("name")
	if formFileName != "" {
		fileName = fmt.Sprintf("%s%s", formFileName, fileExt)
	}

	err = c.SaveFile(file, fmt.Sprintf("./static/%s", fileName))
	if err != nil {
		return err
	}

	newFile := db.File{
		Name:   fileName,
		UserID: userId,
	}

	tx := db.Db.Create(&newFile)
	if tx.Error != nil {
		return tx.Error
	}

	// TODO: figure out how preloading works and preload user data in Create query above to avoid sending 2nd query

	var createdFile db.File
	if tx := db.Db.Preload("User").First(&createdFile, "id = ?", newFile.ID); tx.Error != nil {
		return tx.Error
	}

	return c.JSON(createdFile)
}

func HandleDeleteFile(c *fiber.Ctx) error {
	fileIdParam := c.Params("id")
	fileId, err := strconv.Atoi(fileIdParam)
	if err != nil {
		return errors.New("Invalid file id param")
	}

	userId, err := GetUserIdFromRequest(c)
	if err != nil {
		return errors.New("Not authorized")
	}

	var file db.File
	if tx := db.Db.First(&file, "id = ? AND user_id = ?", fileId, userId); tx.Error != nil {
		return tx.Error
	}

	if tx := db.Db.Delete(&file); tx.Error != nil {
		return tx.Error
	}

	fp := fmt.Sprintf("./static/%s", file.Name)
	err = os.Remove(fp)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"ok": true,
	})
}
