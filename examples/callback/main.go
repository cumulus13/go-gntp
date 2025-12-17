package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/cumulus13/go-gntp"
)

func main() {
	fmt.Println("=== GNTP Callback Example ===\n")
	
	// Create client
	client := gntp.NewClient("Callback Example").
		WithDebug(true)
	
	// Set up callback handler
	fmt.Println("Setting up callback handler...")
	err := client.WithCallback(func(info gntp.CallbackInfo) {
		fmt.Printf("\n🔔 CALLBACK RECEIVED!\n")
		fmt.Printf("   Type: %s\n", info.Type)
		fmt.Printf("   Notification ID: %s\n", info.NotificationID)
		fmt.Printf("   Context: %s\n", info.Context)
		fmt.Printf("   Timestamp: %s\n\n", info.Timestamp.Format(time.RFC3339))
		
		switch info.Type {
		case gntp.CallbackClick:
			fmt.Println("   → User CLICKED the notification!")
		case gntp.CallbackClose:
			fmt.Println("   → User CLOSED the notification")
		case gntp.CallbackTimeout:
			fmt.Println("   → Notification TIMED OUT")
		}
	})
	
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Callback handler ready\n")
	
	defer client.Close()
	
	// Load icon
	icon, err := gntp.LoadResource("icon.png")
	if err != nil {
		fmt.Printf("⚠ Icon not found (optional): %v\n\n", err)
		icon = nil
	} else {
		fmt.Println("✓ Icon loaded\n")
	}
	
	// Define notification type
	notification := gntp.NewNotificationType("alert").
		WithDisplayName("Alert Notification").
		WithIcon(icon)
	
	// Register
	fmt.Println("Registering with Growl...")
	if err := client.Register([]*gntp.NotificationType{notification}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Registered\n")
	
	// Send notification with callback
	fmt.Println("Sending notification with callback...")
	options := gntp.NewNotifyOptions().
		WithSticky(true).
		WithCallbackContext("user_data_123").
		WithCallbackTarget("https://example.com")
	
	if err := client.NotifyWithOptions(
		"alert",
		"Click Me!",
		"This notification has a callback. Click it to see the callback in action!",
		options,
	); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Sent\n")
	
	fmt.Println("🎯 Waiting for callbacks...")
	fmt.Println("   Click the notification to trigger callback!")
	fmt.Println("   Press Ctrl+C to exit\n")
	
	// Wait indefinitely for callbacks
	select {}
}