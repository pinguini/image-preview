package cache

import (
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const cacheDir = "./cache"

var (
	cache *lru.Cache
	mu    sync.Mutex
)

func init() {
	// Инициализация LRU-кэша с максимальным размером 100 элементов
	cache, _ = lru.New(100)
}

func SaveToCache(imgPath, hashString string) error {
	// Сохранение изображения в кэш
	mu.Lock()
	defer mu.Unlock()

	cache.Add(hashString, imgPath)
	return nil
}

func ImageExistsInCache(hashString string) (string, bool) {
	// Проверка наличия изображения в кэше
	mu.Lock()
	defer mu.Unlock()

	if imgPath, ok := cache.Get(hashString); ok {
		return imgPath.(string), true
	}
	return "", false
}

func EvictLRU() {
	// Функция для удаления наименее используемого элемента в LRU-кэше
	for {
		time.Sleep(time.Minute) // Периодическая проверка
		mu.Lock()
		// TODO: заимплементить конфиггурацию на основе переменных окружение с дефолтным значением 10
		if cache.Len() > 100 {
			cache.RemoveOldest()
		}
		mu.Unlock()
	}
}
