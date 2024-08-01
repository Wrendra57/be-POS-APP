package be

import (
	"fmt"
	be "github.com/Wrendra57/Pos-app-be"
	"github.com/Wrendra57/Pos-app-be/config"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	fmt.Println("server starting")
	config.InitConfig()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	app, cleanup, err := be.InitializeApp()

	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer cleanup()

	fmt.Println("Server Ready")
	log.Fatal(app.Listen(":8080"))
}
