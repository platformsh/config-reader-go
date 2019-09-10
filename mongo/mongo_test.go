package mongo_test

import (
	psh "github.com/platformsh/config-reader-go"
	helper "github.com/platformsh/config-reader-go/testdata"
	libpq "github.com/platformsh/config-reader-go/mongo"
	"testing"
)

func TestMongoDriverFormatterCalled(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	credentials, err := config.Credentials("mongodb")
	helper.Ok(t, err)

	formatted, err := libpq.FormattedCredentials(credentials)
	helper.Ok(t, err)

	helper.Equals(t, "mongodb://main:main@mongodb.internal:27017/main", formatted)
}
