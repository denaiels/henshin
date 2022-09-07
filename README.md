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
# OR
./bin/url-blaster
```

# Used Technology

Go programming language.

Redis as store mechanism for super fast data retrieval.

https://golangci-lint.run/ -> fast, lots of linters (no need to install), integrate with VSCode, etc.

Hashing using SHA256 because it doesn't have any known vulnerabilities that make it insecure and it has not been “broken” unlike some other popular hashing algorithms. SHA256 shortens the input data into a smaller form that cannot be understood by using bitwise operations, modular additions, and compression functions.

Encoding using BASE58 instead of BASE64 because:
- The characters 0,O, I, l are highly confusing when used in certain fonts and are even quite harder to differentiate for people with visuality issues.
- Removing ponctuations characters prevent confusion for line breakers.
- Double-clicking selects the whole number as one word if it's all alphanumeric.
