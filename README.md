# Project lair-api

<!--toc:start-->
- [Project lair-api](#project-lair-api)
  - [Getting Started](#getting-started)
  - [MakeFile](#makefile)
<!--toc:end-->

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See
deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests

```sh
make all build
```

build the application

```sh
make build
```

run the application

```sh
make run
```

Create DB container

```sh
make docker-run
```

Shutdown DB container

```sh
make docker-down
```

live reload the application

```sh
make watch
```

run the test suite

```sh
make test
```

clean up binary from the last build

```sh
make clean
```

