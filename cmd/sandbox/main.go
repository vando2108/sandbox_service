package main

import (
	"context"
	"log"

	"github.com/vando2108/sandbox_service/internal/app/sandbox"
)

func main() {
	sandboxServer := sandbox.NewServer()
	if err := sandboxServer.Start(context.TODO()); err != nil {
		log.Fatalln("failed to start the server:", err)
	}
}
