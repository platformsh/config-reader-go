// The gohelper library provides abstractions for the Platform.sh environment
// to make it easier to configure applications to run on Platform.sh.
// See https://docs.platform.sh/development/variables.html for an in-depth
// description of the available properties and their meaning.
package platformconfig

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var NotValidPlatform = errors.New("No valid platform found.")

var NotRuntimePlatform = errors.New("No valid runtime platform found.")

type EnvList map[string]string

type envReader func(string) string

type Credential struct {
	Scheme   string `json:"scheme"`
	Cluster  string `json:"cluster"`
	Service  string `json:"service"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Path     string `json:"path"`
	Public   bool   `json:"public"`
	Fragment string `json:"fragment"`
	Ip       string `json:"ip"`
	Rel      string `json:"rel"`
	Type     string `json:"type"`
	Port     int    `json:"port"`
	Hostname string `json:"hostname"`
	Query    struct {
		IsMaster bool `json:"is_master"`
	}
}

type Credentials map[string][]Credential

type Route struct {
	OriginalUrl    string            `json:"original_url"`
	Attributes     map[string]string `json:"attributes"`
	Type           string            `json:"type"`
	RestrictRobots bool              `json:"restrict_robots"`
	Tls            struct {
		ClientAuthentication         string   `json:"client_authentication"`
		MinVersion                   int      `json:"min_version"`
		ClientCertificateAuthorities []string `json:"client_certificate_authorities"`
		StrictTransportSecurity      struct {
			IncludeSubdomains bool `json:"include_subdomains"`
			Enabled           bool `json:"enabled"`
			Preload           bool `json:"preload"`
		}
	}
	Upstream string `json:"upstream"`
	Cache    struct {
		Enabled    bool     `json:"enabled"`
		Headers    []string `json:"headers"`
		Cookies    []string `json:"cookies"`
		DefaultTtl int      `json:"default_ttl"`
	}
	HttpAccess struct {
		Addresses []string          `json:"addresses"`
		BasicAuth map[string]string `json:"basic_auth"`
	}
	Primary bool   `json:"primary"`
	Id      string `json:"id"`
	Ssi     struct {
		Enabled bool `json:"enabled"`
	}

	// This field is not part of the JSON definition, but it gets added
	// to the struct from the JSON array key.
	Url string
}

type Routes map[string]*Route

type BuildConfig struct {
	// Prefixed simple values, build or deploy.
	applicationName string
	treeId          string
	appDir          string
	project         string
	projectEntropy  string

	// Prefixed complex values.
	variables   EnvList
	application map[string]interface{}

	// Internal data.
	prefix string
}

type RuntimeConfig struct {
	BuildConfig

	// Prefixed simple values, runtime only.
	branch       string
	environment  string
	documentRoot string
	smtpHost     string
	mode         string

	// Prefixed complex values.
	credentials Credentials
	variables   EnvList
	routes      Routes

	// Unprefixed simple values.
	socket string
	port   string
}

func NewBuildConfigReal(getter envReader, prefix string) (*BuildConfig, error) {
	p := &BuildConfig{}

	p.prefix = prefix

	// If it's not a valid platform, bail out now.
	if getter(prefix+"APPLICATION_NAME") == "" {
		return nil, NotValidPlatform
	}

	// Extract the easy environment variables.
	p.applicationName = getter(p.prefix + "APPLICATION_NAME")
	p.appDir = getter(p.prefix + "APP_DIR")
	p.treeId = getter(p.prefix + "TREE_ID")
	p.project = getter(p.prefix + "PROJECT")
	p.projectEntropy = getter(p.prefix + "PROJECT_ENTROPY")

	// Extract the complex environment variables (serialized JSON strings).

	// Extract the PLATFORM_VARIABLES array.
	if vars := getter(p.prefix + "VARIABLES"); vars != "" {
		parsedVars, err := extractVariables(vars)
		if err != nil {
			return nil, err
		}
		p.variables = parsedVars
	}

	// Extract PLATFORM_APPLICATION.
	// @todo Turn this into a proper struct.
	application := getter(p.prefix + "APPLICATION")
	if application != "" {
		var parsedApplication map[string]interface{}
		jsonApplication, err := base64.StdEncoding.DecodeString(application)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(jsonApplication, &parsedApplication)
		if err != nil {
			return nil, err
		}
		p.application = parsedApplication
	}

	return p, nil
}

func NewRuntimeConfigReal(getter envReader, prefix string) (*RuntimeConfig, error) {
	b, err := NewBuildConfigReal(getter, prefix)

	if err != nil {
		return nil, err
	}

	p := &RuntimeConfig{BuildConfig: *b}

	p.prefix = prefix

	// If it's not a valid platform, bail out now.
	if getter(prefix+"BRANCH") == "" {
		return nil, NotRuntimePlatform
	}

	// Extract the easy environment variables.
	p.documentRoot = getter(p.prefix + "DOCUMENT_ROOT")
	p.branch = getter(p.prefix + "BRANCH")
	p.environment = getter(p.prefix + "ENVIRONMENT")
	p.project = getter(p.prefix + "PROJECT")
	p.smtpHost = getter(p.prefix + "SMTP_HOST")
	p.mode = getter(p.prefix + "MODE")
	p.socket = getter("SOCKET")
	p.port = getter("PORT")

	// Extract the complex environment variables (serialized JSON strings).

	// Extract PLATFORM_RELATIONSHIPS, which we'll call credentials since that's what they are.
	if rels := getter(p.prefix + "RELATIONSHIPS"); rels != "" {
		creds, err := extractCredentials(rels)
		if err != nil {
			return nil, err
		}
		p.credentials = creds
	}

	// Extract the PLATFORM_VARIABLES array.
	if vars := getter(p.prefix + "VARIABLES"); vars != "" {
		parsedVars, err := extractVariables(vars)
		if err != nil {
			return nil, err
		}
		p.variables = parsedVars
	}

	// Extract PLATFORM_ROUTES.
	if routes := getter(p.prefix + "ROUTES"); routes != "" {
		parsedRoutes, err := extractRoutes(routes)
		if err != nil {
			return nil, err
		}
		p.routes = parsedRoutes
	}

	return p, nil
}

// This function returns a new Config object, representing
// the abstracted Platform.sh environment.  If run on not a Platform.sh
// environment (eg, a local computer) then it will return nil and an error.
func NewBuildConfig() (*BuildConfig, error) {
	return NewBuildConfigReal(os.Getenv, "PLATFORM_")
}

// This function returns a new Config object, representing
// the abstracted Platform.sh environment.  If run on not a Platform.sh
// environment (eg, a local computer) then it will return nil and an error.
func NewRuntimeConfig() (*RuntimeConfig, error) {
	return NewRuntimeConfigReal(os.Getenv, "PLATFORM_")
}

// Determines if the current environment is a Platform.sh Enterprise environment.
func (p *RuntimeConfig) OnEnterprise() bool {
	return p.mode == "enterprise"
}

// Determines if the current environment is a production environment.
//
// Note: There may be a few edge cases where this is not entirely correct on Enterprise,
// if the production branch is not named `production`.  In that case you'll need to use
// your own logic.
func (p *RuntimeConfig) OnProduction() bool {
	var prodBranch string
	if p.OnEnterprise() {
		prodBranch = "production"
	} else {
		prodBranch = "master"
	}

	return p.branch == prodBranch
}

// The name of the application, as defined in its configuration.
func (p *BuildConfig) ApplicationName() string {
	return p.applicationName
}

// An ID identifying the application tree before it was built: a unique hash
// is generated based on the contents of the application's files in the
// repository.
func (p *BuildConfig) TreeId() string {
	return p.treeId
}

// The absolute path to the application.
func (p *BuildConfig) AppDir() string {
	return p.appDir
}

// The project ID.
func (p *BuildConfig) Project() string {
	return p.project
}

// A random string generated for each project, useful for generating hash keys.
func (p *BuildConfig) ProjectEntropy() string {
	return p.projectEntropy
}

// The Git branch name.
func (p *RuntimeConfig) Branch() string {
	return p.branch
}

// The environment ID (usually the Git branch plus a hash).
func (p *RuntimeConfig) Environment() string {
	return p.environment
}

// The absolute path to the web root of the application.
func (p *RuntimeConfig) DocumentRoot() string {
	return p.documentRoot
}

// The hostname of the Platform.sh default SMTP server (an empty string if
// emails are disabled on the environment).
func (p *RuntimeConfig) SmtpHost() string {
	return p.smtpHost
}

// The TCP port number the application should listen to for incoming requests.
func (p *RuntimeConfig) Port() string {
	return p.port
}

// The Unix socket the application should listen to for incoming requests.
func (p *RuntimeConfig) Socket() string {
	return p.socket
}

// Returns a variable from the VARIABLES array.
//
// Note: variables prefixed with `env:` can be accessed as normal environment variables.
// This method will return such a variable by the name with the prefix still included.
// Generally it's better to access those variables directly.
func (p *BuildConfig) Variable(name string, defaultValue string) string {
	if val, ok := p.variables[name]; ok {
		return val
	}
	return defaultValue
}

// Returns the full variables array.
//
// If you're looking for a specific variable, the Variable() method is a more robust option.
// This method is for cases where you want to scan the whole variables list looking for a pattern.
func (p *BuildConfig) Variables() EnvList {
	return p.variables
}

// Retrieves the credentials for accessing a relationship.
func (p *RuntimeConfig) Credentials(relationship string) (Credential, error) {

	// Non-zero relationship indexes are not currently used, so hard code 0 for now.
	// On the off chance that ever changes, we'll add another method that allows
	// callers to specify an offset.
	if creds, ok := p.credentials[relationship]; ok {
		return creds[0], nil
	}

	return Credential{}, fmt.Errorf("No such relationship: %s", relationship)
}

// Returns the routes definition.
// This is an slice of Route structs.
func (p *RuntimeConfig) Routes() Routes {
	return p.routes
}

// Returns a single route definition.
//
// Note: If no route ID was specified in routes.yaml then it will not be possible
// to look up a route by ID.
func (p *RuntimeConfig) Route(id string) (Route, bool) {
	for _, route := range p.routes {
		if route.Id == id {
			return *route, true
		}
	}

	return Route{}, false
}

// Returns the definition of the primary route.
func (p *RuntimeConfig) PrimaryRoute() (Route, bool) {
	for _, route := range p.routes {
		if route.Primary == true {
			return *route, true
		}
	}

	return Route{}, false
}

func (p *RuntimeConfig) UpstreamRoutes() Routes {
	ret := make(Routes)

	for url, route := range(p.routes) {
		if route.Type == "upstream" {
			ret[url] = route
		}
	}

	return ret
}

// Map the relationships environment variable string into the appropriate data structure.
func extractCredentials(relationships string) (Credentials, error) {
	jsonRelationships, err := base64.StdEncoding.DecodeString(relationships)
	if err != nil {
		return Credentials{}, err
	}

	var rels Credentials

	err = json.Unmarshal([]byte(jsonRelationships), &rels)
	if err != nil {
		return nil, err
	}

	return rels, nil
}

// Map the variables environment variable string into the appropriate data structure.
func extractVariables(vars string) (EnvList, error) {
	jsonVars, err := base64.StdEncoding.DecodeString(vars)
	if err != nil {
		return EnvList{}, err
	}

	var env EnvList

	err = json.Unmarshal([]byte(jsonVars), &env)
	if err != nil {
		return nil, err
	}

	return env, nil
}

// Map the routes environment variable string into the appropriate data structure.
func extractRoutes(routesString string) (Routes, error) {
	jsonRoutes, err := base64.StdEncoding.DecodeString(routesString)
	if err != nil {
		return Routes{}, err
	}

	var routes Routes

	err = json.Unmarshal([]byte(jsonRoutes), &routes)
	if err != nil {
		return nil, err
	}

	// Normalize the URL of each route into the struct, so that it's available
	// when requesting a route individually.
	for url, _ := range routes {
		routes[url].Url = url
	}

	return routes, nil
}
