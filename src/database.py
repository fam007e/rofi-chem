import sqlite3
import os
import logging

logger = logging.getLogger(__name__)

class DatabaseManager:
    def __init__(self, db_path="data/chemdata.db"):
        self.db_path = db_path
        self._ensure_db_dir()

    def _ensure_db_dir(self):
        db_dir = os.path.dirname(self.db_path)
        if db_dir and not os.path.exists(db_dir):
            os.makedirs(db_dir)

    def get_connection(self):
        return sqlite3.connect(self.db_path)

    def init_database(self):
        """Initialize the database schema."""
        with self.get_connection() as conn:
            cursor = conn.cursor()

            # Elements Table
            cursor.execute('''
                CREATE TABLE IF NOT EXISTS elements (
                    id INTEGER PRIMARY KEY,
                    symbol TEXT UNIQUE NOT NULL,
                    name TEXT NOT NULL,
                    atomic_number INTEGER UNIQUE NOT NULL,
                    atomic_mass REAL,
                    density REAL,
                    melting_point REAL,
                    boiling_point REAL,
                    electronegativity REAL,
                    electron_configuration TEXT,
                    oxidation_states TEXT,
                    atomic_radius REAL,
                    discovery_year INTEGER,
                    group_number INTEGER,
                    period INTEGER,
                    category TEXT
                )
            ''')

            # Compounds Table
            cursor.execute('''
                CREATE TABLE IF NOT EXISTS compounds (
                    id INTEGER PRIMARY KEY,
                    name TEXT NOT NULL,
                    formula TEXT NOT NULL,
                    molecular_weight REAL,
                    density REAL,
                    melting_point REAL,
                    boiling_point REAL,
                    solubility TEXT,
                    appearance TEXT,
                    iupac_name TEXT,
                    cas_number TEXT,
                    pubchem_cid INTEGER
                )
            ''')
            conn.commit()

    def insert_element(self, element_data):
        """Insert or replace an element in the database."""
        with self.get_connection() as conn:
            cursor = conn.cursor()
            columns = ', '.join(element_data.keys())
            placeholders = ', '.join(['?'] * len(element_data))
            sql = f"INSERT OR REPLACE INTO elements ({columns}) VALUES ({placeholders})"
            cursor.execute(sql, list(element_data.values()))
            conn.commit()

    def insert_compound(self, compound_data):
        """Insert or replace a compound in the database."""
        with self.get_connection() as conn:
            cursor = conn.cursor()
            columns = ', '.join(compound_data.keys())
            placeholders = ', '.join(['?'] * len(compound_data))
            sql = f"INSERT OR REPLACE INTO compounds ({columns}) VALUES ({placeholders})"
            cursor.execute(sql, list(compound_data.values()))
            conn.commit()

    def search_elements(self, query):
        """Search elements by name or symbol."""
        with self.get_connection() as conn:
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            sql = "SELECT * FROM elements WHERE name LIKE ? OR symbol LIKE ? OR atomic_number LIKE ?"
            search_query = f"%{query}%"
            cursor.execute(sql, (search_query, search_query, search_query))
            return [dict(row) for row in cursor.fetchall()]

    def search_compounds(self, query):
        """Search compounds by name or formula."""
        with self.get_connection() as conn:
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            sql = "SELECT * FROM compounds WHERE name LIKE ? OR formula LIKE ?"
            search_query = f"%{query}%"
            cursor.execute(sql, (search_query, search_query))
            return [dict(row) for row in cursor.fetchall()]

    def get_all_elements(self):
        """Get all elements."""
        with self.get_connection() as conn:
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM elements ORDER BY atomic_number")
            return [dict(row) for row in cursor.fetchall()]

    def get_compound_by_name(self, name):
        """Get a compound by its exact name."""
        with self.get_connection() as conn:
            conn.row_factory = sqlite3.Row
            cursor = conn.cursor()
            cursor.execute("SELECT * FROM compounds WHERE name = ?", (name,))
            row = cursor.fetchone()
            return dict(row) if row else None
