package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	
	"github.com/cumulus13/go-gntp"
)

func getIconPath() (string, error) {
	// Try to find icon in various locations
	possiblePaths := []string{
		"icon.png",
		"growl.png",
		filepath.Join("examples", "with_icon", "icon.png"),
		filepath.Join("..", "..", "icon.png"),
		filepath.Join("..", "..", "growl.png"),
	}
	
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	
	return "", fmt.Errorf("icon not found in: %v", possiblePaths)
}

func main() {
	fmt.Println("=== GNTP With Icon Example ===\n")
	
	// Create client with DataURL mode (best for icon display)
	client := gntp.NewClient("Icon Example App").
		WithIconMode(gntp.IconModeDataURL)
	
	// Try to load icon
	fmt.Println("Loading icon...")
	iconPath, err := getIconPath()
	if err != nil {
		fmt.Printf("⚠ %v\n", err)
		fmt.Println("ℹ Continuing without icon...\n")
	} else {
		fmt.Printf("Found icon at: %s\n", iconPath)
	}
	
	var icon *gntp.Resource
	if iconPath != "" {
		icon, err = gntp.LoadResource(iconPath)
		if err != nil {
			fmt.Printf("✗ Failed to load icon: %v\n\n", err)
			icon = nil
		} else {
			fmt.Printf("✓ Icon loaded successfully (%d bytes)\n\n", len(icon.Data))
		}
	}
	
	// Define notification type
	// IMPORTANT: Only attach icon HERE (not to client or options)
	notification := gntp.NewNotificationType("alert").
		WithDisplayName("Alert Notification").
		WithEnabled(true)
	
	if icon != nil {
		notification = notification.WithIcon(icon)
		fmt.Println("✓ Icon attached to notification type")
	}
	
	// Register
	fmt.Println("Registering with Growl...")
	if err := client.Register([]*gntp.NotificationType{notification}); err != nil {
		log.Fatalf("Registration failed: %v", err)
	}
	fmt.Println("✓ Registered successfully\n")
	
	// Send notification WITHOUT icon in options
	// (icon already in notification type)
	fmt.Println("Sending notification...")
	if err := client.Notify(
		"alert",
		"Alert with Icon!",
		"This notification includes an icon from the notification type",
	); err != nil {
		log.Fatalf("Notification failed: %v", err)
	}
	fmt.Println("✓ Notification sent\n")
	
	fmt.Println("✅ Example completed!")
	fmt.Println("\nNote: To test with a custom icon:")
	fmt.Println("  1. Place 'icon.png' or 'growl.png' in the project root")
	fmt.Println("  2. Run: go run examples/with_icon/main.go")
	fmt.Println("\nIcon delivery mode: DataURL (base64 embedded)")
	fmt.Println("This works best across all platforms including Android.")
}