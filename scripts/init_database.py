import sys
import os

# Add src to path
sys.path.append(os.path.join(os.getcwd(), 'src'))

from database import DatabaseManager
from fetcher import DataFetcher
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def init():
    db = DatabaseManager()
    fetcher = DataFetcher()

    logger.info("Initializing database...")
    db.init_database()

    logger.info("Fetching and inserting elements...")
    elements = fetcher.get_all_elements()
    for e in elements:
        db.insert_element(e)

    logger.info(f"Inserted {len(elements)} elements.")

    # Read compounds from external file
    compounds_file = os.path.join(os.getcwd(), 'data', 'compounds.txt')
    if not os.path.exists(compounds_file):
        logger.warning(f"Compounds file not found: {compounds_file}")
        return

    with open(compounds_file, 'r') as f:
        compounds_names = [line.strip() for line in f if line.strip() and not line.strip().startswith('#')]

    logger.info(f"Found {len(compounds_names)} compound names in {compounds_file}")

    count = 0
    skipped = 0
    for name in compounds_names:
        # Check if already in DB
        if db.get_compound_by_name(name):
            skipped += 1
            continue

        logger.info(f"Fetching data for: {name}")
        compound_data = fetcher.fetch_compound_data(name)
        if compound_data:
            db.insert_compound(compound_data)
            count += 1
        else:
            logger.warning(f"Could not fetch data for: {name}")

    logger.info(f"Inserted {count} new compounds. Skipped {skipped} already in database.")
    logger.info("Database initialization complete.")

if __name__ == "__main__":
    init()
