# Payments API

## Usage

Before running the tests or the app, run `make gen` to generated the models code. (see ##How-this-works)

Run `make run` to start server with docker (requires: docker, docker-compose).

(alternative: run `go run main.go` (requires go 1.11+, postgres installed locally and tweaking config files))

The API is accessible at http://localhost:36480

#### Requirements

- Go 1.11+
- docker
- docker-compose

## Testing

Run tests with `make test`

## How this works

You'll notice models are entirely absent from this repo. That's because we'll get them using code generation.
Models (and their tests) are generated from sql migrations applied to the database, which means:

- code generation must happen in docker since the database runs in docker
- the generated code gets written back to the host via a docker volume
- a second `docker build` is necessary to compile the Go app including the newly generated code

So, before you can run `make test` or `make run`, 
you have to run `make gen` once after cloning the project.
If you change the database migrations folder, you'll have to run `make gen` again to update the generated models.

Generated models are not committed to source control.

## License: MIT



    