package amqp

import (
  "fmt"
  psh "github.com/platformsh/config-reader-go/v2"
)

// AMQP requires a specfic connection string to connect to PostgreSQL.
func FormattedCredentials(creds psh.Credential) (string, error) {
  formatted := fmt.Sprintf("%s://%s:%s@%s:%d/", creds.Scheme, creds.Username,
    creds.Password, creds.Host, creds.Port)
  return formatted, nil
}
