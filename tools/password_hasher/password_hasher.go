package main

import (
	"fmt"
	"os"
	"wallet-app-server/app/util"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("This is a small tool to help you to generate your password hash for wallet-app system")
		fmt.Println("Usage: password_hasher <input_password>")
		os.Exit(-1)
	}
	fmt.Println(util.HashPassword(os.Args[1]))
}
