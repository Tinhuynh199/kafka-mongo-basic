package main

import (
	"fmt"
	"kafkamongo/internal/app"
	"strconv"
	"context"
)

func main() {
	
	// Loading Config
	config := app.GetConfig("./configs/config.yml")

	// Loading App
	app := &app.App{}
	app.Initialize(context.Background())

	// Start Server
	server := ""
	if config.Server.Port != nil {
		server = ":" + strconv.FormatInt(*config.Server.Port, 10)
	}
	fmt.Println("Start server")
	app.Run(server)
}

