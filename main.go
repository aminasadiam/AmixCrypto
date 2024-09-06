package main

import (
	"AmixCrypto/cmd"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	Token string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
	Token = os.Getenv("TOKEN")
}

func main() {
	go func(token string) {
		err := cmd.Execute(token)
		if err != nil {
			logrus.Fatalln(err)
			return
		}
	}(Token)
	fmt.Println("Telegram Bot Started...")

	mux, err := cmd.Serve()
	if err != nil {
		logrus.Fatalln(err)
		return
	}
	fmt.Println("Server Started at 8080...")
	http.ListenAndServe(":8081", mux)
}
