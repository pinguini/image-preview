// server/server.go
package server

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/pinguini/image-preview/cache"
	"github.com/pinguini/image-preview/image"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// const string as byte
	w.Write([]byte("this is default handler"))
}

func FillHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	re := regexp.MustCompile(`^/fill/(\d+)/(\d+)/(.+)$`)
	parts := re.FindStringSubmatch(r.URL.Path)

	if len(parts) != 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	newWidth, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "Invalid new width", http.StatusBadRequest)
		return
	}

	newHeight, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid new height", http.StatusBadRequest)
		return
	}

	imageURL := "http://" + parts[3]
	_, err = url.Parse(imageURL)
	if err != nil {
		http.Error(w, "Invalid URL: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Генерируем хэш SHA-256 от URL
	hash := sha256.Sum256([]byte(imageURL))
	hashString := string(hash[:])

	sourceImageReader := getFileFromCache(hashString)
	if sourceImageReader != nil {
		defer sourceImageReader.Close()
	} else {
		// Делаем запрос на сервер для получения изображения
		clientRequest, err := http.NewRequest("GET", imageURL, nil)
		if err != nil {
			http.Error(w, "Failed to create client request: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Добавляем заголовки из оригинального запроса клиента
		addRequestHeaders(clientRequest, r.Header)

		resp, err := http.DefaultClient.Do(clientRequest)
		if err != nil {
			http.Error(w, "Failed to fetch image: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Проверка размера
		maxSourceImageSize := 50 * 1024 * 1024 // 10MB
		contentSizeStr := resp.Header.Get("Content-Length")
		contentSize, err := strconv.Atoi(contentSizeStr)
		if err != nil {
			http.Error(w, "Failed to get content size: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if contentSize > maxSourceImageSize {
			http.Error(w, "Source image size exceeds the maximum allowed size", http.StatusBadRequest)
			return
		}

		// Определяем тип содержимого
		contentType := resp.Header.Get("Content-Type")
		allowedTypes := []string{"image/jpeg", "image/gif", "image/png"}
		imageType := ""
		for _, t := range allowedTypes {
			if contentType == t {
				imageType = t
			}
		}
		if imageType == "" {
			http.Error(w, "Failed to fetch image: Incorrect ContentType", http.StatusBadGateway)
			return
		}
		// Сохраняем изображение в кэш
		err = cache.SaveToCache(resp.Body, hashString)
		if err != nil {
			fmt.Printf("Failed to save image to cache: %s\n", err)
		}
		sourceImageReader = resp.Body
		putFileToCache(img, hashString)
	}

	// Обработка изображения с изменением размера
	img, err := image.ResizeImage(sourceImageReader, newWidth, newHeight)
	if err != nil {
		http.Error(w, "Failed to resize image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(img)
}

func getFileFromCache(hashString string) io.ReadCloser {
	imgPath, found := cache.ImageExistsInCache(hashString)
	if !found {
		return nil
	}
	fmt.Println("Get from cache: " + imgPath)
	file, err := os.Open(imgPath)
	if err != nil {
		return nil
	}

	return file
}

func putFileToCache(imgPath, hashString string) error {
	// Сохранение изображения в кэш
	err := cache.SaveToCache(imgPath, hashString)
	if err != nil {
		return err
}

func addRequestHeaders(req *http.Request, headers http.Header) {
	// Добавление заголовков из оригинального запроса клиента
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}
