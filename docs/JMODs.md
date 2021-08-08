# Jablko Mods (JMODs)

JMODs allow for users to add their own custom functionality to Jablko. JMODs are subprocesses spawned by Jablko that communicate through HTTP requests.

## Requirements

1. Use environment variables provided by Jablko to listen and communicate with the core
  - `JABLKO_CORE_PORT`: Port that Jablko Core is listening on
  - `JABLKO_MOD_PORT`: Port that the JMOD should listen on for HTTP data
  - `JABLKO_MOD_KEY`: Key used to authenticate requests sent by the JMOD to Jablko core
  - `JABLKO_MOD_DATA_DIR`: Directory that a JMOD can use to store persistent data/files
  - `JABLKO_MOD_CONFIG`: JSON string that contains configuration data for the JMOD

