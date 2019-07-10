# Platform.sh Config Reader (Go)

This library provides a streamlined and easy to use way to interact with a Platform.sh environment. It defines structs for Routes and Relationships and offers utility methods to access them more cleanly than reading the raw environment variables yourself.

This library is best installed using Go modules in Go 1.11 and later.

## Install

Add a dependency on `github.com/platformsh/config-reader-go` to your application. We recommend giving it an explicit import name.

## Usage Example

Example:

```go
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	sqldsn "github.com/platformsh/config-reader-go/v2/sqldsn"
	psh "github.com/platformsh/config-reader-go/v2"
	"net/http"
)

func main() {

	// Creating a psh.RuntimeConfig struct
	config, err := psh.NewRuntimeConfig()
	if err != nil {
		panic("Not in a Platform.sh Environment.")
	}

	// Accessing the database relationship Credentials struct
	credentials, err := config.Credentials("database")
	if err != nil {
		panic(err)
	}

	// Using the sqldsn formatted credentials package
	formatted, err := sqldsn.FormattedCredentials(credentials)
	if err != nil {
		panic(err)
	}

  // Connect to the database using the formatted credentials
	db, err := sql.Open("mysql", formatted)
	if err != nil {
		panic(err)
	}

  // Use the db connection here.

	// Set up an extremely simple web server response.
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		// ...
	})

    // Note the Port value used here.
	http.ListenAndServe(":"+config.Port(), nil)
}
```

## API Reference

### Create a config object

There are two separate constructor functions depending on whether you intend to be in a build environment or runtime environment.

```go
// In a build hook, run:
buildConfig, err := psh.NewBuildConfig()
if err != nil {
    panic("Not in a Platform.sh Environment.")
}
```

`buildConfig` is now a `psh.BuildConfig` struct that provides access to the Platform.sh build environment context.  If `err` is not `nil` it means the library is not running on Platform.sh, so other commands would not run.

```go
// At runtime, run:
runtimeConfig, err := psh.NewRuntimeConfig()
if err != nil {
    panic("Not in a Platform.sh Environment.")
}
```

`runtimeConfig` is now a `psh.RuntimeConfig` struct that provides access to the Platform.sh runtime environment context.  That includes everything available in the Build context as well as information only meaningful at runtime.

### Inspect the environment

The following methods return `true` or `false` to help determine in what context the code is running:

```go
runtimeConfig.OnEnterprise()

runtimeConfig.OnProduction()
```

### Read environment variables

The following methods return the corresponding environment variable value.  See the [Platform.sh documentation](https://docs.platform.sh/development/variables.html) for a description of each.

The following are available both in Build and at Runtime:

```go
buildConfig.ApplicationName()

buildConfig.AppDir()

buildConfig.Project()

buildConfig.TreeId()

buildConfig.ProjectEntropy()
```

The following are available only on a `RuntimeConfig` struct:

```go
runtimeConfig.Branch()

runtimeConfig.DocumentRoot()

runtimeConfig.SmtpHost()

runtimeConfig.Environment()

runtimeConfig.Socket()

runtimeConfig.Port()
```

### Reading service credentials

[Platform.sh services](https://docs.platform.sh/configuration/services.html) are defined in a `services.yaml` file, and exposed to an application by listing a `relationship` to that service in the application's `.platform.app.yaml` file.  User, password, host, etc. information is then exposed to the running application in the `PLATFORM_RELATIONSHIPS` environment variable, which is a base64-encoded JSON string.  The following method allows easier access to credential information than decoding the environment variable yourself.

```go
if creds, ok := runtimeConfig.Credentials("database"); ok {
	// ...
}
```

The return value of `Credentials()` is a `Credential` struct, which includes the appropriate user, password, host, database name, and other pertinent information.  See the [Service documentation](https://docs.platform.sh/configuration/services.html) for your service for the exact structure and meaning of each property.  In most cases that information can be passed directly to whatever other client library is being used to connect to the service.

If `ok` is false it means the specified relationship was not defined so no credentials are available.

## Formatted service credentials

In some cases the library being used to connect to a service wants its credentials formatted in a specific way; it could be a DSN string of some sort or it needs certain values concatenated to the database name, etc. For those cases you can use "Credential Formatters".  A Credential Formatter is a package within `config-reader-go` that contains a function that takes a `Credential` object and returns the specified type for the library it connects to.

This library comes with a few formatters out of the box:

* `amqp`: produces the connection string for using the [AMPQ library](https://github.com/streadway/amqp) to connect to RabbitMQ.
* `gomemcache`: produces a connection string for connecting to Memcached with the [gomemcache library](https://github.com/bradfitz/gomemcache).
* `gosolr`: produces a connection string that includes the full collection path for using the [`go-solr` library](https://github.com/rtt/Go-Solr) to connect to Solr.
* `libpq`: produces the [`lib/pq` library](https://github.com/lib/pq) connection string for PostgreSQL.
* `mongo`: produces the connection string for using MongoDB's [`mongo-driver`](https://github.com/mongodb/mongo-go-driver) for Go.
* `sqldsn`: produces an SQL connection string appropriate for use with many common Go database tools, including the [go-sql-driver](https://github.com/go-sql-driver/mysql).

A formatter package can be used in your application by importing it

```go
import (
	sqldsn "github.com/platformsh/config-reader-go/v2/sqldsn"
)
```

and passing a `Credentials` struct to the package's `FormattedCredentials()` function.

```go
formatted, err := sqldsn.FormattedCredentials(credentials)
```

### Registering Credential formatters

Unlike Platform.sh's other Config Reader libraries, `config-reader-go` does not include an equivalent `RegisterFormatter` function for registering new formatters due to Go's reliance on package imports and type preservation.

New formatter packages will be added periodically, but open a pull request if you would like to see a utility function you have written to connect to a service library in the `config-reader-go` namespace.

### Reading Platform.sh variables

Platform.sh allows you to define arbitrary variables that may be available at build time, runtime, or both.  They are stored in the `PLATFORM_VARIABLES` environment variable, which is a base64-encoded JSON string.  

The following two methods allow access to those values from your code without having to bother decoding the values yourself:

```go
runtimeConfig.Variables()
```

This method returns a `map[string]string` of all variables defined.  Usually this method is not necessary and `config.Variable()` is preferred.

```go
runtimeConfig.Variable("foo", "default")
```

This method looks for the "foo" variable.  If found, it is returned.  If not, the second parameter is returned as a default.

Note that both methods are available on both Build and Runtime, although different values may be defined and avaialble for use.

### Reading Routes

[Routes](https://docs.platform.sh/configuration/routes.html) on Platform.sh define how a project will handle incoming requests; that primarily means what application container will serve the request, but it also includes cache configuration, TLS settings, etc.  Routes may also have an optional ID, which is the preferred way to access them.

```go
runtimeConfig.Route("main")
```

The `Route()` method takes a single string for the route ID ("main" in this case) and returns the corresponding `route` struct.  Its second return is a boolean indicating if the route was found.  It is best used like so:

```go
if route, ok := runtimeConfig.Route("main"); ok {
	// The route was found, so do stuff with `route`
}
```

To access all routes, or to search for a route that has no ID, the `Routes()` method returns a `map[string]Route` of URLs to `Route` objects.  That mirrors the structure of the `PLATFORM_ROUTES` environment variable.

```go
routes := runtimeConfig.Routes()
```
