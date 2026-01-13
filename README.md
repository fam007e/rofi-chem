# Rofi-Chem - Chemical Elements & Compounds Rofi Plugin

A Rofi mode plugin for quick lookup of periodic table elements and chemical compounds with data from NIST, Mendeleev, and PubChem.

## Features

- Interactive periodic table lookup integrated with Rofi.
- Instant access to element properties (atomic number, mass, density, etc.).
- Chemical compound information (formula, properties, molecular weight).
- Configurable display fields.
- Fuzzy search by name, symbol, or formula.
- Local SQLite cache for fast offline access.

## Installation

### Option 1: Direct Binary (Recommended)
Download the latest `rofi-chem` binary from the [Releases](https://github.com/fam007e/rofi-chem/releases) page and place it in your `$PATH` (e.g., `~/.local/bin/`).

### Option 2: Build from Source
```bash
git clone https://github.com/fam007e/rofi-chem.git
cd rofi-chem
go build -o rofi-chem cmd/rofi-chem/main.go
mv rofi-chem ~/.local/bin/
```

> [!NOTE]
> No database setup is required. The element data is embedded inside the binary and will be automatically extracted on first run.

## Data & Contributions

This branch (`main`) focuses on the Rofi interface and Go core.

- **To add new compounds**: Please switch to the `data` branch, which contains the specialized Python pipeline for data enrichment.
- **To contribute**: See [CONTRIBUTING.md](CONTRIBUTING.md) for details on our dual-branch strategy.

## Usage

```bash
rofi -modi "chem:rofi-chem" -show chem
```

## Configuration

Configuration is stored in `~/.config/rofi/rofi-chem/config.yaml`.

## License

MIT
