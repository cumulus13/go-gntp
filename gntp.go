// Package gntp provides a production-ready GNTP (Growl Notification Transport Protocol) client
// for Go with full Windows/Android compatibility and callback support.
//
// Features:
//   - Multiple icon delivery modes (Binary, FileURL, DataURL, HttpURL)
//   - Callback support (click, close, timeout)
//   - Windows Growl compatibility
//   - Android Growl compatibility
//   - Retry mechanism
//   - Resource deduplication
package gntp

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	// GNTPVersion is the GNTP protocol version
	GNTPVersion = "1.0"
	
	// DefaultPort is the standard GNTP port
	DefaultPort = 23053
	
	// CRLF is the line ending for GNTP protocol
	CRLF = "\r\n"
)

// IconMode specifies how icons should be delivered
type IconMode int

const (
	// IconModeBinary sends icons as binary resources (GNTP spec compliant)
	// Not recommended for Windows Growl due to bugs
	IconModeBinary IconMode = iota
	
	// IconModeFileURL references icons via file:// URLs
	// Requires icon files to exist on disk
	IconModeFileURL
	
	// IconModeDataURL embeds icons as base64 data URLs (RECOMMENDED)
	// Best compatibility across all platforms
	IconModeDataURL
	
	// IconModeHttpURL references icons via http:// or https:// URLs
	// Best for remote servers and Android
	IconModeHttpURL
	
	// IconModeAuto automatically selects best mode (defaults to DataURL)
	IconModeAuto
)

// CallbackType represents the type of callback event
type CallbackType string

const (
	// CallbackClick is triggered when notification is clicked
	CallbackClick CallbackType = "CLICK"
	
	// CallbackClose is triggered when notification is closed
	CallbackClose CallbackType = "CLOSE"
	
	// CallbackTimeout is triggered when notification times out
	CallbackTimeout CallbackType = "TIMEOUT"
)

// CallbackInfo contains information about a callback event
type CallbackInfo struct {
	Type              CallbackType
	NotificationID    string
	Context           string
	ContextType       string
	Timestamp         time.Time
}

// CallbackHandler is a function that handles callback events
type CallbackHandler func(info CallbackInfo)

// Resource represents an icon resource
type Resource struct {
	Identifier string
	Data       []byte
	SourcePath string
	MimeType   string
}

// NotificationType defines a type of notification
type NotificationType struct {
	Name        string
	DisplayName string
	Enabled     bool
	Icon        *Resource
}

// NotifyOptions contains options for sending notifications
type NotifyOptions struct {
	Sticky          bool
	Priority        int      // -2 to 2
	Icon            *Resource
	CallbackContext string   // Custom data passed to callback
	CallbackTarget  string   // URL to open on click
}

// Message is a simplified notification structure (for compatibility)
type Message struct {
	Event       string
	Title       string
	Text        string
	Icon        string
	Callback    string
	DisplayName string
	Sticky      bool
	Priority    int
}

// Client is the GNTP client
type Client struct {
	Host             string
	Port             int
	ApplicationName  string
	ApplicationIcon  *Resource
	IconMode         IconMode
	Debug            bool
	Timeout          time.Duration
	registered       bool
	callbackListener net.Listener
	callbackHandler  CallbackHandler
	callbackURL      string
}

// NewClient creates a new GNTP client
func NewClient(applicationName string) *Client {
	return &Client{
		Host:            "localhost",
		Port:            DefaultPort,
		ApplicationName: applicationName,
		IconMode:        IconModeDataURL, // Safe default
		Debug:           false,
		Timeout:         10 * time.Second,
		registered:      false,
	}
}

// WithHost sets the Growl server hostname
func (c *Client) WithHost(host string) *Client {
	c.Host = host
	return c
}

// WithPort sets the Growl server port
func (c *Client) WithPort(port int) *Client {
	c.Port = port
	return c
}

// WithIconMode sets the icon delivery mode
func (c *Client) WithIconMode(mode IconMode) *Client {
	c.IconMode = mode
	return c
}

// WithIcon sets the application icon
func (c *Client) WithIcon(icon *Resource) *Client {
	c.ApplicationIcon = icon
	return c
}

// WithDebug enables debug output
func (c *Client) WithDebug(debug bool) *Client {
	c.Debug = debug
	return c
}

// WithTimeout sets connection timeout
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.Timeout = timeout
	return c
}

// WithCallback sets up callback handler
func (c *Client) WithCallback(handler CallbackHandler) error {
	c.callbackHandler = handler
	
	// Start callback listener
	listener, err := net.Listen("tcp", ":0") // Random port
	if err != nil {
		return fmt.Errorf("failed to start callback listener: %w", err)
	}
	
	c.callbackListener = listener
	
	// Get callback URL
	addr := listener.Addr().(*net.TCPAddr)
	c.callbackURL = fmt.Sprintf("http://%s:%d", getLocalIP(), addr.Port)
	
	// Start accepting callbacks
	go c.acceptCallbacks()
	
	return nil
}

// acceptCallbacks handles incoming callback connections
func (c *Client) acceptCallbacks() {
	for {
		conn, err := c.callbackListener.Accept()
		if err != nil {
			return // Listener closed
		}
		go c.handleCallback(conn)
	}
}

// handleCallback processes a single callback
func (c *Client) handleCallback(conn net.Conn) {
	defer conn.Close()
	
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	
	response := string(buf[:n])
	lines := strings.Split(response, CRLF)
	
	info := CallbackInfo{
		Timestamp: time.Now(),
	}
	
	for _, line := range lines {
		if strings.HasPrefix(line, "Notification-Callback-Result: ") {
			info.Type = CallbackType(strings.TrimPrefix(line, "Notification-Callback-Result: "))
		} else if strings.HasPrefix(line, "Notification-ID: ") {
			info.NotificationID = strings.TrimPrefix(line, "Notification-ID: ")
		} else if strings.HasPrefix(line, "Notification-Callback-Context: ") {
			info.Context = strings.TrimPrefix(line, "Notification-Callback-Context: ")
		} else if strings.HasPrefix(line, "Notification-Callback-Context-Type: ") {
			info.ContextType = strings.TrimPrefix(line, "Notification-Callback-Context-Type: ")
		}
	}
	
	// Send OK response
	conn.Write([]byte(fmt.Sprintf("GNTP/%s -OK NONE%s%s", GNTPVersion, CRLF, CRLF)))
	
	// Call handler
	if c.callbackHandler != nil {
		c.callbackHandler(info)
	}
}

// getLocalIP gets the local IP address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	
	return "127.0.0.1"
}

// LoadResource loads an icon from a file
func LoadResource(path string) (*Resource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	mimeType := guessMimeType(path)
	
	return &Resource{
		Identifier: uuid.New().String(),
		Data:       data,
		SourcePath: path,
		MimeType:   mimeType,
	}, nil
}

// LoadResourceFromBytes creates a resource from byte data
func LoadResourceFromBytes(data []byte, mimeType string) *Resource {
	return &Resource{
		Identifier: uuid.New().String(),
		Data:       data,
		MimeType:   mimeType,
	}
}

// guessMimeType guesses MIME type from file extension
func guessMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".ico":
		return "image/x-icon"
	case ".svg":
		return "image/svg+xml"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// getReference returns the icon reference string based on mode
func (r *Resource) getReference(mode IconMode) string {
	switch mode {
	case IconModeBinary:
		return fmt.Sprintf("x-growl-resource://%s", r.Identifier)
	case IconModeFileURL:
		if r.SourcePath != "" {
			path := strings.ReplaceAll(r.SourcePath, "\\", "/")
			return fmt.Sprintf("file:///%s", path)
		}
		return r.toDataURL()
	case IconModeHttpURL:
		if r.SourcePath != "" {
			return r.SourcePath // Assume it's already a URL
		}
		return r.toDataURL()
	default: // IconModeDataURL, IconModeAuto
		return r.toDataURL()
	}
}

// toDataURL converts resource to base64 data URL
func (r *Resource) toDataURL() string {
	encoded := base64.StdEncoding.EncodeToString(r.Data)
	
	// Add line breaks every 76 characters (MIME standard)
	var formatted strings.Builder
	for i := 0; i < len(encoded); i += 76 {
		end := i + 76
		if end > len(encoded) {
			end = len(encoded)
		}
		formatted.WriteString(encoded[i:end])
		if end < len(encoded) {
			formatted.WriteString(CRLF)
		}
	}
	
	return fmt.Sprintf("data:%s;base64,%s", r.MimeType, formatted.String())
}

// NewNotificationType creates a new notification type
func NewNotificationType(name string) *NotificationType {
	return &NotificationType{
		Name:    name,
		Enabled: true,
	}
}

// WithDisplayName sets the display name
func (nt *NotificationType) WithDisplayName(displayName string) *NotificationType {
	nt.DisplayName = displayName
	return nt
}

// WithIcon sets the icon
func (nt *NotificationType) WithIcon(icon *Resource) *NotificationType {
	nt.Icon = icon
	return nt
}

// WithEnabled sets enabled state
func (nt *NotificationType) WithEnabled(enabled bool) *NotificationType {
	nt.Enabled = enabled
	return nt
}

// NewNotifyOptions creates new notify options
func NewNotifyOptions() *NotifyOptions {
	return &NotifyOptions{
		Priority: 0,
	}
}

// WithSticky sets sticky mode
func (no *NotifyOptions) WithSticky(sticky bool) *NotifyOptions {
	no.Sticky = sticky
	return no
}

// WithPriority sets priority (-2 to 2)
func (no *NotifyOptions) WithPriority(priority int) *NotifyOptions {
	if priority < -2 {
		priority = -2
	} else if priority > 2 {
		priority = 2
	}
	no.Priority = priority
	return no
}

// WithIcon sets the notification icon
func (no *NotifyOptions) WithIcon(icon *Resource) *NotifyOptions {
	no.Icon = icon
	return no
}

// WithCallbackContext sets callback context data
func (no *NotifyOptions) WithCallbackContext(context string) *NotifyOptions {
	no.CallbackContext = context
	return no
}

// WithCallbackTarget sets the URL to open on click
func (no *NotifyOptions) WithCallbackTarget(target string) *NotifyOptions {
	no.CallbackTarget = target
	return no
}

// Close closes the callback listener
func (c *Client) Close() error {
	if c.callbackListener != nil {
		return c.callbackListener.Close()
	}
	return nil
}