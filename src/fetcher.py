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
                "oxidation_states": str(e.oxistates),
                "atomic_radius": e.atomic_radius,
                "discovery_year": e.discovery_year,
                "group_number": e.group_id,
                "period": e.period,
                "category": e.series
            }
        except Exception as ex:
            logger.error(f"Error fetching element data for {symbol}: {ex}")
            return None

    def fetch_experimental_properties(self, cid):
        """Fetch experimental properties from PubChem View API."""
        try:
            url = f"https://pubchem.ncbi.nlm.nih.gov/rest/pug_view/data/compound/{cid}/JSON?heading=Experimental+Properties"
            response = requests.get(url, timeout=10)
            if response.status_code != 200:
                return {}

            data = response.json()
            properties = {}

            # Navigate to Experimental Properties section
            sections = data.get("Record", {}).get("Section", [])
            for sec in sections:
                if sec.get("TOCHeading") == "Chemical and Physical Properties":
                    for sub_sec in sec.get("Section", []):
                        if sub_sec.get("TOCHeading") == "Experimental Properties":
                            # iterate through properties
                            for prop in sub_sec.get("Section", []):
                                heading = prop.get("TOCHeading")
                                info = prop.get("Information", [])
                                if not info:
                                    continue

                                # Extract first valid string value
                                value = None
                                for item in info:
                                    if "Value" in item:
                                        val_obj = item["Value"]
                                        if "StringWithMarkup" in val_obj:
                                            value = val_obj["StringWithMarkup"][0]["String"]
                                            break
                                        elif "Number" in val_obj and "Unit" in val_obj:
                                            value = f"{val_obj['Number'][0]} {val_obj['Unit']}"
                                            break

                                if value:
                                    if heading == "Boiling Point":
                                        properties["boiling_point"] = value
                                    elif heading == "Melting Point":
                                        properties["melting_point"] = value
                                    elif heading == "Density":
                                        properties["density"] = value
                                    elif heading == "Solubility":
                                        properties["solubility"] = value
                                    elif heading == "Physical Description" or heading == "Color/Form":
                                        properties["appearance"] = value

            return properties
        except Exception as ex:
            logger.error(f"Error fetching experimental properties for CID {cid}: {ex}")
            return {}

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
            if response.status_code != 200:
                logger.warning(f"Could not find CID for {query}")
                return None

            cid_data = response.json()

            if "IdentifierList" not in cid_data:
                return None

            cid = cid_data["IdentifierList"]["CID"][0]

            # Fetch basic properties
            props_url = f"https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/cid/{cid}/property/IUPACName,MolecularFormula,MolecularWeight,Complexity/JSON"
            props_response = requests.get(props_url, timeout=10)
            props_response.raise_for_status()
            props_data = props_response.json()["PropertyTable"]["Properties"][0]

            # Fetch experimental properties
            exp_props = self.fetch_experimental_properties(cid)

            # Combine data (parsing numeric values where possible could be added here)
            # For now, we store the raw strings from PubChem as the database schema
            # might expect REAL for some, but PubChem returns strings like "78 C".
            # To keep it simple and safe, we'll try to extract the first number for REAL fields
            # or just store the string if we change the DB schema (which is currently REAL).

            # Simple helpers to extract float from string
            def extract_float(s):
                if not s: return None
                import re
                # simplistic: take first float found
                match = re.search(r"[-+]?\d*\.\d+|\d+", s)
                return float(match.group()) if match else None

            data = {
                "name": query,
                "formula": props_data.get("MolecularFormula"),
                "molecular_weight": props_data.get("MolecularWeight"),
                "iupac_name": props_data.get("IUPACName"),
                "pubchem_cid": cid,
                "density": extract_float(exp_props.get("density")),
                "melting_point": extract_float(exp_props.get("melting_point")),
                "boiling_point": extract_float(exp_props.get("boiling_point")),
                "solubility": exp_props.get("solubility"),
                "appearance": exp_props.get("appearance"),
                "cas_number": None # CAS is harder to get reliably without another call, skipping for now
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
                "oxidation_states": str(e.oxistates),
                "atomic_radius": e.atomic_radius,
                "discovery_year": e.discovery_year,
                "group_number": e.group_id,
                "period": e.period,
                "category": e.series
            })
        return elements
