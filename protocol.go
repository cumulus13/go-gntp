package gntp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"crypto/md5"
	"time"
	"io"
)

// Register registers the application and notification types with Growl
func (c *Client) Register(notifications []*NotificationType) error {
	var packet strings.Builder
	resources := make([]*Resource, 0)
	seenIDs := make(map[string]bool)
	
	// Build REGISTER packet
	packet.WriteString(fmt.Sprintf("GNTP/%s REGISTER NONE%s", GNTPVersion, CRLF))
	packet.WriteString(fmt.Sprintf("Application-Name: %s%s", c.ApplicationName, CRLF))
	
	// Application icon
	if c.ApplicationIcon != nil {
		iconRef := c.ApplicationIcon.getReference(c.IconMode)
		packet.WriteString(fmt.Sprintf("Application-Icon: %s%s", iconRef, CRLF))
		
		if c.IconMode == IconModeBinary && !seenIDs[c.ApplicationIcon.Identifier] {
			resources = append(resources, c.ApplicationIcon)
			seenIDs[c.ApplicationIcon.Identifier] = true
		}
	}
	
	// Callback URL if handler is set
	if c.callbackURL != "" {
		packet.WriteString(fmt.Sprintf("Notification-Callback-Target: %s%s", c.callbackURL, CRLF))
	}
	
	packet.WriteString(fmt.Sprintf("Notifications-Count: %d%s", len(notifications), CRLF))
	packet.WriteString(CRLF)
	
	// Each notification type
	for _, notif := range notifications {
		packet.WriteString(fmt.Sprintf("Notification-Name: %s%s", notif.Name, CRLF))
		
		if notif.DisplayName != "" {
			packet.WriteString(fmt.Sprintf("Notification-Display-Name: %s%s", notif.DisplayName, CRLF))
		}
		
		enabled := "False"
		if notif.Enabled {
			enabled = "True"
		}
		packet.WriteString(fmt.Sprintf("Notification-Enabled: %s%s", enabled, CRLF))
		
		if notif.Icon != nil {
			iconRef := notif.Icon.getReference(c.IconMode)
			packet.WriteString(fmt.Sprintf("Notification-Icon: %s%s", iconRef, CRLF))
			
			if c.IconMode == IconModeBinary && !seenIDs[notif.Icon.Identifier] {
				resources = append(resources, notif.Icon)
				seenIDs[notif.Icon.Identifier] = true
			}
		}
		
		packet.WriteString(CRLF)
	}
	
	// Binary resources
	if c.IconMode == IconModeBinary {
		for _, res := range resources {
			packet.WriteString(fmt.Sprintf("Identifier: %s%s", res.Identifier, CRLF))
			packet.WriteString(fmt.Sprintf("Length: %d%s", len(res.Data), CRLF))
			packet.WriteString(CRLF)
		}
	}
	
	if c.Debug {
		fmt.Printf("\n=== REGISTER PACKET (Mode: %d) ===\n", c.IconMode)
		fmt.Println(packet.String())
		fmt.Printf("Resources: %d\n", len(resources))
		fmt.Println("======================================\n")
	}
	
	// Send packet
	var err error
	if c.IconMode == IconModeBinary {
		_, err = c.sendPacketWithResources(packet.String(), resources)
	} else {
		_, err = c.sendPacket(packet.String())
	}
	
	if err != nil {
		return err
	}
	
	c.registered = true
	return nil
}

// Notify sends a notification
func (c *Client) Notify(notificationName, title, text string) error {
	return c.NotifyWithOptions(notificationName, title, text, NewNotifyOptions())
}

// NotifyWithOptions sends a notification with options
func (c *Client) NotifyWithOptions(notificationName, title, text string, options *NotifyOptions) error {
	if !c.registered {
		return fmt.Errorf("must call Register() before Notify()")
	}
	
	var packet strings.Builder
	resources := make([]*Resource, 0)
	
	// Generate notification ID for callbacks
	notificationID := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s:%s:%d", c.ApplicationName, notificationName, time.Now().UnixNano()))))
	
	packet.WriteString(fmt.Sprintf("GNTP/%s NOTIFY NONE%s", GNTPVersion, CRLF))
	packet.WriteString(fmt.Sprintf("Application-Name: %s%s", c.ApplicationName, CRLF))
	packet.WriteString(fmt.Sprintf("Notification-Name: %s%s", notificationName, CRLF))
	packet.WriteString(fmt.Sprintf("Notification-ID: %s%s", notificationID, CRLF))
	packet.WriteString(fmt.Sprintf("Notification-Title: %s%s", title, CRLF))
	packet.WriteString(fmt.Sprintf("Notification-Text: %s%s", text, CRLF))
	
	if options.Sticky {
		packet.WriteString(fmt.Sprintf("Notification-Sticky: True%s", CRLF))
	}
	
	if options.Priority != 0 {
		packet.WriteString(fmt.Sprintf("Notification-Priority: %d%s", options.Priority, CRLF))
	}
	
	if options.Icon != nil {
		iconRef := options.Icon.getReference(c.IconMode)
		packet.WriteString(fmt.Sprintf("Notification-Icon: %s%s", iconRef, CRLF))
		
		if c.IconMode == IconModeBinary {
			resources = append(resources, options.Icon)
		}
	}
	
	// Callback settings
	if c.callbackURL != "" {
		packet.WriteString(fmt.Sprintf("Notification-Callback-Context: %s%s", options.CallbackContext, CRLF))
		packet.WriteString(fmt.Sprintf("Notification-Callback-Context-Type: string%s", CRLF))
		
		if options.CallbackTarget != "" {
			packet.WriteString(fmt.Sprintf("Notification-Callback-Target: %s%s", options.CallbackTarget, CRLF))
		}
	}
	
	packet.WriteString(CRLF)
	
	// Binary resources
	if c.IconMode == IconModeBinary {
		for _, res := range resources {
			packet.WriteString(fmt.Sprintf("Identifier: %s%s", res.Identifier, CRLF))
			packet.WriteString(fmt.Sprintf("Length: %d%s", len(res.Data), CRLF))
			packet.WriteString(CRLF)
		}
	}
	
	if c.Debug {
		fmt.Printf("\n=== NOTIFY PACKET (Mode: %d) ===\n", c.IconMode)
		fmt.Println(packet.String())
		fmt.Println("====================================\n")
	}
	
	// Send packet
	if c.IconMode == IconModeBinary {
		_, err := c.sendPacketWithResources(packet.String(), resources)
		return err
	}
	
	_, err := c.sendPacket(packet.String())
	return err
}

// SendMessage sends a notification using Message struct (compatibility method)
func (c *Client) SendMessage(msg *Message) error {
	// Load icon if specified
	var icon *Resource
	if msg.Icon != "" {
		var err error
		icon, err = LoadResource(msg.Icon)
		if err != nil {
			return fmt.Errorf("failed to load icon: %w", err)
		}
	}
	
	// Create notification type if not registered
	if !c.registered {
		displayName := msg.DisplayName
		if displayName == "" {
			displayName = msg.Event
		}
		
		notif := NewNotificationType(msg.Event).
			WithDisplayName(displayName).
			WithIcon(icon)
		
		if err := c.Register([]*NotificationType{notif}); err != nil {
			return err
		}
	}
	
	// Send notification
	options := NewNotifyOptions().
		WithSticky(msg.Sticky).
		WithPriority(msg.Priority)
	
	if icon != nil {
		options.WithIcon(icon)
	}
	
	if msg.Callback != "" {
		options.WithCallbackTarget(msg.Callback)
	}
	
	return c.NotifyWithOptions(msg.Event, msg.Title, msg.Text, options)
}

// sendPacket sends a text-only packet
func (c *Client) sendPacket(packet string) (string, error) {
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	
	if c.Debug {
		fmt.Printf("Connecting to %s...\n", address)
	}
	
	conn, err := net.DialTimeout("tcp", address, c.Timeout)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	defer conn.Close()
	
	// Set deadlines
	conn.SetDeadline(time.Now().Add(c.Timeout))
	
	// Send packet
	if _, err := conn.Write([]byte(packet)); err != nil {
		return "", fmt.Errorf("failed to send packet: %w", err)
	}
	
	if c.Debug {
		fmt.Println("Packet sent, waiting for response...")
	}
	
	// Read response
	reader := bufio.NewReader(conn)
	response := strings.Builder{}
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			// Connection closed is OK for some Growl versions
			if strings.Contains(err.Error(), "connection") {
				break
			}
			return "", fmt.Errorf("failed to read response: %w", err)
		}
		response.WriteString(line)
		
		// Check for end of response
		if strings.TrimSpace(line) == "" {
			break
		}
	}
	
	responseStr := response.String()
	
	if c.Debug {
		fmt.Printf("Response:\n%s\n", responseStr)
	}
	
	// Check for errors
	if strings.Contains(responseStr, "-ERROR") {
		return "", fmt.Errorf("server error: %s", responseStr)
	}
	
	return responseStr, nil
}

// sendPacketWithResources sends a packet with binary resources
func (c *Client) sendPacketWithResources(packet string, resources []*Resource) (string, error) {
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	
	conn, err := net.DialTimeout("tcp", address, c.Timeout)
	if err != nil {
		return "", fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	defer conn.Close()
	
	conn.SetDeadline(time.Now().Add(c.Timeout))
	
	// Send text packet
	if _, err := conn.Write([]byte(packet)); err != nil {
		return "", fmt.Errorf("failed to send packet: %w", err)
	}
	
	// Send binary resources
	for _, res := range resources {
		if _, err := conn.Write(res.Data); err != nil {
			return "", fmt.Errorf("failed to send resource data: %w", err)
		}
		if _, err := conn.Write([]byte(CRLF)); err != nil {
			return "", fmt.Errorf("failed to send resource CRLF: %w", err)
		}
	}
	
	// Message termination
	if _, err := conn.Write([]byte(CRLF)); err != nil {
		return "", fmt.Errorf("failed to send termination: %w", err)
	}
	
	// Read response
	reader := bufio.NewReader(conn)
	response := strings.Builder{}
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			if strings.Contains(err.Error(), "connection") {
				break
			}
			return "", fmt.Errorf("failed to read response: %w", err)
		}
		response.WriteString(line)
		
		if strings.TrimSpace(line) == "" {
			break
		}
	}
	
	responseStr := response.String()
	
	if strings.Contains(responseStr, "-ERROR") {
		return "", fmt.Errorf("server error: %s", responseStr)
	}
	
	return responseStr, nil
}