package main

import (
	"fmt"
	"log"
	"os"
	"time"
	
	"github.com/cumulus13/go-gntp"
)

func main() {
	fmt.Println("=== Growl for Android Example ===\n")
	
	// Get Android device IP from environment
	androidHost := os.Getenv("ANDROID_HOST")
	if androidHost == "" {
		fmt.Println("⚠ ANDROID_HOST not set, using default")
		fmt.Println("  Set with: export ANDROID_HOST=192.168.1.100\n")
		androidHost = "192.168.1.100"
	}
	
	fmt.Printf("Target Android device: %s\n", androidHost)
	
	// Create client optimized for Android
	client := gntp.NewClient("Android Example").
		WithHost(androidHost).
		WithPort(23053).
		WithIconMode(gntp.IconModeDataURL).  // Best for Android
		WithTimeout(15 * time.Second).        // Longer timeout for mobile
		WithDebug(false)
	
	// Try to load icon (optional)
	var icon *gntp.Resource
	iconPath := "icon.png"
	
	if _, err := os.Stat(iconPath); err == nil {
		icon, err = gntp.LoadResource(iconPath)
		if err != nil {
			fmt.Printf("⚠ Failed to load icon: %v\n", err)
		} else {
			fmt.Printf("✓ Icon loaded: %s (%d bytes)\n", iconPath, len(icon.Data))
		}
	} else {
		fmt.Println("ℹ No icon found (optional)")
	}
	
	fmt.Println()
	
	// Define notification type with icon
	notification := gntp.NewNotificationType("android").
		WithDisplayName("Android Notification")
	
	if icon != nil {
		notification = notification.WithIcon(icon)
		fmt.Println("✓ Icon attached to notification")
	}
	
	// Register with retry (Android may need retry due to network)
	fmt.Println("Registering with Growl for Android...")
	
	var registerOK bool
	for attempt := 1; attempt <= 3; attempt++ {
		err := client.Register([]*gntp.NotificationType{notification})
		if err == nil {
			if attempt > 1 {
				fmt.Printf("✓ Registered successfully (attempt %d)\n\n", attempt)
			} else {
				fmt.Println("✓ Registered successfully\n")
			}
			registerOK = true
			break
		}
		
		if attempt < 3 {
			fmt.Printf("⚠ Attempt %d failed, retrying... (%v)\n", attempt, err)
			time.Sleep(2 * time.Second)
		} else {
			fmt.Printf("❌ Registration failed after 3 attempts: %v\n", err)
			fmt.Println("\nTroubleshooting:")
			fmt.Println("  1. Is Growl for Android running?")
			fmt.Printf("  2. Is %s the correct IP address?\n", androidHost)
			fmt.Println("  3. Are both devices on the same network?")
			fmt.Println("  4. Check Android firewall settings")
			log.Fatal(err)
		}
	}
	
	if !registerOK {
		log.Fatal("Registration failed")
	}
	
	// Send notification with options
	fmt.Println("Sending notification...")
	options := gntp.NewNotifyOptions().
		WithSticky(false).  // Don't make it sticky on mobile
		WithPriority(1)     // High priority
	
	err := client.NotifyWithOptions(
		"android",
		"Hello Android!",
		"This notification was sent from Go",
		options,
	)
	
	if err != nil {
		fmt.Printf("❌ Failed to send: %v\n", err)
		log.Fatal(err)
	}
	
	fmt.Println("✓ Notification sent\n")
	fmt.Println("✅ Check your Android device for the notification!")
}