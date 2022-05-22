package main

import (
	"log"
	"os"
	functions "soteria-functions"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	funcframework.RegisterHTTPFunction("/CreateBtcAccount", functions.CreateBtcAccount)
	funcframework.RegisterHTTPFunction("/CreateEthAccount", functions.CreateEthAccount)
	funcframework.RegisterHTTPFunction("/AddBtcBalance", functions.AddBtcBalance)
	funcframework.RegisterHTTPFunction("/AddUsdsBalance", functions.AddUsdsBalance)
	funcframework.RegisterHTTPFunction("/ConvertToken", functions.ConvertToken)

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
