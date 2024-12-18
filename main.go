package main

import (
	"context"
	"fmt"
	"github.com/opentreehole/backend/model"
	"io"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	model.Init()
	var images []model.OriginalImageTable
	result := model.OriginalDB.Find(&images)
	slog.Info("find all original images successfully")

	if result.Error != nil {
		log.Fatal(result.Error)
	}

	for _, image := range images {
		imageURL := "https://pic.jingyijun.xyz:8443/i/" + image.Path + "/" + image.Name
		fmt.Println("Downloading image from:", imageURL)

		imageData, err := downloadImage(imageURL)
		if err != nil {
			log.Printf("Error downloading image: %v\n", err)
			continue
		}
		fmt.Println("Image downloaded successfully.")

		originalFileName := image.OriginName // 用户上传的文件名
		imageFullName := image.Name          // 66f2cbaf9c143.png
		imageIdentifier := strings.TrimSuffix(imageFullName, filepath.Ext(imageFullName))
		fileExtension := strings.TrimPrefix(filepath.Ext(imageFullName), ".")

		createdAt := image.CreatedAt
		updatedAt := image.UpdatedAt

		// -------------------------------------------------------

		err = storeImageInDatabase(originalFileName, fileExtension, imageData, imageIdentifier, createdAt, updatedAt)
		if err != nil {
			log.Printf("Error storing image in database: %v\n", err)
		} else {
			fmt.Println("Image stored in database successfully.")
		}

	}
}

func downloadImage(imageURL string) ([]byte, error) {

	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// 读取图片数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imageData, nil
}

func storeImageInDatabase(originalFileName, fileExtension string, imageData []byte, imageIdentifier string, createdAt time.Time, updatedAt time.Time) error {

	uploadedImage := &model.NewImageTable{
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		ImageIdentifier: imageIdentifier,
		// 待替换
		OriginalFileName: originalFileName,
		ImageType:        fileExtension,
		ImageFileData:    imageData,
	}

	err := model.NewDB.Create(&uploadedImage).Error
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Database cannot store the image",
			slog.String("err", err.Error()), slog.String("fileName", originalFileName))
		return err
	}

	return nil
}
