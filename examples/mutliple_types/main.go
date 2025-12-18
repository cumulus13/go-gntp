package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/cumulus13/go-gntp"
)

func main() {
	fmt.Println("=== Multiple Notification Types Example ===\n")
	
	// Create GNTP client
	client := gntp.NewClient("Multi-Type App").
		WithIconMode(gntp.IconModeDataURL)
	
	// Define multiple notification types
	info := gntp.NewNotificationType("info").
		WithDisplayName("Information").
		WithEnabled(true)
	
	warning := gntp.NewNotificationType("warning").
		WithDisplayName("Warning").
		WithEnabled(true)
	
	errorType := gntp.NewNotificationType("error").
		WithDisplayName("Error").
		WithEnabled(true)
	
	// Register all types at once
	fmt.Println("Registering notification types...")
	err := client.Register([]*gntp.NotificationType{info, warning, errorType})
	if err != nil {
		log.Fatalf("Registration failed: %v", err)
	}
	fmt.Println("✓ Registered 3 notification types\n")
	
	// Send different types of notifications
	fmt.Println("Sending info notification...")
	if err := client.Notify(
		"info",
		"System Information",
		"Application started successfully",
	); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Info sent\n")
	
	time.Sleep(2 * time.Second)
	
	fmt.Println("Sending warning notification...")
	if err := client.Notify(
		"warning",
		"Low Disk Space",
		"Only 10% disk space remaining",
	); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Warning sent\n")
	
	time.Sleep(2 * time.Second)
	
	fmt.Println("Sending error notification...")
	if err := client.Notify(
		"error",
		"Critical Error",
		"Database connection failed",
	); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Error sent\n")
	
	fmt.Println("✅ Example completed!")
	fmt.Println("\nYou should see 3 different notifications on your screen.")
}