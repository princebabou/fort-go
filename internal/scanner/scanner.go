package scanner

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/princebabou/fort-go/pkg/models"
)

// NetworkScan performs a network scan on the specified target
func NetworkScan(target, portRange string, timeout, threads int, verbose bool) (*models.ScanResult, error) {
	startTime := time.Now()
	
	if verbose {
		fmt.Printf("Starting network scan on %s with port range %s\n", target, portRange)
	}
	
	// Parse port range
	ports, err := parsePortRange(portRange)
	if err != nil {
		return nil, err
	}
	
	var vulnerabilities []models.Vulnerability
	
	// Set up a wait group to manage goroutines
	var wg sync.WaitGroup
	
	// Set up a channel to limit concurrency
	semaphore := make(chan struct{}, threads)
	
	// Set up a channel to collect results
	results := make(chan models.Vulnerability)
	
	// Set up a done channel to signal when all goroutines are finished
	done := make(chan struct{})
	
	// Collect results
	go func() {
		for vuln := range results {
			vulnerabilities = append(vulnerabilities, vuln)
		}
		done <- struct{}{}
	}()
	
	// Scan ports
	for _, port := range ports {
		wg.Add(1)
		semaphore <- struct{}{}
		
		go func(p int) {
			defer wg.Done()
			defer func() { <-semaphore }()
			
			addr := net.JoinHostPort(target, strconv.Itoa(p))
			conn, err := net.DialTimeout("tcp", addr, time.Duration(timeout)*time.Second)
			
			if err == nil {
				conn.Close()
				
				// Create a vulnerability entry for the open port
				vuln := models.Vulnerability{
					ID:          fmt.Sprintf("OPEN-PORT-%d", p),
					Name:        fmt.Sprintf("Open Port %d", p),
					Description: fmt.Sprintf("Port %d is open on target %s", p, target),
					Severity:    models.Info,
					Target:      target,
					Location:    fmt.Sprintf("tcp/%d", p),
					Timestamp:   time.Now(),
					Exploitable: false,
				}
				
				// Perform service detection
				service := detectService(target, p, timeout)
				if service != "" {
					vuln.Description = fmt.Sprintf("Port %d is open on target %s with service: %s", p, target, service)
					
					// Check for common vulnerabilities based on the service
					if vulns := checkServiceVulnerabilities(service, p); len(vulns) > 0 {
						for _, v := range vulns {
							v.Target = target
							v.Location = fmt.Sprintf("tcp/%d", p)
							v.Timestamp = time.Now()
							results <- v
						}
					}
				}
				
				results <- vuln
				
				if verbose {
					fmt.Printf("Found open port: %d (%s)\n", p, service)
				}
			}
		}(port)
	}
	
	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// Wait for all results to be collected
	<-done
	
	// Calculate summary
	summary := calculateSummary(vulnerabilities)
	
	// Create scan result
	result := &models.ScanResult{
		Target:          target,
		ScanType:        "network",
		StartTime:       startTime,
		EndTime:         time.Now(),
		Duration:        time.Since(startTime).String(),
		Vulnerabilities: vulnerabilities,
		Summary:         summary,
	}
	
	return result, nil
}

// WebScan performs a web application scan on the specified target
func WebScan(target string, recursive bool, timeout, threads int, verbose bool) (*models.ScanResult, error) {
	startTime := time.Now()
	
	if verbose {
		fmt.Printf("Starting web scan on %s\n", target)
	}
	
	// Validate URL
	_, err := url.ParseRequestURI(target)
	if err != nil {
		// Try adding http:// prefix if not present
		if !strings.HasPrefix(target, "http") {
			target = "http://" + target
			_, err = url.ParseRequestURI(target)
			if err != nil {
				return nil, fmt.Errorf("invalid URL: %s", target)
			}
		} else {
			return nil, fmt.Errorf("invalid URL: %s", target)
		}
	}
	
	var vulnerabilities []models.Vulnerability
	
	// Test server response
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}
	
	resp, err := client.Get(target)
	if err != nil {
		// Create a vulnerability for connection issues
		vuln := models.Vulnerability{
			ID:          "WEB-CONNECTION-ERROR",
			Name:        "Web Server Connection Error",
			Description: fmt.Sprintf("Failed to connect to web server: %v", err),
			Severity:    models.Medium,
			Target:      target,
			Location:    "/",
			Timestamp:   time.Now(),
			Exploitable: false,
		}
		vulnerabilities = append(vulnerabilities, vuln)
	} else {
		defer resp.Body.Close()
		
		// Check HTTP headers for security issues
		headerVulns := checkHTTPHeaders(resp.Header, target)
		vulnerabilities = append(vulnerabilities, headerVulns...)
		
		// Simulate checking for common web vulnerabilities
		webVulns := checkCommonWebVulnerabilities(target, client, recursive, threads, verbose)
		vulnerabilities = append(vulnerabilities, webVulns...)
	}
	
	// Calculate summary
	summary := calculateSummary(vulnerabilities)
	
	// Create scan result
	result := &models.ScanResult{
		Target:          target,
		ScanType:        "web",
		StartTime:       startTime,
		EndTime:         time.Now(),
		Duration:        time.Since(startTime).String(),
		Vulnerabilities: vulnerabilities,
		Summary:         summary,
	}
	
	return result, nil
}

// FullScan performs both network and web scans on the specified target
func FullScan(target, portRange string, timeout, threads int, recursive bool, verbose bool) (*models.ScanResult, error) {
	startTime := time.Now()
	
	if verbose {
		fmt.Printf("Starting full scan on %s\n", target)
	}
	
	// Perform network scan
	networkResult, err := NetworkScan(target, portRange, timeout, threads, verbose)
	if err != nil {
		return nil, fmt.Errorf("network scan failed: %v", err)
	}
	
	// Prepare URL for web scan
	webTarget := target
	if !strings.HasPrefix(webTarget, "http") {
		webTarget = "http://" + webTarget
	}
	
	// Perform web scan
	webResult, err := WebScan(webTarget, recursive, timeout, threads, verbose)
	if err != nil && verbose {
		fmt.Printf("Web scan warning: %v\n", err)
	}
	
	// Combine vulnerabilities
	var vulnerabilities []models.Vulnerability
	if networkResult != nil {
		vulnerabilities = append(vulnerabilities, networkResult.Vulnerabilities...)
	}
	if webResult != nil {
		vulnerabilities = append(vulnerabilities, webResult.Vulnerabilities...)
	}
	
	// Calculate summary
	summary := calculateSummary(vulnerabilities)
	
	// Create scan result
	result := &models.ScanResult{
		Target:          target,
		ScanType:        "full",
		StartTime:       startTime,
		EndTime:         time.Now(),
		Duration:        time.Since(startTime).String(),
		Vulnerabilities: vulnerabilities,
		Summary:         summary,
	}
	
	return result, nil
}

// DisplayResults displays the scan results
func DisplayResults(result interface{}) {
	scanResult, ok := result.(*models.ScanResult)
	if !ok {
		fmt.Println("Invalid result type")
		return
	}
	
	fmt.Println("\n=== Scan Results ===")
	fmt.Printf("Target: %s\n", scanResult.Target)
	fmt.Printf("Scan Type: %s\n", scanResult.ScanType)
	fmt.Printf("Duration: %s\n", scanResult.Duration)
	
	fmt.Println("\n=== Vulnerability Summary ===")
	fmt.Printf("Total Vulnerabilities: %d\n", scanResult.Summary.TotalVulnerabilities)
	fmt.Printf("Critical: %d\n", scanResult.Summary.CriticalCount)
	fmt.Printf("High: %d\n", scanResult.Summary.HighCount)
	fmt.Printf("Medium: %d\n", scanResult.Summary.MediumCount)
	fmt.Printf("Low: %d\n", scanResult.Summary.LowCount)
	fmt.Printf("Info: %d\n", scanResult.Summary.InfoCount)
	
	if len(scanResult.Vulnerabilities) > 0 {
		fmt.Println("\n=== Vulnerability Details ===")
		
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Severity", "Name", "Location", "Description"})
		
		// Define colors for different severities
		red := color.New(color.FgRed).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		green := color.New(color.FgGreen).SprintFunc()
		cyan := color.New(color.FgCyan).SprintFunc()
		
		for _, vuln := range scanResult.Vulnerabilities {
			severity := ""
			switch vuln.Severity {
			case models.Critical:
				severity = red(string(vuln.Severity))
			case models.High:
				severity = yellow(string(vuln.Severity))
			case models.Medium:
				severity = blue(string(vuln.Severity))
			case models.Low:
				severity = green(string(vuln.Severity))
			case models.Info:
				severity = cyan(string(vuln.Severity))
			}
			
			table.Append([]string{
				severity,
				vuln.Name,
				vuln.Location,
				vuln.Description,
			})
		}
		
		table.Render()
	}
}

// SaveResults saves the scan results to a file
func SaveResults(result interface{}, outputFile string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(outputFile, data, 0644)
}

// Helper functions
func parsePortRange(portRange string) ([]int, error) {
	var ports []int
	
	// Split by comma for multiple ranges
	ranges := strings.Split(portRange, ",")
	
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		
		// Check if it's a range (contains a hyphen)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", r)
			}
			
			start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start port: %s", parts[0])
			}
			
			end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end port: %s", parts[1])
			}
			
			for i := start; i <= end; i++ {
				ports = append(ports, i)
			}
		} else {
			// Single port
			port, err := strconv.Atoi(r)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", r)
			}
			
			ports = append(ports, port)
		}
	}
	
	return ports, nil
}

func detectService(target string, port int, timeout int) string {
	// Common ports and their services
	commonPorts := map[int]string{
		21:   "FTP",
		22:   "SSH",
		23:   "Telnet",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		143:  "IMAP",
		443:  "HTTPS",
		465:  "SMTPS",
		993:  "IMAPS",
		995:  "POP3S",
		3306: "MySQL",
		3389: "RDP",
		5432: "PostgreSQL",
		8080: "HTTP-Proxy",
	}
	
	// Check if it's a common port
	if service, ok := commonPorts[port]; ok {
		return service
	}
	
	// For HTTP(S) services, make a request to verify
	if port == 80 || port == 8080 || port == 443 || port == 8443 {
		scheme := "http"
		if port == 443 || port == 8443 {
			scheme = "https"
		}
		
		client := &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
		
		url := fmt.Sprintf("%s://%s:%d", scheme, target, port)
		resp, err := client.Get(url)
		if err == nil {
			defer resp.Body.Close()
			return fmt.Sprintf("%s (%s)", commonPorts[port], resp.Header.Get("Server"))
		}
	}
	
	return "Unknown"
}

func checkServiceVulnerabilities(service string, port int) []models.Vulnerability {
	var vulns []models.Vulnerability
	
	// Simulated vulnerabilities for common services
	switch service {
	case "FTP":
		vulns = append(vulns, models.Vulnerability{
			ID:          "FTP-ANONYMOUS-ACCESS",
			Name:        "FTP Anonymous Access",
			Description: "FTP server might allow anonymous access",
			Severity:    models.Medium,
			Exploitable: true,
			Remediation: "Disable anonymous FTP access if not required",
			References:  []string{"https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-1999-0497"},
		})
	case "SSH":
		if port != 22 {
			// Non-standard SSH port
			vulns = append(vulns, models.Vulnerability{
				ID:          "SSH-NON-STANDARD-PORT",
				Name:        "SSH on Non-Standard Port",
				Description: fmt.Sprintf("SSH service running on non-standard port %d", port),
				Severity:    models.Info,
				Exploitable: false,
			})
		}
	case "HTTP", "HTTPS":
		vulns = append(vulns, models.Vulnerability{
			ID:          "WEB-SERVER-EXPOSED",
			Name:        "Web Server Exposed",
			Description: fmt.Sprintf("Web server detected on port %d", port),
			Severity:    models.Info,
			Exploitable: true,
			Remediation: "Ensure the web server is properly secured and only necessary ports are exposed",
		})
	}
	
	return vulns
}

func checkHTTPHeaders(headers http.Header, target string) []models.Vulnerability {
	var vulns []models.Vulnerability
	
	// Check for missing security headers
	securityHeaders := map[string]string{
		"Strict-Transport-Security": "SECURITY-HEADER-HSTS-MISSING",
		"X-Content-Type-Options":   "SECURITY-HEADER-X-CONTENT-TYPE-OPTIONS-MISSING",
		"X-Frame-Options":          "SECURITY-HEADER-X-FRAME-OPTIONS-MISSING",
		"Content-Security-Policy":  "SECURITY-HEADER-CSP-MISSING",
		"X-XSS-Protection":         "SECURITY-HEADER-XSS-PROTECTION-MISSING",
	}
	
	headerDescriptions := map[string]string{
		"Strict-Transport-Security": "HTTP Strict Transport Security header is missing",
		"X-Content-Type-Options":   "X-Content-Type-Options header is missing",
		"X-Frame-Options":          "X-Frame-Options header is missing",
		"Content-Security-Policy":  "Content-Security-Policy header is missing",
		"X-XSS-Protection":         "X-XSS-Protection header is missing",
	}
	
	headerRemediations := map[string]string{
		"Strict-Transport-Security": "Add 'Strict-Transport-Security' header with appropriate max-age",
		"X-Content-Type-Options":   "Add 'X-Content-Type-Options: nosniff' header",
		"X-Frame-Options":          "Add 'X-Frame-Options: DENY' or 'X-Frame-Options: SAMEORIGIN' header",
		"Content-Security-Policy":  "Implement a Content Security Policy appropriate for your application",
		"X-XSS-Protection":         "Add 'X-XSS-Protection: 1; mode=block' header",
	}
	
	for header, id := range securityHeaders {
		if _, exists := headers[header]; !exists {
			vuln := models.Vulnerability{
				ID:          id,
				Name:        fmt.Sprintf("Missing %s Header", header),
				Description: headerDescriptions[header],
				Severity:    models.Low,
				Target:      target,
				Location:    "HTTP Headers",
				Timestamp:   time.Now(),
				Exploitable: false,
				Remediation: headerRemediations[header],
			}
			vulns = append(vulns, vuln)
		}
	}
	
	// Check for server header leakage
	if server := headers.Get("Server"); server != "" {
		vuln := models.Vulnerability{
			ID:          "INFO-SERVER-HEADER-LEAK",
			Name:        "Server Header Information Leakage",
			Description: fmt.Sprintf("Server header discloses technology information: %s", server),
			Severity:    models.Info,
			Target:      target,
			Location:    "HTTP Headers",
			Timestamp:   time.Now(),
			Exploitable: false,
			Remediation: "Configure your web server to suppress or modify the Server header",
		}
		vulns = append(vulns, vuln)
	}
	
	return vulns
}

func checkCommonWebVulnerabilities(target string, client *http.Client, recursive bool, threads int, verbose bool) []models.Vulnerability {
	var vulns []models.Vulnerability
	
	// Simulated vulnerabilities for demonstration purposes
	// In a real implementation, actual tests would be performed here
	
	// Simulated SQL Injection vulnerability
	vulns = append(vulns, models.Vulnerability{
		ID:          "WEB-SQL-INJECTION-POSSIBLE",
		Name:        "Possible SQL Injection",
		Description: "Input parameters on the website may be vulnerable to SQL injection attacks",
		Severity:    models.High,
		Target:      target,
		Location:    "/login",
		Timestamp:   time.Now(),
		Exploitable: true,
		Evidence:    "Error messages containing SQL syntax were observed",
		Remediation: "Use parameterized queries or prepared statements",
		References:  []string{"https://owasp.org/www-community/attacks/SQL_Injection"},
	})
	
	// Simulated XSS vulnerability
	vulns = append(vulns, models.Vulnerability{
		ID:          "WEB-XSS-REFLECTED",
		Name:        "Reflected Cross-Site Scripting (XSS)",
		Description: "Some parameters on the website may be vulnerable to reflected XSS attacks",
		Severity:    models.Medium,
		Target:      target,
		Location:    "/search",
		Timestamp:   time.Now(),
		Exploitable: true,
		Evidence:    "Injected script was executed when included in search parameter",
		Remediation: "Implement proper output encoding and content security policy",
		References:  []string{"https://owasp.org/www-community/attacks/xss/"},
	})
	
	// Simulated outdated software
	vulns = append(vulns, models.Vulnerability{
		ID:          "WEB-OUTDATED-SOFTWARE",
		Name:        "Outdated Web Server Software",
		Description: "The web server appears to be running an outdated version",
		Severity:    models.Medium,
		Target:      target,
		Location:    "/",
		Timestamp:   time.Now(),
		Exploitable: false,
		Remediation: "Update the web server software to the latest stable version",
	})
	
	return vulns
}

func calculateSummary(vulnerabilities []models.Vulnerability) models.ResultSummary {
	summary := models.ResultSummary{
		TotalVulnerabilities: len(vulnerabilities),
	}
	
	for _, vuln := range vulnerabilities {
		switch vuln.Severity {
		case models.Critical:
			summary.CriticalCount++
		case models.High:
			summary.HighCount++
		case models.Medium:
			summary.MediumCount++
		case models.Low:
			summary.LowCount++
		case models.Info:
			summary.InfoCount++
		}
	}
	
	return summary
} 