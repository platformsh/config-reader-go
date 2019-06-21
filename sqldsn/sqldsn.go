package sqldsn

import (
  "fmt"
  psh "github.com/platformsh/config-reader-go/v2"
)

// SqlDsn produces an SQL connection string appropriate for use with many
// common Go database tools.
func FormattedCredentials(creds psh.Credential) (string, error) {

  formatted := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", creds.Username, creds.Password, creds.Host, creds.Port, creds.Path)
	return formatted, nil

}
