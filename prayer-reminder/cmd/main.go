package main

import (
	"log"

	_ "crosswand.com/prayer-reminder"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	port := "8080"
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}