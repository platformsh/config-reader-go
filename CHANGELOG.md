# Changelog

## [2.3.2] - 2021-02-03

### Added

* GitHub actions for tests (`quality-assurance.yaml`).

### Changed 

* named variable`prefix` on constructor renamed to `varPrefix`.

### Removed

* CircleCI action config. 

## [2.3.1] - 2019-11-04

### Added

* `CHANGELOG` added.
* `OnDedicated` method that determines if the current environment is a Platform.sh Dedicated environment. Replaces deprecated `OnEnterprise` method.

### Changed

* Deprecates `OnEnterprise` method - which is for now made to wrap around the added `OnDedicated` method. `OnEnterprise` **will be removed** in a future release, so update your projects to use `OnDedicated` instead as soon as possible.

## [2.3.0] - 2019-09-19

### Added

* `PrimaryRoute` method for accessing routes marked "primary" in `routes.yaml`.
* `UpstreamRoutes` method returns an object map that includes only those routes that point to a valid upstream.
* `UpstreamRoutesForApp(appName)` method returns an object map that includes only those routes that point to a valid upstream only for the current application where the code is being run.

## [2.2.3] - 2019-07-10

### Added

* Adds a number of formatted credentials packages:
  * `amqp`: produces the connection string for using the [AMPQ library](https://github.com/streadway/amqp) to connect to RabbitMQ.
  * `gomemcache`: produces a connection string for connecting to Memcached with the [gomemcache library](https://github.com/bradfitz/gomemcache).
  * `gosolr`: produces a connection string that includes the full collection path for using the [`go-solr` library](https://github.com/rtt/Go-Solr) to connect to Solr.
  * `libpq`: produces the [`lib/pq` library](https://github.com/lib/pq) connection string for PostgreSQL.
  * `mongo`: produces the connection string for using MongoDB's [`mongo-driver`](https://github.com/mongodb/mongo-go-driver) for Go.

### Changed

* Fixes issue where the `Route.Url` property wasn't being populated correctly.
* Updates `README` to correctly document `v2` imports.

## [2.2.2] - 2019-06-25

### Changed

* More gracefully handles the `PLATFORM_APPLICATION` variable being undefined on local environments. The other compound variables were already handled.

## [2.2.1] - 2019-06-21

### Changed

* Bumps CircleCI version to 1.12.
* Updates CircleCI `go get` to use modules.
* Uses `v2` syntax for imports.

## [2.2.0] - 2019-06-21

### Added

* `sqldsn` formatted credentials package. Produces an SQL connection string appropriate for use with many common Go database tools, including the [go-sql-driver](https://github.com/go-sql-driver/mysql).
* `BuildConfig` struct that provides access to the Platform.sh build environment context.
* `RuntimeConfig` struct that provides access to the Platform.sh runtime environment context.
