package main

import (
	"fmt"
	"log"
	
	"github.com/cumulus13/go-gntp"
)

func main() {
	fmt.Println("=== Basic GNTP Notification ===\n")
	
	// Create client
	client := gntp.NewClient("Go GNTP Example")
	
	// Define notification type
	notification := gntp.NewNotificationType("alert").
		WithDisplayName("Alert Notification")
	
	// Register
	fmt.Println("Registering with Growl...")
	if err := client.Register([]*gntp.NotificationType{notification}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Registered\n")
	
	// Send notification
	fmt.Println("Sending notification...")
	if err := client.Notify("alert", "Hello from Go!", "This is a test notification"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Sent\n")
	
	fmt.Println("✅ Done!")
}