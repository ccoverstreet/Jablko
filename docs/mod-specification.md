# Mod Specification

## Backend

This is the specification for the process that Jablko spawns and represents the brains of the mod. The mod has the following requirements

- HTTP API
	- `/mod/<some route for functions>`
		- General requests to be passed through from Jablko to the mod must be prefixed with the `/mod` route.
		- The target mod is specified by passing the mod name as a URL argument (specified in [frontend](#frontend))
	- `/mod/webComponent`
		- Should return a string containing the javascript for a webComponent

## Frontend

- To target the specified mod from the frontend or any other client, the mod name must be specified as a URL argument 
	- ex. `/mod/getTemperature?modName=somemodname`
- The webComponent for the frontend display card for each module must have the following methods
	- `constructor()`
	- `init(modName)`
		- Takes in the modName as an argument and stores it
		- This modName value is used when sending requests as specified above
	- It is also recommended to use a connectedCallback to set the innerHTML of the element
- 

