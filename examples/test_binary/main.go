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
	fmt.Println("=== Testing Binary Mode (Force) ===\n")
	
	// Load icon
	// icon, err := gntp.LoadResource("icon.png")
	icon1, err := getIconPath()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("✓ Icon loaded: %d bytes\n\n", len(icon1.Data))
	fmt.Printf("✓ Icon loaded: %d bytes\n\n", len(icon1))
	
	icon, err := gntp.LoadResource(icon1)

	// Binary mode with debug
	client := gntp.NewClient("Binary Test").
		WithIconMode(gntp.IconModeBinary).
		WithDebug(true)
	
	notification := gntp.NewNotificationType("alert").
		WithDisplayName("Alert").
		WithIcon(icon)
	
	fmt.Println("Registering with Binary mode...")
	if err := client.Register([]*gntp.NotificationType{notification}); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("\nSending notification...")
	if err := client.Notify("alert", "Binary Mode Test", "Does icon show?"); err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("\n✅ Done! Check if icon appears.")
}