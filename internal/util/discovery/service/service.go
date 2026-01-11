package service

import (
	"os"
)

var Name = "main"

func init() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return
	}

	Name = port
}
