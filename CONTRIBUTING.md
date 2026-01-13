# Contributing To Rofi-Chem

Thank you for helping build a better tool for the scientific community. To keep the project maintainable and performant, we use a strictly separated dual-branch architecture.

## 🏗 Project Architecture

| Branch | Purpose | Primary Tech | Target Audience |
| :--- | :--- | :--- | :--- |
| **`main`** | Core Rofi Plugin & Distribution | Go | End Users / UI Devs |
| **`data`** | Data Enrichment Factory | Python / SQLite | Data Scientists / Maintainers |

### Which branch should I use?
- **Bug in Rofi display or lookup logic?** Use `main`.
- **Typo in chemical mass or missing compounds?** Use `data`.
- **Requesting a new feature?** Open an Issue first.

---

## 🚀 Development Workflow

### 1. Code Contributions (`main`)
The `main` branch is optimized for a zero-dependency user experience.
- **Requirements**: Go 1.22+
- **Build**: `go build -o rofi-chem cmd/rofi-chem/main.go`
- **Embedding**: Any data changes must be merged from `data` and then manually updated in `internal/db/data/chemdata.db` before rebuilding.

### 2. Data Contributions (`data`)
The `data` branch contains the Python pipeline for scientific data fetching.
- **Requirements**: Python 3.10+, `mendeleev`, `pubchempy`
- **Setup**: `pip install -r requirements.txt`
- **Workflow**:
  1. Add compound names to `data/compounds.txt`.
  2. Run `python scripts/init_database.py` to fetch properties and update the local database.

---

## 🛠 Submission Guidelines

1. **Self-Documentation**: Ensure any UI changes are reflected in the `main` branch README, and any pipeline changes in the `data` branch README.
2. **Commit Hygiene**: Use descriptive commit messages. GPG signing is preferred but not required for community PRs.
3. **Draft PRs**: Feel free to open a Draft PR if you want early feedback on a complex data pipeline or UI refactor.

## ⚖ Code of Conduct
Respect for scientific accuracy and community members is paramount. See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
