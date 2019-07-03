package libpq

import (
  "fmt"
  psh "github.com/platformsh/config-reader-go/v2"
)

// lib/pq requires a specfic connection string to connect to PostgreSQL.
func FormattedCredentials(creds psh.Credential) (string, error) {
  formatted := fmt.Sprintf("host=%s port=%d user=%s " + "password=%s dbname=%s sslmode=disable",
    creds.Host, creds.Port, creds.Username, creds.Password, creds.Path)
  return formatted, nil
}
