package report

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"strconv"
	"strings"
	
	"github.com/princebabou/fort-go/pkg/models"
)

// GenerateTimestamp returns a formatted timestamp for use in filenames
func GenerateTimestamp() string {
	return time.Now().Format("20060102_150405")
}

// GeneratePDFReport generates a PDF report from scan or exploit results
func GeneratePDFReport(inputFile, outputFile, templateFile string, digitallySigned bool, verbose bool) error {
	if verbose {
		fmt.Printf("Generating PDF report from %s to %s\n", inputFile, outputFile)
		if templateFile != "" {
			fmt.Printf("Using custom template: %s\n", templateFile)
		}
		fmt.Printf("Digital signature: %v\n", digitallySigned)
	}
	
	// Load data from input file
	data, err := loadDataFromFile(inputFile)
	if err != nil {
		return err
	}
	
	// In a real implementation, this would use a PDF generation library
	// For demonstration purposes, we'll just write a text representation to the output file
	
	report := generateReportText(data, "PDF")
	
	// Add digital signature information if requested
	if digitallySigned {
		report += "\n\n[This report has been digitally signed to prevent tampering]\n"
		report += fmt.Sprintf("Digital Signature: %s\n", generateMockSignature(report))
		report += fmt.Sprintf("Timestamp: %s\n", time.Now().Format(time.RFC3339))
	}
	
	// Write to file
	if err := os.WriteFile(outputFile, []byte(report), 0644); err != nil {
		return err
	}
	
	fmt.Printf("PDF report generated successfully: %s\n", outputFile)
	return nil
}

// GenerateHTMLReport generates an HTML report from scan or exploit results
func GenerateHTMLReport(inputFile, outputFile, templateFile string, verbose bool) error {
	if verbose {
		fmt.Printf("Generating HTML report from %s to %s\n", inputFile, outputFile)
		if templateFile != "" {
			fmt.Printf("Using custom template: %s\n", templateFile)
		}
	}
	
	// Load data from input file
	data, err := loadDataFromFile(inputFile)
	if err != nil {
		return err
	}
	
	// In a real implementation, this would use an HTML template
	// For demonstration purposes, we'll just generate a simple HTML file
	
	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	html += "<title>FortiCore Security Report</title>\n"
	html += "<style>\n"
	html += "body { font-family: Arial, sans-serif; margin: 20px; }\n"
	html += "h1, h2 { color: #336699; }\n"
	html += "table { border-collapse: collapse; width: 100%; margin-top: 20px; }\n"
	html += "th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }\n"
	html += "th { background-color: #f2f2f2; }\n"
	html += ".critical { color: #d9534f; }\n"
	html += ".high { color: #f0ad4e; }\n"
	html += ".medium { color: #5bc0de; }\n"
	html += ".low { color: #5cb85c; }\n"
	html += ".info { color: #5bc0de; }\n"
	html += "</style>\n"
	html += "</head>\n<body>\n"
	
	html += "<h1>FortiCore Security Report</h1>\n"
	html += "<p>Generated: " + time.Now().Format(time.RFC3339) + "</p>\n"
	
	// Add summary and details based on data type
	switch v := data.(type) {
	case *models.ScanResult:
		html += "<h2>Scan Results</h2>\n"
		html += "<p><strong>Target:</strong> " + v.Target + "</p>\n"
		html += "<p><strong>Scan Type:</strong> " + v.ScanType + "</p>\n"
		html += "<p><strong>Duration:</strong> " + v.Duration + "</p>\n"
		
		html += "<h2>Vulnerability Summary</h2>\n"
		html += "<p><strong>Total Vulnerabilities:</strong> " + strconv.Itoa(v.Summary.TotalVulnerabilities) + "</p>\n"
		html += "<p><strong>Critical:</strong> " + strconv.Itoa(v.Summary.CriticalCount) + "</p>\n"
		html += "<p><strong>High:</strong> " + strconv.Itoa(v.Summary.HighCount) + "</p>\n"
		html += "<p><strong>Medium:</strong> " + strconv.Itoa(v.Summary.MediumCount) + "</p>\n"
		html += "<p><strong>Low:</strong> " + strconv.Itoa(v.Summary.LowCount) + "</p>\n"
		html += "<p><strong>Info:</strong> " + strconv.Itoa(v.Summary.InfoCount) + "</p>\n"
		
		if len(v.Vulnerabilities) > 0 {
			html += "<h2>Vulnerability Details</h2>\n"
			html += "<table>\n"
			html += "<tr><th>Severity</th><th>Name</th><th>Location</th><th>Description</th><th>Remediation</th></tr>\n"
			
			for _, vuln := range v.Vulnerabilities {
				severityClass := strings.ToLower(string(vuln.Severity))
				html += "<tr>\n"
				html += fmt.Sprintf("<td class=\"%s\">%s</td>\n", severityClass, vuln.Severity)
				html += "<td>" + vuln.Name + "</td>\n"
				html += "<td>" + vuln.Location + "</td>\n"
				html += "<td>" + vuln.Description + "</td>\n"
				html += "<td>" + vuln.Remediation + "</td>\n"
				html += "</tr>\n"
			}
			
			html += "</table>\n"
		}
	case *models.ExploitResult:
		html += "<h2>Exploitation Results</h2>\n"
		html += "<p><strong>Target:</strong> " + v.Target + "</p>\n"
		html += "<p><strong>Exploit Type:</strong> " + v.ExploitType + "</p>\n"
		html += "<p><strong>Safe Mode:</strong> " + strconv.FormatBool(v.SafeMode) + "</p>\n"
		html += "<p><strong>Duration:</strong> " + v.Duration + "</p>\n"
		html += "<p><strong>Success Count:</strong> " + strconv.Itoa(v.SuccessCount) + "</p>\n"
		html += "<p><strong>Fail Count:</strong> " + strconv.Itoa(v.FailCount) + "</p>\n"
		
		if len(v.Vulnerabilities) > 0 {
			html += "<h2>Exploitation Details</h2>\n"
			html += "<table>\n"
			html += "<tr><th>Status</th><th>Name</th><th>Location</th><th>Details</th></tr>\n"
			
			for _, vuln := range v.Vulnerabilities {
				status := "FAILED"
				statusClass := "critical"
				if vuln.Exploited {
					status = "SUCCESS"
					statusClass = "high"
				}
				
				html += "<tr>\n"
				html += fmt.Sprintf("<td class=\"%s\">%s</td>\n", statusClass, status)
				html += "<td>" + vuln.Name + "</td>\n"
				html += "<td>" + vuln.Location + "</td>\n"
				html += "<td>" + vuln.ExploitInfo + "</td>\n"
				html += "</tr>\n"
			}
			
			html += "</table>\n"
		}
	default:
		html += "<p>Unknown report type</p>\n"
	}
	
	html += "<footer>\n<p>FortiCore - Automated Penetration Testing Tool</p>\n</footer>\n"
	html += "</body>\n</html>"
	
	// Write to file
	if err := os.WriteFile(outputFile, []byte(html), 0644); err != nil {
		return err
	}
	
	fmt.Printf("HTML report generated successfully: %s\n", outputFile)
	return nil
}

// GenerateTextReport generates a plain text report from scan or exploit results
func GenerateTextReport(inputFile, outputFile string, verbose bool) error {
	if verbose {
		fmt.Printf("Generating text report from %s to %s\n", inputFile, outputFile)
	}
	
	// Load data from input file
	data, err := loadDataFromFile(inputFile)
	if err != nil {
		return err
	}
	
	// Generate report text
	report := generateReportText(data, "TEXT")
	
	// Write to file
	if err := os.WriteFile(outputFile, []byte(report), 0644); err != nil {
		return err
	}
	
	fmt.Printf("Text report generated successfully: %s\n", outputFile)
	return nil
}

// Helper functions

// loadDataFromFile loads and parses JSON data from the specified file
func loadDataFromFile(filePath string) (interface{}, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %v", err)
	}
	
	// Try to parse as ScanResult
	var scanResult models.ScanResult
	err1 := json.Unmarshal(fileData, &scanResult)
	if err1 == nil && scanResult.ScanType != "" {
		return &scanResult, nil
	}
	
	// Try to parse as ExploitResult
	var exploitResult models.ExploitResult
	err2 := json.Unmarshal(fileData, &exploitResult)
	if err2 == nil && exploitResult.ExploitType != "" {
		return &exploitResult, nil
	}
	
	return nil, fmt.Errorf("failed to parse input file: %v or %v", err1, err2)
}

// generateReportText generates a text representation of the report data
func generateReportText(data interface{}, reportType string) string {
	var report strings.Builder
	
	report.WriteString("===============================================\n")
	report.WriteString("            FORTICORE SECURITY REPORT          \n")
	report.WriteString("===============================================\n\n")
	report.WriteString(fmt.Sprintf("Report Type: %s\n", reportType))
	report.WriteString(fmt.Sprintf("Generated On: %s\n\n", time.Now().Format(time.RFC3339)))
	
	switch v := data.(type) {
	case *models.ScanResult:
		report.WriteString("SCAN RESULTS\n")
		report.WriteString("------------\n\n")
		report.WriteString(fmt.Sprintf("Target: %s\n", v.Target))
		report.WriteString(fmt.Sprintf("Scan Type: %s\n", v.ScanType))
		report.WriteString(fmt.Sprintf("Start Time: %s\n", v.StartTime.Format(time.RFC3339)))
		report.WriteString(fmt.Sprintf("End Time: %s\n", v.EndTime.Format(time.RFC3339)))
		report.WriteString(fmt.Sprintf("Duration: %s\n\n", v.Duration))
		
		report.WriteString("VULNERABILITY SUMMARY\n")
		report.WriteString("---------------------\n\n")
		report.WriteString(fmt.Sprintf("Total Vulnerabilities: %d\n", v.Summary.TotalVulnerabilities))
		report.WriteString(fmt.Sprintf("Critical: %d\n", v.Summary.CriticalCount))
		report.WriteString(fmt.Sprintf("High: %d\n", v.Summary.HighCount))
		report.WriteString(fmt.Sprintf("Medium: %d\n", v.Summary.MediumCount))
		report.WriteString(fmt.Sprintf("Low: %d\n", v.Summary.LowCount))
		report.WriteString(fmt.Sprintf("Info: %d\n\n", v.Summary.InfoCount))
		
		if len(v.Vulnerabilities) > 0 {
			report.WriteString("VULNERABILITY DETAILS\n")
			report.WriteString("---------------------\n\n")
			
			for i, vuln := range v.Vulnerabilities {
				report.WriteString(fmt.Sprintf("Vulnerability #%d\n", i+1))
				report.WriteString(fmt.Sprintf("  Severity: %s\n", vuln.Severity))
				report.WriteString(fmt.Sprintf("  Name: %s\n", vuln.Name))
				report.WriteString(fmt.Sprintf("  Description: %s\n", vuln.Description))
				report.WriteString(fmt.Sprintf("  Location: %s\n", vuln.Location))
				
				if vuln.CVSSScore > 0 {
					report.WriteString(fmt.Sprintf("  CVSS Score: %.1f\n", vuln.CVSSScore))
				}
				
				if vuln.CVEID != "" {
					report.WriteString(fmt.Sprintf("  CVE ID: %s\n", vuln.CVEID))
				}
				
				if vuln.Remediation != "" {
					report.WriteString(fmt.Sprintf("  Remediation: %s\n", vuln.Remediation))
				}
				
				if len(vuln.References) > 0 {
					report.WriteString("  References:\n")
					for _, ref := range vuln.References {
						report.WriteString(fmt.Sprintf("    - %s\n", ref))
					}
				}
				
				report.WriteString("\n")
			}
		}
	case *models.ExploitResult:
		report.WriteString("EXPLOITATION RESULTS\n")
		report.WriteString("--------------------\n\n")
		report.WriteString(fmt.Sprintf("Target: %s\n", v.Target))
		report.WriteString(fmt.Sprintf("Exploit Type: %s\n", v.ExploitType))
		report.WriteString(fmt.Sprintf("Safe Mode: %v\n", v.SafeMode))
		report.WriteString(fmt.Sprintf("Start Time: %s\n", v.StartTime.Format(time.RFC3339)))
		report.WriteString(fmt.Sprintf("End Time: %s\n", v.EndTime.Format(time.RFC3339)))
		report.WriteString(fmt.Sprintf("Duration: %s\n", v.Duration))
		report.WriteString(fmt.Sprintf("Success Count: %d\n", v.SuccessCount))
		report.WriteString(fmt.Sprintf("Fail Count: %d\n\n", v.FailCount))
		
		if len(v.Vulnerabilities) > 0 {
			report.WriteString("EXPLOITATION DETAILS\n")
			report.WriteString("--------------------\n\n")
			
			for i, vuln := range v.Vulnerabilities {
				status := "FAILED"
				if vuln.Exploited {
					status = "SUCCESS"
				}
				
				report.WriteString(fmt.Sprintf("Exploitation Attempt #%d\n", i+1))
				report.WriteString(fmt.Sprintf("  Status: %s\n", status))
				report.WriteString(fmt.Sprintf("  Name: %s\n", vuln.Name))
				report.WriteString(fmt.Sprintf("  Description: %s\n", vuln.Description))
				report.WriteString(fmt.Sprintf("  Location: %s\n", vuln.Location))
				
				if vuln.ExploitInfo != "" {
					report.WriteString(fmt.Sprintf("  Details: %s\n", vuln.ExploitInfo))
				}
				
				if vuln.Evidence != "" {
					report.WriteString(fmt.Sprintf("  Evidence: %s\n", vuln.Evidence))
				}
				
				if vuln.Remediation != "" {
					report.WriteString(fmt.Sprintf("  Remediation: %s\n", vuln.Remediation))
				}
				
				report.WriteString("\n")
			}
		}
	default:
		report.WriteString("Unknown report type\n")
	}
	
	report.WriteString("===============================================\n")
	report.WriteString("       GENERATED BY FORTICORE - FORT-GO        \n")
	report.WriteString("===============================================\n")
	
	return report.String()
}

// generateMockSignature generates a mock digital signature for demonstration purposes
func generateMockSignature(data string) string {
	// In a real implementation, this would use a proper digital signature algorithm
	// For demonstration purposes, we'll just create a mock signature
	return fmt.Sprintf("MOCK-SIG-%x", time.Now().UnixNano())
} 