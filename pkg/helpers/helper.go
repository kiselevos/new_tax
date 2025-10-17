package helpers

import "os"

// GetGRPCWebPort хэлпер для заполнения адреса по дефолту
func GetGRPCWebPort() string {
	if port := os.Getenv("BACKEND_PORT"); port != "" {
		return ":" + port
	}
	return ":8081"
}
