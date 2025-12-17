package main

import (
	"fmt"
	"log"
	
	"github.com/cumulus13/go-gntp"
)

func main() {
	fmt.Println("=== Message Struct Example ===\n")
	
	// Create client
	client := gntp.NewClient("Message Example")
	
	// Simple message struct (compatible API)
	msg := &gntp.Message{
		Event:       "alert",
		Title:       "Simple Message",
		Text:        "Using Message struct for easy notifications",
		Icon:        "icon.png",  // Optional
		Callback:    "https://github.com/cumulus13/go-gntp",  // URL to open on click
		DisplayName: "Alert",
		Sticky:      false,
		Priority:    1,
	}
	
	fmt.Println("Sending message...")
	if err := client.SendMessage(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✓ Sent\n")
	
	fmt.Println("✅ Done!")
	fmt.Println("\nThe Message struct is a simplified API compatible with gntplib.")
}