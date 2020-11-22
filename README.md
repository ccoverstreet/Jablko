# Jablko

Jablko is a smart home system written in Go that is extendible by user created Jablko Mods. The system is designed to be very simple, but offer flexibility to suit whatever needs you may come up with. The main server can communicate through network requests to any physical modules you may have, or you can use a custom communication protocol to communicate with smart home devices. User-written Jablko Mods provide an interface between your smart home dashboard and the rest of the world.

## Reason for Using Go

The switch from NodeJS to Go was made to improve performance, increase stability, enforce a uniform Jablko Mod interface, and reduce development time. A major issue with the NodeJS version was that not all critical bugs were caught before runtime (even with TypeScript), leading to unpredictable errors and seemingly unrelated issues. This lead to prolonged development time as debugging seemingly innocent segments of code took mych longer than expected. Using Go, a majority of errors is caught at compile-time and the strong, static-typing prevents dangerous operations from occuring. The language support for concurrency also makes it far easier to maximize performance and prevent race conditions. An unexpected side-effect of using Go was the reduction in lines-of-code needed to complete comparable tasks. This is not solely due to the change of language as the Javascript version was not as optimized as it could be; however, using Go forced a more concise design which makes future development much easier and maintanable.

## Installing

In Progress

## Jablko Mods

In Progress

## Future work

- Messaging functionality (likely groupme, possible email) 
- Robust jablko mod manager.
  - Users should be able to install from dashboard (end goal)
  - Terminal usage is a short-term goal
  
