import json
import os
import requests
from mendeleev import element
import logging

# Set up logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class DataFetcher:
    def __init__(self, cache_dir="data/cache"):
        self.cache_dir = cache_dir
        if not os.path.exists(self.cache_dir):
            os.makedirs(self.cache_dir)

    def fetch_element_data(self, symbol):
        """Fetch element data using mendeleev library."""
        try:
            e = element(symbol)
            return {
                "symbol": e.symbol,
                "name": e.name,
                "atomic_number": e.atomic_number,
                "atomic_mass": e.atomic_weight,
                "density": e.density,
                "melting_point": e.melting_point,
                "boiling_point": e.boiling_point,
                "electronegativity": e.electronegativity(),
                "electron_configuration": str(e.ec),
                "oxidation_states": str(e.oxidation_states),
                "atomic_radius": e.atomic_radius,
                "discovery_year": e.discovery_year,
                "group_number": e.group_id,
                "period": e.period,
                "category": e.series
            }
        except Exception as ex:
            logger.error(f"Error fetching element data for {symbol}: {ex}")
            return None

    def fetch_compound_data(self, query):
        """Fetch compound data from PubChem API with local caching."""
        cache_file = os.path.join(self.cache_dir, f"{query}.json")

        if os.path.exists(cache_file):
            with open(cache_file, 'r') as f:
                return json.load(f)

        try:
            # First, search for the CID
            search_url = f"https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/name/{query}/cids/JSON"
            response = requests.get(search_url, timeout=10)
            response.raise_for_status()
            cid_data = response.json()

            if "IdentifierList" not in cid_data:
                return None

            cid = cid_data["IdentifierList"]["CID"][0]

            # Fetch detailed properties
            props_url = f"https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/cid/{cid}/property/IUPACName,MolecularFormula,MolecularWeight,Complexity/JSON"
            props_response = requests.get(props_url, timeout=10)
            props_response.raise_for_status()
            props_data = props_response.json()["PropertyTable"]["Properties"][0]

            # Fetch experimental properties (like melting point, boiling point)
            # This is more complex as it's in the full record
            data = {
                "name": query,
                "formula": props_data.get("MolecularFormula"),
                "molecular_weight": props_data.get("MolecularWeight"),
                "iupac_name": props_data.get("IUPACName"),
                "pubchem_cid": cid,
                # Placeholders for fields requiring full record parsing
                "density": None,
                "melting_point": None,
                "boiling_point": None,
                "solubility": None,
                "appearance": None,
                "cas_number": None
            }

            # Save to cache
            with open(cache_file, 'w') as f:
                json.dump(data, f)

            return data
        except Exception as ex:
            logger.error(f"Error fetching compound data for {query}: {ex}")
            return None

    def get_all_elements(self):
        """Fetch all 118 elements."""
        from mendeleev import get_all_elements
        elements = []
        for e in get_all_elements():
            elements.append({
                "symbol": e.symbol,
                "name": e.name,
                "atomic_number": e.atomic_number,
                "atomic_mass": e.atomic_weight,
                "density": e.density,
                "melting_point": e.melting_point,
                "boiling_point": e.boiling_point,
                "electronegativity": e.electronegativity(),
                "electron_configuration": str(e.ec),
                "oxidation_states": str(e.oxidation_states),
                "atomic_radius": e.atomic_radius,
                "discovery_year": e.discovery_year,
                "group_number": e.group_id,
                "period": e.period,
                "category": e.series
            })
        return elements
