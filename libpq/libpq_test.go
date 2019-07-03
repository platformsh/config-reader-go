package libpq_test

import (
	psh "github.com/platformsh/config-reader-go/v2"
	helper "github.com/platformsh/config-reader-go/v2/testdata"
	libpq "github.com/platformsh/config-reader-go/v2/libpq"
	"testing"
)

func TestLibPQFormatterCalled(t *testing.T){
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	credentials, err := config.Credentials("postgresql")
	helper.Ok(t, err)

	formatted, err := libpq.FormattedCredentials(credentials)
	helper.Ok(t, err)

	helper.Equals(t, "host=postgresql.internal port=5432 user=main password=main dbname=main sslmode=disable", formatted)
}
