package main

import "github.com/franklindh/simedis-api/internal/router"

func main() {
	r := router.New()

	r.Run(":3000")
}
