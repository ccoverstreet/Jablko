# Jablko

**UPDATE**: The new version of Jablko based on spawning subprocesses and communicating through HTTP has been merged onto the master branch. This implementation is still not finished, but a working version can be run using `go run main.go`. Right now, performance has not been optimized as no caching mechanisms are being used for critical file reads. Despite this, the time to interactive measured by Lighthouse is still < 0.8 seconds. With caching and bundling the `standard.css` file with the sent HTML files on initial request will significantly reduce load time.

Jablko is a smart home system written in Go that is extendible by user created Jablko Mods. The system is designed to be very simple, but offer flexibility to suit whatever needs you may come up with. The main server can communicate through network requests to any physical modules you may have, or you can use a custom communication protocol to communicate with smart home devices. User-written Jablko Mods provide an interface between your smart home dashboard and the rest of the world.

## News

0.3.0 is the current development goal. Issues and suggested features should go in issues. 

Different architectures are currently being tested. The master branch represents a Go only version using go plugins that are dynamically loaded. This design has the downside that the plugins have be built using the exact same dev setup and forces the restart of the main Jablko process if a module needs to be reloaded. 

An alternative, microservice-esque design is being developed to evaluate performance, scalability, and maintainability. The goal of this design is to decouple the actions of individual modules from the core Jablko process. In this design, Jablko acts more like a reverse-proxy layer that handles user authentication and provides a common target for jablko mods. Having a Jablko as the common framework that modules can call interact with makes it easier to provide a front-end for each module.

## Installing

### Prerequisites
- Go

### Instructions

1. `git clone https://github.com/ccoverstreet/Jablko`
2. `go build .`
3. `./Jablko`
4. Follow any directions that may pop up.


## Jablko Mods

In Progress

Mods must comply with the interfaces descrived in types/types.go

## Future work

- Messaging functionality (likely groupme, possible email) 
- Robust Jablko Mod manager.
  - Users should be able to install from dashboard (end goal)
  - Terminal usage is a short-term goal

## Reason for Using Go

The switch from NodeJS to Go was made to improve performance, increase stability, enforce a uniform Jablko Mod interface, and reduce development time. A major issue with the NodeJS version was that not all critical bugs were caught before runtime (even with TypeScript), leading to unpredictable errors and seemingly unrelated issues (switching arguments). This lead to prolonged development time as debugging seemingly innocent segments of code took much longer than expected. Using Go, a majority of errors is caught at compile-time and the strong, static-typing prevents dangerous operations from occuring. The language support for concurrency also makes it far easier to maximize performance and prevent race conditions. An unexpected side-effect of using Go was the reduction in lines-of-code needed to complete comparable tasks. This is not solely due to the change of language as the Javascript version was not as optimized as it could be; however, using Go forced a more concise design which makes future development much easier and maintanable.

Also NPM.
