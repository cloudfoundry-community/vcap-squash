# vcap-squash [![Build Status](https://travis-ci.org/joshq00/vcap-squash.svg?branch=master)](https://travis-ci.org/joshq00/vcap-squash)
Turn VCAP_SERVICES into flat env vars

## Usage
This application will parse the `VCAP_SERVICES` environment variable and output the unix exports of the flattened version.

Flattened vcap environment variables start with the service name and append `_` for each nested credential.
For example:
```sh
$ VCAP_SERVICES='{ "user-provided": [ {
  "name": "myservice",
  "credentials": {
    "url": "myservice.com",
    "username": "josh",
    "password": "secret",
    "nested": {
      "key": "value",
      "number": 123
    }
  }
} ] }' ./vcap-squash

export MYSERVICE_URL="myservice.com"
export MYSERVICE_USERNAME="josh"
export MYSERVICE_PASSWORD="secret"
export MYSERVICE_NESTED_KEY="value"
export MYSERVICE_NESTED_NUMBER=123
```

To set your environment variables using this output, use
```sh
eval $(./vcap-squash)
```

### In Cloud Foundry
Add the proper `vcap-squash` binary to your project root (depending on cf stack)

Create a `.profile.d/setenv.sh` file to push along with your repo
```sh
#!/bin/sh
eval $(./vcap-squash)
```

## Development
To run all basic tasks, use
```sh
$ make
```

### Dependencies
To download needed libraries:
```sh
$ make deps
```

### Test
To run the test suite, use
```sh
$ make test
```

Run the test suite in development/watch mode:
```sh
$ make watch
```

### Build local
To build a binary using your go env:
```sh
$ make build-local
```

### Build all
Build the binary for all systems:
```sh
$ make build
```
_the binaries will be placed in `./out`_

