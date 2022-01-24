# Jablko Specification

This is the primary specification document for the Docker rewrite of Jablko.


## Jablko Structure

The core process serves as a manager for the Docker processes spawned. Each spawned docker process is called a JMOD and has several requirements for interop with Jablko. The core of Jablko serves primarily as HTTP proxy server that takes care of authentication, process management, and frontend integration which removes some responsibility from each JMOD. 

Each JMOD is represented by a card on the dashboard, the JavaScript WebComponent for which is provided by each JMOD. Each JMOD's webcomponent is rendered on client page load (this may be changed if static webcomponent rendering becomes easier).

## JMOD Interface Requirements

The responsibilities for JMODs are summarized in Table 1.

On startup several environment variables are provided to each JMOD: Jablko Core IP address and port (ex. 192.168.1.203:8080), JMOD key (used to authenticate a JMOD with the core process). Jablko will not enforce SSL connections between JMODs and the core as JMODs should be running on localhost. If someone has access to the connections, then you have much bigger problems to worry about.

The config file for each JMOD will be stored in the mounted `/data` directory within the image as `jmodconfig.json`. On first JMOD install, this file will not exist and it is the JMOD's responsibility to create this file and ensure it is named correctly. User-specified config changes work by allowing the user to edit the current version of this file and upon application, Jablko will stop the JMOD and restart the image with the new config file in the `/data` directory.

The webcomponent for each JMOD is to be POSTed to the core process on the route `/service/webcomponent` by each JMOD after JMOD startup. Jablko will cache the webcomponents file and cash the integrated dashboard code.
