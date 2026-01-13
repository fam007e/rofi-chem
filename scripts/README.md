# Rofi-Chem Data Enrichment Logic

This directory contains the Python scripts and utilities for the **`data` branch** of Rofi-Chem.

## Purpose
The logic here is responsible for:
1.  Fetching chemical element data from the `mendeleev` scientific library.
2.  Fetching compound and molecule data from the **PubChem REST API**.
3.  Generating and enriching the SQLite `chemdata.db` file.

## Key Files
-   **`init_database.py`**: The main entry point for database population. It reads from `data/compounds.txt` and fetches missing data.
-   **`../data/compounds.txt`**: A list of chemical names to be included in the database.
-   **`../src/`**: (Only on the `data` branch) Contains the source code for the Python fetcher and database manager.

## Usage (On Data Branch)
1.  Install dependencies: `pip install requests mendeleev`
2.  Add compounds to `data/compounds.txt`.
3.  Run: `python scripts/init_database.py`
4.  Copy the updated `data/chemdata.db` to the `main` branch to ship it with the Go binary.
