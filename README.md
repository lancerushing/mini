[![Go Report Card](https://goreportcard.com/badge/lancerushing/mini)](https://goreportcard.com/report/lancerushing/mini)

[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)

# "Mini" web application

A small web application used to skeleton a new project.

## Features

* New user signup
* User login
* Forgot password functionality

## Requirements

* go >1.16
* postgresql server
* make (useful)

## Development

```bash
go get github.com/cortesi/modd/cmd/modd
go get golang.org/x/lint/golint

curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.35.2
```


```bash
make create-db  # first time only

# optional for dev, if you want email to work
export SENDGRID_API_KEY="Your-key-here" 

make dev
```

Then open your browser to http://localhost:4000


#### Prerequisites

* Install Make ( apt install make) 
* Install go ( https://golang.org/doc/install ) 
* postgres ( https://www.postgresql.org/download/linux/debian )
* gin `$ go get github.com/codegangsta/gin`
* ci-lint `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`

## Todo

* How to handle configs? yaml,env,build flags,...,???
  * Eventually the app will need ~10 config vars
    * 5 for DB connection
    * 1 for send gride
    * 4 for the secure cookie inputs
    * 1 for schema://domain_name:port
* How to handle logging?
  * logrus
  * chi mux logging middleware
* Deploying?
  * Do we bundle templates into the binary?
  * deploy to /opt or /usr/local/{app} ?
  * separate config file in /etc ? or shove all ENV vars into .service
  