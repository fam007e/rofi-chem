# Contributing to Rofi-Chem

Thank you for your interest in contributing to Rofi-Chem! We welcome contributions from everyone.

## Branching Strategy

This project uses a dual-branch strategy for clarity and performance:

- **`main`**: The primary branch for users. It contains the Go source code and a pre-built, embedded chemical database. Focus here is on the Rofi interface, performance, and user experience.
- **`data`**: The enrichment branch. It contains Python scripts and tools used to fetch, process, and enrich the chemical database from scientific sources like PubChem and Mendeleev. Focus here is on data engineering and expansion.

## How to Contribute

### Reporting Bugs
Please use the [Bug Report](.github/ISSUE_TEMPLATE/bug_report.md) template when opening an issue.

### Suggesting Features
Please use the [Feature Request](.github/ISSUE_TEMPLATE/feature_request.md) template.

### Pull Requests
1. Fork the repository.
2. Create a new branch for your feature or fix.
3. If your change is data-related, target the `data` branch.
4. If your change is code-related (Go), target the `main` branch.
5. Ensure your code follows the existing style and is well-documented.
6. Submit a Pull Request with a clear description of the changes.

## Development Setup

### Go (Main Branch)
- Requires Go 1.22+.
- `go build -o rofi-chem cmd/rofi-chem/main.go`

### Python (Data Branch)
- Requires Python 3.10+.
- Recommended to use a virtual environment.
- `pip install -r requirements.txt` (Available on the `data` branch).

## Code of Conduct
Please note that this project is released with a [Contributor Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project you agree to abide by its terms.
