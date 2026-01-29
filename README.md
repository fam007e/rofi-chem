# Rofi-Chem: Data Factory

This branch (`data`) is the **Data Engineering Hub** for Rofi-Chem. It contains the Python pipeline used to fetch, process, and enrich the chemical database from scientific sources like PubChem and Mendeleev.

> [!IMPORTANT]
> This branch is intended for **contributors and data maintainers**. If you just want to use the Rofi plugin, switch to the `main` branch.

## Purpose

The "Data Factory" allows us to:
- Expand the chemical database beyond the 118 periodic elements by adding common compounds.
- Fetch molecular properties automatically via the PubChem API.
- **New**: Retrieves experimental properties like Density, Melting Point, and Boiling Point.
- Maintain a local SQLite cache to avoid repeated network requests.
- Generate the final `chemdata.db` that gets embedded into the Go binary on the `main` branch.

## Environment Setup

### Prerequisites
- Python 3.10+
- Go 1.22+ (for final binary verification)

### Installation
1. Install Python dependencies:
   ```bash
   pip install -r requirements.txt
   ```
2. (Optional) Initialize a virtual environment:
   ```bash
   python -m venv venv
   source venv/bin/activate
   pip install -r requirements.txt
   ```

## How to Enrich Data

1. **Update Molecule List**: Add common chemical names to [data/compounds.txt](data/compounds.txt).
2. **Run the Pipeline**:
   ```bash
   python scripts/init_database.py
   ```
   This script will:
   - Read the names from `compounds.txt`.
   - Check if the compound already exists in `data/chemdata.db`.
   - Fetch missing data from PubChem.
   - Cache results in `data/cache/`.
3. **Verify**: Check the file size or contents of `data/chemdata.db`.

## Syncing with Main

Once you have enriched the database:
1. Copy the updated `data/chemdata.db` to the `internal/db/data/` directory.
2. Commit the database changes.
3. Switch to the `main` branch.
4. Merge or manually update the database file in `main`.
5. Rebuild the Go binary to include the new data:
   ```bash
   go build -o rofi-chem cmd/rofi-chem/main.go
   ```

## Data Architecture & Sources

The enrichment pipeline utilizes specialized scientific libraries to interface with upstream databases:

- **[Mendeleev Python Library](https://github.com/lmmentel/mendeleev)**: Used to pull deep atomic data. We cite the package as a primary source for standardized elemental properties.
- **[PubChem PUG REST](https://pubchem.ncbi.nlm.nih.gov/docs/pug-rest)**: Used to fetch molecular weights, IUPAC names, and formulas for compounds directly from the NIH PubChem database.
- **NIST Standard Reference Data**: Elements are cross-referenced with NIST data where available to ensure scientific accuracy.

All data fetched is stored locally in `data/chemdata.db` (SQLite) to minimize API overhead and ensure project stability.

## License

This project is licensed under the [MIT License](LICENSE).
