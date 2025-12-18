package main

import (
	"fmt"
	"log"
	"time"
	
	"github.com/cumulus13/go-gntp"
)

func testMode(mode gntp.IconMode, modeName string, iconPath string) {
	fmt.Printf("\n=== Testing %s Mode ===\n\n", modeName)
	
	client := gntp.NewClient(fmt.Sprintf("Test %s", modeName)).
		WithIconMode(mode).
		WithDebug(false)
	
	icon, err := gntp.LoadResource(iconPath)
	if err != nil {
		log.Printf("Failed to load icon: %v", err)
		return
	}
	
	fmt.Printf("Icon loaded: %d bytes\n", len(icon.Data))
	
	// Show icon reference format using public method
	iconRef := icon.GetReference(mode)
	if len(iconRef) > 80 {
		fmt.Printf("Icon ref: %s...\n", iconRef[:80])
	} else {
		fmt.Printf("Icon ref: %s\n", iconRef)
	}
	
	notification := gntp.NewNotificationType("test").
		WithDisplayName(fmt.Sprintf("Test %s", modeName)).
		WithIcon(icon)
	
	fmt.Println("\nRegistering...")
	if err := client.Register([]*gntp.NotificationType{notification}); err != nil {
		log.Printf("Registration failed: %v", err)
		return
	}
	fmt.Println("✓ Registered")
	
	fmt.Println("Sending notification...")
	if err := client.Notify("test", fmt.Sprintf("%s Mode Test", modeName), "Does icon show?"); err != nil {
		log.Printf("Notification failed: %v", err)
		return
	}
	fmt.Println("✓ Sent")
	
	fmt.Printf("\n✅ Check if icon appears in %s mode!\n", modeName)
	fmt.Println("Waiting 3 seconds before next test...")
	time.Sleep(3 * time.Second)
}

func main() {
	iconPath := "icon.png"
	
	fmt.Println("=== Icon Mode Testing ===")
	fmt.Printf("Using icon: %s\n", iconPath)
	
	// Test all modes
	testMode(gntp.IconModeBinary, "Binary", iconPath)
	testMode(gntp.IconModeDataURL, "DataURL", iconPath)
	testMode(gntp.IconModeFileURL, "FileURL", iconPath)
	
	fmt.Println("\n=== All tests completed! ===")
	fmt.Println("\nResults:")
	fmt.Println("  Which mode showed the icon?")
	fmt.Println("    • Binary   - GNTP spec (binary data)")
	fmt.Println("    • DataURL  - Base64 embedded (no line breaks)")
	fmt.Println("    • FileURL  - Absolute file path (file:///...)")
	fmt.Println("\nRecommendations:")
	fmt.Println("  • Windows Growl → Binary mode (works best)")
	fmt.Println("  • Android Growl → DataURL mode")
	fmt.Println("  • macOS Growl   → Binary mode")
}