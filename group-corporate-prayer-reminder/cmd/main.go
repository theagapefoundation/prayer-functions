package main

import (
	"log"

	_ "theagapefoundation.org/group-corporate-prayer-reminder"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	port := "8080"
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}