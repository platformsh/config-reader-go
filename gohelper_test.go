package platformconfig_test

import (
	psh "github.com/platformsh/config-reader-go/v2"
	helper "github.com/platformsh/config-reader-go/v2/testdata"
	"testing"
)

func TestNotOnPlatformReturnsError(t *testing.T) {

	_, err := psh.NewBuildConfigReal(helper.NonPlatformEnv(), "PLATFORM_")

	if err == nil {
		t.Fail()
	}
}

func TestBuildConfigInRuntimeReturnsSuccessfully(t *testing.T) {

	_, err := psh.NewBuildConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)
}

func TestRuntimeConfigInBuildReturnsError(t *testing.T) {

	_, err := psh.NewRuntimeConfigReal(helper.BuildEnv(psh.EnvList{}), "PLATFORM_")

	if err == nil {
		t.Fail()
	}
}

func TestOnDedicatedReturnsTrueOnDedicated(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{"PLATFORM_MODE": "enterprise"}), "PLATFORM_")
	helper.Ok(t, err)

	if !config.OnDedicated() {
		t.Fail()
	}
}

func TestOnDedicatedReturnsFalseOnStandard(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	if config.OnDedicated() {
		t.Fail()
	}
}

func TestOnProductionOnDedicatedProdReturnsTrue(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{
		"PLATFORM_MODE":   "enterprise",
		"PLATFORM_BRANCH": "production",
	}), "PLATFORM_")
	helper.Ok(t, err)

	helper.Assert(t, config.OnProduction(), "OnProduction() returned false when it should be true.")
}

func TestOnProductionOnDedicatedStagingReturnsFalse(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{
		"PLATFORM_MODE":   "enterprise",
		"PLATFORM_BRANCH": "staging",
	}), "PLATFORM_")
	helper.Ok(t, err)

	helper.Assert(t, !config.OnProduction(), "OnProduction() returned true when it should be false.")
}

func TestOnProductionOnStandardProdReturnsTrue(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{
		"PLATFORM_BRANCH": "master",
	}), "PLATFORM_")
	helper.Ok(t, err)

	helper.Assert(t, config.OnProduction(), "OnProduction() returned false when it should be true.")
}

func TestOnProductionOnStandardStagingReturnsFalse(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	helper.Assert(t, !config.OnProduction(), "OnProduction() returned true when it should be false.")
}

func TestBuildPropertyInBuildExists(t *testing.T) {
	config, err := psh.NewBuildConfigReal(helper.BuildEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	helper.Equals(t, "/app", config.AppDir())
	helper.Equals(t, "app", config.ApplicationName())
	helper.Equals(t, "test-project", config.Project())
	helper.Equals(t, "abc123", config.TreeId())
	helper.Equals(t, "def789", config.ProjectEntropy())
}

func TestBuildAndRuntimePropertyInRuntimeExists(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	helper.Equals(t, "/app", config.AppDir())
	helper.Equals(t, "app", config.ApplicationName())
	helper.Equals(t, "test-project", config.Project())
	helper.Equals(t, "abc123", config.TreeId())
	helper.Equals(t, "def789", config.ProjectEntropy())

	helper.Equals(t, "feature-x", config.Branch())
	helper.Equals(t, "feature-x-hgi456", config.Environment())
	helper.Equals(t, "/app/web", config.DocumentRoot())
	helper.Equals(t, "1.2.3.4", config.SmtpHost())
	helper.Equals(t, "8080", config.Port())
	helper.Equals(t, "unix://tmp/blah.sock", config.Socket())
}

func TestReadingExistingVariableWorks(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	helper.Equals(t, "someval", config.Variable("somevar", ""))
}

func TestReadingMissingVariableReturnsDefault(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	helper.Equals(t, "default-val", config.Variable("missing", "default-val"))
}

func TestVariablesReturnsMapWithData(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	list := config.Variables()

	helper.Equals(t, "someval", list["somevar"])
}

func TestCredentialsForExistingRelationshipReturns(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	creds, err := config.Credentials("database")
	helper.Ok(t, err)

	helper.Equals(t, "mysql", creds.Scheme)
}

//public function test_credentials_missing_relationship_throws() : void
func TestCredentialsForMissingRelationshipErrrors(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	_, err = config.Credentials("does-not-exist")

	if err == nil {
		t.Fail()
	}
}

func TestGetAllRoutesAtRuntimeWorks(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	routes := config.Routes()

	helper.Equals(t, "upstream", routes["https://www.master-7rqtwti-gcpjkefjk4wc2.us-2.platformsh.site/"].Type)
}

func TestGetPrimaryRouteWorks(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	route, ok := config.PrimaryRoute()

	helper.Equals(t, true, ok)
	helper.Equals(t, true, route.Primary)
	helper.Equals(t, "main", route.Id)
}

func TestGetUpstreamRoutesWorks(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	routes := config.UpstreamRoutes()

	helper.Equals(t, 3, len(routes))
	helper.Equals(t, "https://www.{default}/", routes["https://www.master-7rqtwti-gcpjkefjk4wc2.us-2.platformsh.site/"].OriginalUrl)
}

func TestGetUpstreamRoutesForAppWorks(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	routes := config.UpstreamRoutesForApp("app")

	helper.Equals(t, 2, len(routes))
	helper.Equals(t, "https://www.{default}/", routes["https://www.master-7rqtwti-gcpjkefjk4wc2.us-2.platformsh.site/"].OriginalUrl)
}

func TestGetRouteByIdWorks(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	route, ok := config.Route("main")

	helper.Equals(t, true, ok)
	helper.Equals(t, "upstream", route.Type)
}

func TestRouteUrlIsAddedProperly(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	route, ok := config.Route("main")

	helper.Equals(t, true, ok)
	helper.Equals(t, "https://www.master-7rqtwti-gcpjkefjk4wc2.us-2.platformsh.site/", route.Url)
}

func TestGetNonExistentRouteErrors(t *testing.T) {
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{}), "PLATFORM_")
	helper.Ok(t, err)

	_, ok := config.Route("missing")

	helper.Equals(t, false, ok)
}

func TestLocalNoRelationships(t *testing.T) {
	// This should not fail due to the missing relationships.
	config, err := psh.NewRuntimeConfigReal(helper.RuntimeEnv(psh.EnvList{
		"PLATFORM_RELATIONSHIPS": "",
		"PLATFORM_VARIABLES":     "",
		"PLATFORM_APPLICATION":   "",
		"PLATFORM_ROUTES":        "",
	}), "PLATFORM_")
	helper.Ok(t, err)

	// But this should error out.
	_, err = config.Credentials("does-not-exist")

	if err == nil {
		t.Fail()
	}
}
