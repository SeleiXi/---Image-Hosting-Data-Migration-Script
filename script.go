package main

import (
	"context"
	"fmt"
	"github.com/opentreehole/backend/model"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	model.Init()

	for {
		// 1.存createdAt和uploadedAt 2.存Original file name

		imageURL := "http://localhost:8000/api/i/2024/12/17/6298d7df669f33.png"

		imageData, err := downloadImage(imageURL)
		if err != nil {
			log.Printf("Error downloading image: %v\n", err)
			continue
		}
		fmt.Println("Image downloaded successfully.")

		fileName := filepath.Base(imageURL) // 62986b5df57607d82a1c7ae10.jpg
		imageIdentifier := strings.TrimSuffix(fileName, filepath.Ext(fileName))
		fileExtension := strings.TrimPrefix(filepath.Ext(fileName), ".")

		// --------------------------------
		// 解析 URL
		parsedURL, err := url.Parse(imageURL)
		if err != nil {
			fmt.Printf("Error parsing URL: %v\n", err)
			return
		}

		urlPath := path.Clean(parsedURL.Path) // 提取 "/api/i/2024/12/17/62986b5df57607d82a1c7ae10.jpg"
		parts := strings.Split(urlPath, "/")  // 分割路径

		// 确保路径包含足够的段以提取日期
		if len(parts) < 5 {
			fmt.Println("URL does not contain a valid date path")
			return
		}

		// 提取并组合日期字符串
		dateString := fmt.Sprintf("%s/%s/%s", parts[len(parts)-4], parts[len(parts)-3], parts[len(parts)-2]) // "2024/12/17"

		// 解析日期
		createdAt, err := time.Parse("2006/01/02", dateString)
		if err != nil {
			fmt.Printf("Error parsing date: %v\n", err)
			return
		}

		// 待替换
		updatedAt := createdAt

		// -------------------------------------------------------

		if err != nil {
			fmt.Printf("Error extracting date: %v\n", err)
			return
		}

		err = storeImageInDatabase(fileName, fileExtension, imageData, imageIdentifier, createdAt, updatedAt)
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

func storeImageInDatabase(fileName, fileExtension string, imageData []byte, imageIdentifier string, createdAt time.Time, updatedAt time.Time) error {

	uploadedImage := &model.ImageTable{
		CreatedAt:       createdAt,
		UpdatedAt:       updatedAt,
		ImageIdentifier: imageIdentifier,
		// 待替换（从其他地方获取文件的Original Name）
		OriginalFileName: fileName,
		ImageType:        fileExtension,
		ImageFileData:    imageData,
	}

	err := model.DB.Create(&uploadedImage).Error
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Database cannot store the image",
			slog.String("err", err.Error()), slog.String("fileName", fileName))
		return err
	}

	return nil
}

// package main
//
// import (
// 	"fmt"
// 	"github.com/opentreehole/backend/import_image/model"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"path/filepath"
// 	"strings"
// 	"log/slog"
// )
//
// func main() {
// 	model.Init()
//
// 	for {
// 		imageURL := ""
// 		resp, err := http.Get(imageURL)
// 		if err != nil {
// 			log.Printf(err.Error())
// 		}
// 		if resp.StatusCode != http.StatusOK {
// 			log.Printf("bad status: %s", resp.Status)
// 		}
//
// 		if err != nil {
// 			fmt.Println("Error downloading image.")
// 		} else {
// 			fmt.Println("Image downloaded.")
// 		}
//
// 		imageData, err := ioutil.ReadAll(resp.Body)
// 		file := ""
// 		fileExtension := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
// 		fileContent, err := file.Open()
//
//
// 		imageUrl :=
// 		uploadedImage := &ImageTable{
// 			ImageIdentifier: imageIdentifier,
// 			BaseName:        originalFileName,
// 			ImageType:       fileExtension,
// 			ImageFileData:   imageData,
// 		}
// 		err = model.DB.Create(&uploadedImage).Error
//
// 		if err != nil {
// 			slog.LogAttrs(context.Background(), slog.LevelError, "Database cannot store the image", slog.String("err", err.Error()))
// 		}
//
//
// 	}
// }
