package gosolr

import (
  "fmt"
  psh "github.com/platformsh/config-reader-go"
)

// Go-solr requires a string that includes the full collection path to  connect to Solr.
func FormattedCredentials(creds psh.Credential) (string, error) {

  formatted := fmt.Sprintf("http://%s:%d/%s", creds.Host, creds.Port, creds.Path)

  return formatted, nil

}
