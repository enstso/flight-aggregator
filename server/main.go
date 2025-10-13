/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"aggregator/internal/config"
	"aggregator/internal/health"
	"fmt"
	"net/http"
)

func main() {
	//load the config .env
	config.Load()

	http.HandleFunc("/health", health.HealthHandler)
	fmt.Println("Server running on :8080")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
