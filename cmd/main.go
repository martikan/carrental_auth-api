package main

import (
	"github.com/martikan/carrental_auth-api/api"
)

func main() {
	server := api.InitApi()
	server.Start()
}
