# Introduction
            .,*/(#####(/*,.                               .,*((###(/*.
        .*(%%%%%%%%%%%%%%#/.                           .*#%%%%####%%%%#/.
      ./#%%%%#(/,,...,,***.           .......          *#%%%#*.   ,(%%%#/.
     .(#%%%#/.                    .*(#%%%%%%%##/,.     ,(%%%#*    ,(%%%#*.
    .*#%%%#/.    ..........     .*#%%%%#(/((#%%%%(,     ,/#%%%#(/#%%%#(,
    ./#%%%(*    ,#%%%%%%%%(*   .*#%%%#*     .*#%%%#,      *(%%%%%%%#(,.
    ./#%%%#*    ,(((##%%%%(*   ,/%%%%/.      .(%%%#/   .*#%%%#(*/(#%%%#/,
     ,#%%%#(.        ,#%%%(*   ,/%%%%/.      .(%%%#/  ,/%%%#/.    .*#%%%(,
      *#%%%%(*.      ,#%%%(*   .*#%%%#*     ./#%%%#,  ,(%%%#*      .(%%%#*
       ,(#%%%%%##(((##%%%%(*    .*#%%%%#(((##%%%%(,   .*#%%%##(///(#%%%#/.
         .*/###%%%%%%%###(/,      .,/##%%%%%##(/,.      .*(##%%%%%%##(*,
              .........                ......                .......
An API generator for Go API development for [go8](https://github.com/gmhafiz/go8). Full documentation in that link.

# Features

This kit is composed of standard Go library together with some well-known libraries to manage things like router, database query and migration support.

- [x] Framework-less and net/http compatible handler
- [x] Router/Mux with [Chi Router](https://github.com/go-chi/chi)
- [x] Database Operations with [sqlx](https://github.com/jmoiron/sqlx)
- [x] Database migration with [golang-migrate](https://github.com/golang-migrate/migrate/)
- [x] Input [validation](https://github.com/go-playground/validator) that returns multiple error strings
- [x] Read all configurations using a single `.env` file or environment variable
- [x] Clear directory structure, so you know where to find the middleware, domain, server struct, handle, business logic, store, configuration files, migrations etc.
- [x] (optional) Request log that logs each user uniquely based on host address
- [x] CORS
- [x] Custom model JSON output
- [x] Filters (input port), Resource (output port) for pagination and custom response respectively.
- [x] Uses [Task](https://taskfile.dev) to simplify various tasks like mocking, linting, test coverage, hot reload etc

# Quick Start

Install tool. Recommended go >= v1.17

    go install github.com/gmhafiz/go8gen/cmd/go8@latest

Create new project using new argument followed by project name. Go modules name will be the same as project name.

    go8 new <projectName>

example:

    go8 new myReSTAPI

Optionally, specify go modules name

    go8 new <projectName> <moduleName>

example:

    go8 new myReSTAPI github.com/gmhafiz/rest

To run:

    cd <projectName>
    go run cmd/<projectName>/<projectName>.go
    
    2021/01/01 13:48:01 API version: v0.1.0
    2021/01/01 13:48:01 serving at 0.0.0.0:3080
    2021/01/01 13:48:01 path: /health/liveness method: GET
    2021/01/01 13:48:01 path: /health/readiness method: GET

    curl -XGET http://localhost:3080/health/liveness

    2021/01/01 14:12:09 "GET http://localhost:3080/health/liveness HTTP/1.1" from 127.0.0.1:32996 - 000 0B in 35.22µs

To add a domain,

    go8 domain <name>

example:

    go8 domain books

    go run cmd/<projectName>/<projectName>.go


# Version

### v0.14.8

  * Fix domain creation

### v0.14.5

 * New layout

### v0.4.0

 * Changed into a code generation tool
 * Go + Postgres + Chi Router + sqlx
 
### v.0.3.0

 * Go + Postgres + Chi Router + sqlx

### v0.2.0

 * elasticsearch support
 * jobs (queue)
 * gRPC support 
 
### v0.1.0

 * Initial starter kit boilerplate for Go + Postgres + Chi Router + SqlBoiler  API Development
 

# License

MIT

Refer to LICENSE file.