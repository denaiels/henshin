# URL Blaster

URL Blaster is a tool to shorten URLs.

For example you have a link like this:

"https://source.golabs.io/daniel.santoso/url-blaster"

You can make a shorter version of it like this:

"https://localhost:9080/URLBlaster" -> just an example

# Installation

## Prerequisites

1. Go 1.18
2. Docker compose

## Clone Repo

In your terminal, go to the directory where you want to clone the repo

```sh-session
cd Your/Directory/Of/Choice
```

Then type this in your terminal and hit enter:

```sh-session
git clone https://source.golabs.io/daniel.santoso/url-blaster.git
```

## Setup

To setup, please run the following commands:

Install package dependencies

```sh-session
make setup
```

## Building, Testing, and Running the service

To build the service, please run the following commands:

```sh-session
# This will generate an executable in ./out directory.
make build

# This will run all the tests.
make test
```

## Running the service

To run the service, please run the following commands:

```sh-session
make run
```

# Used Technology

Go programming language.

Redis as store mechanism for super fast data retrieval.

https://golangci-lint.run/ -> fast, lots of linters (no need to install), integrate with VSCode, etc.
