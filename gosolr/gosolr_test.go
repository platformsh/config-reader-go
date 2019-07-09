package gosolr_test

import (
	psh "github.com/platformsh/config-reader-go/v2"
	helper "github.com/platformsh/config-reader-go/v2/testdata"
	gosolr "github.com/platformsh/config-reader-go/v2/gosolr"
	"testing"
)

func TestGoSolrFormatterCalled(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	credentials, err := config.Credentials("solr")
	helper.Ok(t, err)

	formatted, err := gosolr.FormattedCredentials(credentials)
	helper.Ok(t, err)

  helper.Equals(t, "http://solr.internal:8080/solr/", formatted.Url)
  helper.Equals(t, "collection1", formatted.Collection)
}
