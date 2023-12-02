// server/server.go
package server

import (
	"fmt"
	"net/http"
	"project/cache"
	"project/image"
	"strconv"
	"strings"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 6 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	newWidth, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid new width", http.StatusBadRequest)
		return
	}

	newHeight, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid new height", http.StatusBadRequest)
		return
	}

	imageURL := parts[4]

	// Генерируем хэш SHA-256 от URL
	hash := sha256.Sum256([]byte(imageURL))
	hashString := fmt.Sprintf("%x", hash)

	// Проверяем, есть ли изображение в кэше
	if imgPath, found := cache.ImageExistsInCache(hashString); found {
		// Изображение найдено в кэше, отправляем его клиенту
		http.ServeFile(w, r, imgPath)
		return
	}

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
		http.Error(w, "Failed to fetch image: "+err.Error(), http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	// Определяем тип содержимого
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// Если заголовок Content-Type отсутствует, пытаемся определить тип изображения на основе данных
		contentType = image.DetectImageType(resp.Body)
		if contentType == "" {
			http.Error(w, "Unable to determine content type", http.StatusNotFound)
			return
		}
	}

	// Проверяем "magic number" перед изменением размера изображения
	if !image.CheckMagicNumberForImage(contentType, imageURL) {
		http.Error(w, "Invalid image type", http.StatusNotFound)
		return
	}

	// Обработка изображения с изменением размера
	img, err := image.ResizeImage(resp.Body, newWidth, newHeight)
	if err != nil {
		http.Error(w, "Failed to resize image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохранение изображения в кэш
	err = cache.SaveToCache(img, hashString)
	if err != nil {
		fmt.Printf("Failed to save image to cache: %s\n", err)
	}

	// Отправляем измененное изображение клиенту
	// ...
}