# Jablko

**UPDATE**: The new version of Jablko based on spawning subprocesses and communicating through HTTP has been merged onto the master branch. This implementation is still not finished, but a working version can be run using `go run main.go`. Right now, performance has not been optimized as no caching mechanisms are being used for critical file reads. Despite this, the time to interactive measured by Lighthouse is still < 0.8 seconds. With caching and bundling the `standard.css` file with the sent HTML files on initial request will significantly reduce load time.

Jablko is a smart home system written in Go that is extendible by user created Jablko Mods. The system is designed to be very simple, but offer flexibility to suit whatever needs you may come up with. The main server can communicate through network requests to any physical modules you may have, or you can use a custom communication protocol to communicate with smart home devices. User-written Jablko Mods provide an interface between your smart home dashboard and the rest of the world.

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

## Development

- core/app is deprecated. Work is currently using core/app2
