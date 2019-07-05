package amqp_test

import (
	psh "github.com/platformsh/config-reader-go/v2"
	helper "github.com/platformsh/config-reader-go/v2/testdata"
	amqp "github.com/platformsh/config-reader-go/v2/amqp"
	"testing"
)

func TestAMQPFormatterCalled(t *testing.T){
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	credentials, err := config.Credentials("rabbitmq")
	helper.Ok(t, err)

	formatted, err := amqp.FormattedCredentials(credentials)
	helper.Ok(t, err)

  helper.Equals(t, "amqp://guest:guest@rabbitmq.internal:5672/", formatted)
}
