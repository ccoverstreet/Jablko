# Jablko

**UPDATE**: The new new version will use Docker containers for the JMODs as I realize I have been rewriting aspects of Docker (particularly file access management). Switching to Docker would also give an easy JMOD package management system as users can simply specify the images they want. 

## Installing

### Prerequisites
- Go

### Instructions

1. `git clone https://github.com/ccoverstreet/Jablko`
2. `go build .`
3. `./Jablko`
4. Follow any directions that may pop up.


## Jablko Mods

[Documentation](/docs/JMODs.md)

## Future work

- Messaging functionality (likely groupme, possible email) 
- Robust Jablko Mod manager.
  - Users should be able to install from dashboard (end goal)
  - Terminal usage is a short-term goal

## Development

- core/app is deprecated. Work is currently using core/app2
