# go-gntp

[![Go Reference](https://pkg.go.dev/badge/github.com/cumulus13/go-gntp.svg)](https://pkg.go.dev/github.com/cumulus13/go-gntp)
[![Go Report Card](https://goreportcard.com/badge/github.com/cumulus13/go-gntp)](https://goreportcard.com/report/github.com/cumulus13/go-gntp)
[![Go CI](https://github.com/cumulus13/go-gntp/actions/workflows/go.yml/badge.svg)](https://github.com/cumulus13/go-gntp/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/cumulus13/go-gntp)](https://github.com/cumulus13/go-gntp/blob/master/go.mod)

GNTP (Growl Notification Transport Protocol) client for Go with **full callback support**, Windows/Linux/Mac/Android compatibility, and multiple icon delivery modes.

## ✨ Features

- ✅ **Full GNTP 1.0 protocol implementation**
- ✅ **Callback support** (click, close, timeout events)
- ✅ **Multiple icon delivery modes** (Binary, File URL, Data URL, HTTP URL)
- ✅ **Windows Growl compatibility** with automatic workarounds
- ✅ **Android Growl compatibility** (tested!)
- ✅ **Cross-platform support** (Windows, macOS, Linux, Android)
- ✅ **Simple Message struct API** (gntplib-compatible)
- ✅ **Resource deduplication** to prevent errors
- ✅ **Zero external dependencies** (except uuid)

## 📦 Installation

```bash
go get github.com/cumulus13/go-gntp
```

## 🚀 Quick Start

### Basic Notification

```go
package main

import (
    "log"
    "github.com/cumulus13/go-gntp"
)

func main() {
    // Create client
    client := gntp.NewClient("My App")
    
    // Define notification type
    notification := gntp.NewNotificationType("alert").
        WithDisplayName("Alert")
    
    // Register
    if err := client.Register([]*gntp.NotificationType{notification}); err != nil {
        log.Fatal(err)
    }
    
    // Send notification
    if err := client.Notify("alert", "Hello", "Test notification"); err != nil {
        log.Fatal(err)
    }
}
```

### With Icon

```go
// Load icon
icon, err := gntp.LoadResource("icon.png")
if err != nil {
    log.Fatal(err)
}

client := gntp.NewClient("My App").
    WithIconMode(gntp.IconModeDataURL)  // Best for Windows/Android

notification := gntp.NewNotificationType("alert").
    WithIcon(icon)

client.Register([]*gntp.NotificationType{notification})
client.Notify("alert", "Hello", "With icon!")
```

### With Callbacks (IMPORTANT!)

```go
client := gntp.NewClient("Callback App")

// Set up callback handler
err := client.WithCallback(func(info gntp.CallbackInfo) {
    switch info.Type {
    case gntp.CallbackClick:
        fmt.Printf("User clicked! Context: %s\n", info.Context)
    case gntp.CallbackClose:
        fmt.Println("User closed notification")
    case gntp.CallbackTimeout:
        fmt.Println("Notification timed out")
    }
})
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// Register
notification := gntp.NewNotificationType("alert")
client.Register([]*gntp.NotificationType{notification})

// Send with callback options
options := gntp.NewNotifyOptions().
    WithSticky(true).
    WithCallbackContext("user_data_123").
    WithCallbackTarget("https://example.com")  // URL to open on click

client.NotifyWithOptions("alert", "Click Me!", "Message", options)

// Wait for callbacks
select {}
```

### Using Message Struct (gntplib-compatible)

```go
client := gntp.NewClient("Simple App")

msg := &gntp.Message{
    Event:       "alert",
    Title:       "Title",
    Text:        "Message",
    Icon:        "icon.png",
    Callback:    "https://example.com",
    DisplayName: "Alert",
    Sticky:      true,
    Priority:    1,
}

client.SendMessage(msg)
```

## 🎯 Icon Delivery Modes

### Binary Mode (Recommended for Windows!)

```go
client := gntp.NewClient("App").
    WithIconMode(gntp.IconModeBinary)
```

**Format:** `x-growl-resource://UUID` + binary data
**Best for:** Windows Growl, macOS, Linux (MOST RELIABLE!)
**✅ Tested and working on Windows Growl**

### DataURL Mode

```go
client := gntp.NewClient("App").
    WithIconMode(gntp.IconModeDataURL)
```

**Format:** `data:image/png;base64,iVBORw0KGgo...`
**Best for:** Android, when you can't use binary
**⚠️ May have issues with large icons on Windows**

### File URL Mode

```go
client := gntp.NewClient("App").
    WithIconMode(gntp.IconModeFileURL)
```

**Format:** `file:///C:/full/path/to/icon.png`
**Best for:** Shared icons on disk
**⚠️ Requires absolute path, file must be accessible to Growl**

### HTTP URL Mode

```go
client := gntp.NewClient("App").
    WithIconMode(gntp.IconModeHttpURL)
```

**Format:** `http://example.com/icon.png`
**Best for:** Web-hosted icons, remote servers

## 🔔 Callback Events

Callbacks are triggered when users interact with notifications:

- **`CallbackClick`** - User clicked the notification
- **`CallbackClose`** - User closed the notification  
- **`CallbackTimeout`** - Notification timed out

### Callback Example

```go
client.WithCallback(func(info gntp.CallbackInfo) {
    fmt.Printf("Event: %s\n", info.Type)
    fmt.Printf("ID: %s\n", info.NotificationID)
    fmt.Printf("Context: %s\n", info.Context)
    fmt.Printf("Time: %s\n", info.Timestamp)
})
```

### Callback with Context Data

```go
options := gntp.NewNotifyOptions().
    WithCallbackContext("order_id:12345")  // Pass custom data

client.NotifyWithOptions("alert", "Order Ready", "Your order is ready!", options)

// In callback handler:
// info.Context will be "order_id:12345"
```

## 📋 API Reference

### Client Methods

```go
client := gntp.NewClient("App Name")
client.WithHost("192.168.1.100")                    // Set remote host
client.WithPort(23053)                              // Set port
client.WithIconMode(gntp.IconModeDataURL)           // Set icon mode
client.WithIcon(icon)                               // Set app icon
client.WithDebug(true)                              // Enable debug
client.WithTimeout(10 * time.Second)                // Set timeout
client.WithCallback(handler)                        // Set callback handler
client.Register(notifications)                      // Register app
client.Notify(name, title, text)                    // Send notification
client.NotifyWithOptions(name, title, text, opts)   // Send with options
client.SendMessage(msg)                             // Send via Message struct
client.Close()                                      // Close callback listener
```

### NotifyOptions

```go
opts := gntp.NewNotifyOptions()
opts.WithSticky(true)                               // Sticky notification
opts.WithPriority(2)                                // Priority (-2 to 2)
opts.WithIcon(icon)                                 // Per-notification icon
opts.WithCallbackContext("custom_data")             // Callback context
opts.WithCallbackTarget("https://example.com")      // URL to open
```

## 🌍 Platform Compatibility

| Platform | Binary Mode | DataURL Mode | FileURL Mode | Callbacks | Recommended |
|----------|-------------|--------------|--------------|-----------|-------------|
| **Windows (Growl for Windows)** | ✅ **Works!** | ⚠️ Issues | ⚠️ Issues | ✅ Works | **Binary** |
| **macOS (Growl)** | ✅ Works | ✅ Works | ✅ Works | ✅ Works | Binary |
| **Linux (Growl-compatible)** | ✅ Works | ✅ Works | ✅ Works | ✅ Works | Binary |
| **Android (Growl for Android)** | ✅ Works | ✅ Works | ⚠️ Issues | ✅ Works | **Binary or DataURL** |

**Testing Results:**
- ✅ **Binary Mode**: Confirmed working on Windows Growl for Windows
- ⚠️ **DataURL Mode**: May fail with large icons (base64 size limit)
- ⚠️ **FileURL Mode**: Requires absolute path, may have permission issues

## 📚 Examples

Run examples:

```bash
# Basic notification
go run examples/basic/main.go

# With icon
go run examples/with_icon/main.go

# With callback
go run examples/callback/main.go

# Message struct
go run examples/message/main.go

# Remote Android
GROWL_HOST=192.168.1.50 go run examples/android/main.go
```

## 🔧 Advanced Usage

### Remote Notifications

```go
client := gntp.NewClient("Remote App").
    WithHost("192.168.1.100").
    WithPort(23053)
```

### Multiple Notification Types

```go
info := gntp.NewNotificationType("info").WithDisplayName("Information")
warning := gntp.NewNotificationType("warning").WithDisplayName("Warning")
error := gntp.NewNotificationType("error").WithDisplayName("Error")

client.Register([]*gntp.NotificationType{info, warning, error})

client.Notify("info", "Info", "Something happened")
client.Notify("warning", "Warning", "Be careful!")
client.Notify("error", "Error", "Something went wrong!")
```

### Loading Icons

```go
// From file
icon, err := gntp.LoadResource("icon.png")

// From bytes
data, _ := os.ReadFile("icon.png")
icon := gntp.LoadResourceFromBytes(data, "image/png")
```

## 🐛 Troubleshooting

### Icon Not Showing

Try different icon modes:

```go
// Try DataURL mode (most compatible)
client.WithIconMode(gntp.IconModeDataURL)

// Or Binary mode (spec compliant)
client.WithIconMode(gntp.IconModeBinary)
```

### Callbacks Not Working

Make sure:
1. Callback handler is set BEFORE Register()
2. Keep program running with `select {}` or similar
3. Close client properly with `defer client.Close()`
4. Firewall allows incoming connections on callback port

### Android Connection Issues

```go
// Use DataURL mode for best Android compatibility
client := gntp.NewClient("Android App").
    WithHost("192.168.1.50").
    WithIconMode(gntp.IconModeDataURL).
    WithTimeout(15 * time.Second)  // Longer timeout for mobile
```

## 🤝 Contributing

Contributions welcome! Please open an issue or PR.

## 📄 License

MIT License - See LICENSE file for details.

## 👤 Author

[Hadi Cahyadi](mailto:cumulus13@gmail.com)
    
[![Buy Me a Coffee](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/cumulus13)

[![Donate via Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/cumulus13)
 
[Support me on Patreon](https://www.patreon.com/cumulus13)

## 🙏 Related Projects

- **Rust gntp**: https://github.com/cumulus13/gntp
- **Python gntplib**: https://github.com/cumulus13/gntplib

---

**Note:** This library implements the full GNTP 1.0 specification with production-ready Windows/Linux/Mac/Android compatibility and callback support!
