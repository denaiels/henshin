# Henshin

Henshin is a tool to transform long URL to a shortened version of it.

For example you have a URL like this:

"https://source.golabs.io/daniel.santoso/url-blaster"

You can make a shorter version of it that can redirect back to the original URL like this:

"https://localhost:9080/Henshin" -> just an example

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

# Features

## Shorten URL

Run this command:

```sh-session
curl --request POST \
--data '{
    "long_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    "user_id" : "e0dba740-fc4b-4977-872c-d360239e6b10"
}' \
  http://localhost:9808/create-short-url
```

And open the short URL generated on your browser.

## Shorten URL with predefined string

Run this command:

```sh-session
curl --request POST \
--data '{
    "long_url": "https://ultra.fandom.com/wiki/Ultraman_Cosmos_(character)?file=Ultraman_Cosmos.png",
    "user_id" : "e0dba740-fc4b-4977-872c-d360239e6b10",
    "predefined_name" : "cosmos"
}' \
  http://localhost:9808/create-short-url
```

And open the short URL generated on your browser.

## Update URL pointed by the short URL

Run this command:

```sh-session
curl --request POST \
--data '{
    "short_url": "cosmos",
    "new_long_url" : "https://ultra.fandom.com/wiki/Ultraman_Cosmos_(character)?file=Cosmos_Luna_to_Corona.gif#Luna"
}' \
  http://localhost:9808/update-url
```

And open the short URL you updated on your browser.

## Remove shortened URL

Run this command:

```sh-session
curl --request POST \
--data '{
    "short_url": "cosmos"
}' \
  http://localhost:9808/remove-url
```

Try to open the short URL you removed using your browser, it will show 404 error.

# Used Technology

Go programming language.

Redis as store mechanism for super fast data retrieval.

Using https://golangci-lint.run/ to lint because it is fast, lots of linters (no need to install), integrates with VSCode, etc.

Hashing using SHA256 because it doesn't have any known vulnerabilities that make it insecure and it has not been “broken” unlike some other popular hashing algorithms. SHA256 shortens the input data into a smaller form that cannot be understood by using bitwise operations, modular additions, and compression functions.

Encoding using BASE58 instead of BASE64 because:
- Doesn't generate the characters "0", "O", "I", and "l" which are highly confusing when used in certain fonts and are even quite harder to differentiate for people with visual issues.
- Removing punctuations characters prevent confusion for line breakers.
- Double-clicking selects the whole number as one word if it's all alphanumeric.
