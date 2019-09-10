package sqldsn_test

import (
	psh "github.com/platformsh/config-reader-go"
	helper "github.com/platformsh/config-reader-go/testdata"
	sqldsn "github.com/platformsh/config-reader-go/sqldsn"
	"testing"
)

func TestSqlDsnFormatterCalled(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	credentials, err := config.Credentials("database")
	helper.Ok(t, err)

	formatted, err := sqldsn.FormattedCredentials(credentials)
	helper.Ok(t, err)

	helper.Equals(t, "user:@tcp(database.internal:3306)/main?charset=utf8", formatted)
}
