package mongo

import (
  "fmt"
  psh "github.com/platformsh/config-reader-go"
)

// The mongo-go-driver requires a specific string to connect to MongoDB.
func FormattedCredentials(creds psh.Credential) (string, error) {
  formatted := fmt.Sprintf("%s://%s:%s@%s:%d/%s",
    creds.Scheme, creds.Username, creds.Password, creds.Host, creds.Port, creds.Path)
  return formatted, nil
}
