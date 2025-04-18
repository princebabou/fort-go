package utils

import (
	"net"
	"net/url"
	"strings"
	"time"
	"strconv"
)

// ValidateIP checks if the given string is a valid IP address
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ValidatePort checks if the given port is in a valid range
func ValidatePort(port int) bool {
	return port > 0 && port < 65536
}

// ValidateURL checks if the given string is a valid URL
func ValidateURL(urlStr string) bool {
	// Add http:// prefix if not present for validation purposes
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "http://" + urlStr
	}
	
	u, err := url.ParseRequestURI(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// FormatDuration formats a time duration to a readable string
func FormatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	
	if h > 0 {
		return formatTime(int(h), "hour") + ", " + formatTime(int(m), "minute") + ", " + formatTime(int(s), "second")
	}
	if m > 0 {
		return formatTime(int(m), "minute") + ", " + formatTime(int(s), "second")
	}
	return formatTime(int(s), "second")
}

// SanitizeInput sanitizes user input to prevent command injection
func SanitizeInput(input string) string {
	// Replace potentially dangerous characters
	sanitized := strings.ReplaceAll(input, ";", "")
	sanitized = strings.ReplaceAll(sanitized, "&", "")
	sanitized = strings.ReplaceAll(sanitized, "|", "")
	sanitized = strings.ReplaceAll(sanitized, ">", "")
	sanitized = strings.ReplaceAll(sanitized, "<", "")
	sanitized = strings.ReplaceAll(sanitized, "$", "")
	sanitized = strings.ReplaceAll(sanitized, "`", "")
	sanitized = strings.ReplaceAll(sanitized, "\\", "")
	
	return strings.TrimSpace(sanitized)
}

// Helper function to format time units with proper pluralization
func formatTime(value int, unit string) string {
	if value == 1 {
		return "1 " + unit
	}
	return strconv.Itoa(value) + " " + unit + "s"
} 