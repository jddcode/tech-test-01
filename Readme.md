# Weather Service

This service allows you to look up the weather for the next two days for cities aound the world.
You can run or build the service via the `makefile` and the following commands:

* `make` - Build the service locally. The binary will be located in the `cmd/tech-test` folder 
and named `service`
* `make test` - Run unit tests for the entire repository
* `make run` - Run the service on your local machine natively - it will  be listening on port `8080`
* `make docker` - Build the service and install it in an Alpine linux container. An image called
`tech-test:latest` will be available on your local system visible with the command `docker images`
* `make docker-run` - Build the service and run it from within a docker container.

### Cache

The system caches the result for each city via a simple in memory cache. The cache is part of the
handler and appears in a subfolder of the handler `internal/handler-weather/cache`. The cache is
fully multi-process safe and implements `sync.RWMutex` to control crashes due to conflicting
read/writes to the map which holds the data.

### Example URL

Here is a typical sample URL you can use to view the output:

`http://127.0.0.1:8080/weather?city=chicago`

## Improvements - commercialisation

If this were a piece of commercial software and not for a tech test I would implement the 
following extra features:

* Structured logging throughout with messages flagged at various levels eg. 
`INFO`, `DEBUG` or `ERROR`
* The use of environment variables for many settings throughout such as the number of days of
weather to return and possibly other useful settings such as the URLs for third party APIs. This 
would allow changes in future without having to go to deeply into the code.
* A `github/workflows` Github Actions yaml file to automatically build and deploy the docker
container into Kubernetes, AWS Lambda, Google Cloud Run or other environment.
