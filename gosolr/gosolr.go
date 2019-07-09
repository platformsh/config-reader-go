package gosolr

import (
  "fmt"
  "strings"
  psh "github.com/platformsh/config-reader-go/v2"
)

type SolrCredentials struct {
  Url        string
  Collection string
}

// go-solr requires two separate strings to connect to Solr, a url and a collection name.
func FormattedCredentials(creds psh.Credential) (SolrCredentials, error) {

  var formatted SolrCredentials

  path := strings.SplitAfter(creds.Path, "/")

  formatted.Url = fmt.Sprintf("http://%s:%d/%s", creds.Host, creds.Port, path[0])
  formatted.Collection = path[1]

  return formatted, nil

}
