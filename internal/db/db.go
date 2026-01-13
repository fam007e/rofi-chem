package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed data/chemdata.db
var embeddedDB embed.FS


type Database struct {
	conn *sql.DB
}

type Element struct {
	Symbol        string
	Name          string
	AtomicNumber  int
	AtomicMass    float64
	Density       float64
	MeltingPoint  float64
	BoilingPoint  float64
	ElectronConfig string
    // Dynamic access to other fields is harder in Go,
    // so we'll map common ones or use a map[string]interface{} if needed later.
    // For now, let's keep it simple and extensible.
    RawData map[string]interface{}
}

type Compound struct {
	Name            string
	Formula         string
	MolecularWeight float64
    RawData map[string]interface{}
}

func NewDatabase() (*Database, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

    // Try default path first
	dbPath := filepath.Join(home, ".config", "rofi", "rofi-chem", "data", "chemdata.db")

    // Fallback to local data/chemdata.db (dev mode)
    if _, err := os.Stat(dbPath); os.IsNotExist(err) {
        cwd, _ := os.Getwd()
        localPath := filepath.Join(cwd, "data", "chemdata.db")
        if _, err := os.Stat(localPath); err == nil {
            dbPath = localPath
        } else {
            // Last resort: extract embedded database to cache
            cacheDir, _ := os.UserCacheDir()
            dbPath = filepath.Join(cacheDir, "rofi-chem", "chemdata.db")
            if _, err := os.Stat(dbPath); os.IsNotExist(err) {
                if err := extractEmbeddedDB(dbPath); err != nil {
                    return nil, fmt.Errorf("failed to extract embedded database: %w", err)
                }
            }
        }
    }


	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &Database{conn: db}, nil
}

func extractEmbeddedDB(dest string) error {
	err := os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	src, err := embeddedDB.Open("data/chemdata.db")
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}


func (d *Database) Close() {
	d.conn.Close()
}

func (d *Database) GetAllElements() ([]map[string]interface{}, error) {
	rows, err := d.conn.Query("SELECT * FROM elements ORDER BY atomic_number")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToMap(rows)
}

func (d *Database) SearchCompounds(query string) ([]map[string]interface{}, error) {
	q := "%" + query + "%"
	rows, err := d.conn.Query("SELECT * FROM compounds WHERE name LIKE ? OR formula LIKE ?", q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToMap(rows)
}

func rowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice of interface{} to hold values
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		results = append(results, entry)
	}
	return results, nil
}

func (d *Database) GetElementBySymbol(symbol string) (map[string]interface{}, error) {
	rows, err := d.conn.Query("SELECT * FROM elements WHERE symbol = ?", symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := rowsToMap(rows)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("element not found")
	}
	return results[0], nil
}

func (d *Database) GetCompoundByName(name string) (map[string]interface{}, error) {
	rows, err := d.conn.Query("SELECT * FROM compounds WHERE name = ?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := rowsToMap(rows)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("compound not found")
	}
	return results[0], nil
}
