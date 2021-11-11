# Introduction

```
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
```
A code generation tool to quickly scaffold a Go API application. Inspired by  [How I write HTTP
 services after eight years](https://pace.dev/blog/2018/05/09/how-I-write-http-services-after
 -eight-years.html) and incorporates ideas from [go-clean-architecture](https://github.com/zhashkevych/go-clean-architecture)

This project has changed from being a boilerplate into a code generator.

# Quickstart

Install tool. Requires go >= v1.13

    go install github.com/gmhafiz/go8gen/cmd/go8@latest
    
Create new project using `new` argument followed by project name. Go modules name will be the
 same as project name. 

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
 
# Motivation

On the topic of API development, there are two opposing camps between a framework (like echo, gin
, buffalo and starting small and only add features you need. However , starting small and adding
 features aren't that straightforward. Also, you will want to structure your project in such a
  way that there are clear separation of functionalities for different files. This is the idea
   behind [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html). This way, it is easy to switch whichever library to another of your choice.
   
   
# Acknowledgements

 * https://github.com/joshuabezaleel/unit-testing-golang-mocks