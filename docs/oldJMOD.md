# Jablko Mods (JMODs)

JMODs allow for users to add their own custom functionality to Jablko. JMODs are subprocesses spawned by Jablko that communicate through HTTP requests.

## Requirements

1. Use environment variables provided by Jablko to listen and communicate with the core
  - `JABLKO_CORE_PORT`: Port that Jablko Core is listening on
  - `JABLKO_MOD_PORT`: Port that the JMOD should listen on for HTTP data
  - `JABLKO_MOD_KEY`: Key used to authenticate requests sent by the JMOD to Jablko core
  - `JABLKO_MOD_DATA_DIR`: Directory that a JMOD can use to store persistent data/files
  - `JABLKO_MOD_CONFIG`: JSON string that contains configuration data for the JMOD

2. Handle HTTP Requests from Jablko Core
  - `/jmod/*...`
    - Requests sent from a client. Jablko handles authentication.
  - `/webComponent`
    - Should return the javascript for a WebComponent.
    - This WebComponent is used to create front-end representation of the JMOD
    - The WebComponent has extra requirements specified below
  - `/instanceData`
    - Should respond with a JSON array of objects
    - Each object is used to create a unique front-end instance
      - Passed to a member function `init` that takes the JMOD Name as the first argument and the config object as the second argument.


### WebComponents for Dashboard

WebComponents are used as a way of encapsulating all javascript used by a JMOD in the dash. As such, the WebComponents for JMODs should not rely on external manipulation and should only operate on DOM elements within their shadowRoot. To maintain uniformity in appearance, JMODs can add the link below in their shadowRoot to have access to colors and common structures. This import comes at no cost as `/assets/standard.css` is already loaded by Jablko

```html
<link rel="stylesheet" href="/assets/standard.css"></link>
```

An example implementation of a JMOD WebComponent can be found in [builtin/demo/webcomponent.js](/builtin/demo/webcomponent.js). It is critical that all functionality for a JMOD front-end is implemented inside of the WebComponent class. 

### Example WebComponent HTML

```html
<link rel="stylesheet" href="/assets/standard.css"></link>
<div class="jmod-wrapper">
	<div class="jmod-header" style="display: flex; ">
		<h1>Test Module</h1>
		<svg viewBox="0 0 360 360">
			<circle cx="180" 
				cy="180" 
				r="120" 
				stroke="var(--clr-accent)" 
				stroke-width="30"
				fill="transparent"/>		
		</svg>
	</div>
	<hr>
	<div class="jmod-body">
		<p>
		Demo for Jablko
		</p>
	</div>
</div>
```



