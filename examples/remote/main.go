package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	
	"github.com/cumulus13/go-gntp"
)

func main() {
	fmt.Println("=== Remote Growl Server Example ===\n")
	
	// Get remote host from environment or use default
	remoteHost := os.Getenv("GROWL_HOST")
	if remoteHost == "" {
		remoteHost = "192.168.1.100"
	}
	
	remotePortStr := os.Getenv("GROWL_PORT")
	remotePort := 23053
	if remotePortStr != "" {
		if p, err := strconv.Atoi(remotePortStr); err == nil {
			remotePort = p
		}
	}
	
	fmt.Printf("Target: %s:%d\n", remoteHost, remotePort)
	fmt.Println("(Set GROWL_HOST and GROWL_PORT environment variables to change)\n")
	
	// Create client for remote server
	client := gntp.NewClient("Remote Example").
		WithHost(remoteHost).
		WithPort(remotePort).
		WithIconMode(gntp.IconModeDataURL)  // Best for remote/Android
	
	// Define notification type
	notification := gntp.NewNotificationType("remote").
		WithDisplayName("Remote Notification")
	
	// Register
	fmt.Println("Registering with remote Growl...")
	if err := client.Register([]*gntp.NotificationType{notification}); err != nil {
		fmt.Printf("❌ Registration failed: %v\n\n", err)
		fmt.Println("Troubleshooting:")
		fmt.Printf("  1. Is Growl running on %s?\n", remoteHost)
		fmt.Printf("  2. Is port %d open in firewall?\n", remotePort)
		fmt.Printf("  3. Is remote host reachable? (ping %s)\n", remoteHost)
		log.Fatal(err)
	}
	fmt.Println("✓ Registered successfully\n")
	
	// Send notification
	fmt.Println("Sending notification...")
	if err := client.Notify(
		"remote",
		"Hello from Remote!",
		"This notification was sent over the network",
	); err != nil {
		log.Fatalf("Notification failed: %v", err)
	}
	fmt.Println("✓ Notification sent\n")
	
	fmt.Println("✅ Example completed!")
}