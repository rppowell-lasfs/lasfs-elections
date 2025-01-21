package main

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	sum := sha256.Sum256([]byte("hello world\n"))
	fmt.Printf("%x", sum)
}
