package gosolr_test

import (
	psh "github.com/platformsh/config-reader-go"
	helper "github.com/platformsh/config-reader-go/testdata"
	gosolr "github.com/platformsh/config-reader-go/gosolr"
	"testing"
)

func TestGoSolrFormatterCalled(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	credentials, err := config.Credentials("solr")
	helper.Ok(t, err)

	formatted, err := gosolr.FormattedCredentials(credentials)
	helper.Ok(t, err)

  helper.Equals(t, "http://solr.internal:8080/solr/collection1", formatted)
}
