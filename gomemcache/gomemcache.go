package gomemcache

import (
  "fmt"
  psh "github.com/platformsh/config-reader-go"
)

// The gomemcache library requires a specific string to connect to Memcached.
func FormattedCredentials(creds psh.Credential) (string, error) {
  formatted := fmt.Sprintf("%s:%d", creds.Host, creds.Port)
  return formatted, nil
}
