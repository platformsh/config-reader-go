package gomemcache_test

import (
	psh "github.com/platformsh/config-reader-go"
	helper "github.com/platformsh/config-reader-go/testdata"
	mem "github.com/platformsh/config-reader-go/gomemcache"
	"testing"
)

func TestGoMemcacheFormatterCalled(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	credentials, err := config.Credentials("memcached")
	helper.Ok(t, err)

	formatted, err := mem.FormattedCredentials(credentials)
	helper.Ok(t, err)

	helper.Equals(t, "memcached.internal:11211", formatted)
}
