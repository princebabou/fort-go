# FortiCore

FortiCore is an automated Penetration Testing Tool (PTT) designed to simplify penetration testing processes. It offers a command-line-based interface for streamlined, automated vulnerability scanning, enabling businesses to efficiently identify potential security threats without the need for advanced technical expertise.

## Features

- **User-Friendly CLI Interface**: Simple command-line interface that allows users to select systems to test, configure scan parameters, and view results.
- **Automated Vulnerability Scanning**: Automatically scans for common vulnerabilities (e.g., SQL injection, XSS, open ports) on user-selected targets.
- **Safe Exploitation**: Automatically exploits identified vulnerabilities in a non-destructive manner to demonstrate the risks and simulate an attack.
- **Detailed Reports**: Generates comprehensive reports with identified vulnerabilities, exploitation attempts, and remediation guidance.

## Installation

### Prerequisites

- Go 1.20 or later
- Git

### Building from source

```bash
git clone https://github.com/princebabou/fort-go.git
cd fort-go
go build -o fort ./cmd/fort
```

## Usage

FortiCore uses `fort` as the command trigger. Here are some examples of how to use it:

### Basic Scanning

```bash
# Perform a full scan on a target
fort scan -t example.com

# Perform a network scan on a specific port range
fort scan -t 192.168.1.1 -y network -p 1-1000

# Perform a web scan with recursive option
fort scan -t https://example.com -y web -r
```

### Exploitation

```bash
# Automatically exploit vulnerabilities found in a previous scan
fort exploit -t example.com

# Manually exploit with a specific payload
fort exploit -t example.com -y manual -p "' OR 1=1 --"
```

### Reporting

```bash
# Generate a PDF report from scan results
fort report -i scan_results.json -f pdf -o report.pdf

# Generate an HTML report
fort report -i scan_results.json -f html -o report.html
```

## Command Options

### Global Options

- `-t, --target`: Target to scan (IP, domain, or URL)
- `-v, --verbose`: Enable verbose output
- `-o, --output`: Output file for reports

### Scan Options

- `-y, --type`: Type of scan to perform (network, web, full)
- `-p, --port`: Port range to scan (network scan only)
- `--timeout`: Timeout in seconds for each scan operation
- `--threads`: Number of concurrent threads to use
- `-r, --recursive`: Recursively scan web applications (web scan only)

### Exploit Options

- `-y, --type`: Type of exploitation (auto, manual)
- `-s, --safe`: Enable safe mode (non-destructive)
- `-p, --payload`: Custom payload for manual exploitation
- `--timeout`: Timeout in seconds for exploitation attempts

### Report Options

- `-f, --format`: Report format (pdf, html, txt)
- `--template`: Custom report template file
- `-i, --input`: Input file with scan results
- `--sign`: Digitally sign the report (PDF only)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Disclaimer

FortiCore is designed for legitimate security testing only. Always ensure you have proper authorization before scanning or testing any system or network. Unauthorized scanning or exploitation of systems is illegal and unethical.
