package cache

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

const CacheTTL = 24 * time.Hour

// InitializeCache sets up the SQLite database and creates the cache table if it doesn't exist.
func InitializeCache() error {
	// Connect to SQLite database (creates cache.db file if it doesn't exist)
	var err error
	db, err = sql.Open("sqlite3", "./cache.db")
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite: %v", err)
	}

	// Create the cache table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS cache (
		url TEXT PRIMARY KEY,
		response TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create cache table: %v", err)
	}

	return nil
}

// CloseCache closes the SQLite database connection.
func CloseCache() {
	if db != nil {
		db.Close()
	}
}

// isExpired checks if a given timestamp is older than the TTL duration.
func isExpired(timestamp string) (bool, error) {
	parsedTimestamp, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return false, fmt.Errorf("failed to parse timestamp: %v", err)
	}
	return time.Since(parsedTimestamp) > CacheTTL, nil
}

// getCachedResponse retrieves a cached response and checks if it is expired.
func getCachedResponse(url string) (string, bool, error) {
	var response string
	var timestamp string
	err := db.QueryRow("SELECT response, timestamp FROM cache WHERE url = ?", url).Scan(&response, &timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			// No cache entry for this URL
			return "", false, nil
		}
		return "", false, fmt.Errorf("error checking cache for URL %s: %v", url, err)
	}

	// Check if the cache entry is expired
	expired, err := isExpired(timestamp)
	if err != nil {
		return "", false, err
	}
	if expired {
		// Cache entry is expired
		return "", false, nil
	}

	// Cache entry is valid
	return response, true, nil
}

// cacheResponse stores a new response for a given URL in the cache, updating the timestamp.
func cacheResponse(url, response string) error {
	_, err := db.Exec(
		"INSERT OR REPLACE INTO cache (url, response, timestamp) VALUES (?, ?, ?)",
		url, response, time.Now().Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("error caching response for URL %s: %v", url, err)
	}
	return nil
}
