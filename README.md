# Jablko Smart Home

## About

Jablko is a smart home system that uses NodeJS and the Express module to run an interface server that communicates to other devices on the local network to control devices like lights/sensors. Jablko has a web interface exposed on an HTTP and optional HTTPS port users can specify in the configuration file. The goal is to have a fully customizable system where any device that can communicate on the local network through JSON requests and responses can be linked to the interface through its own Jablko Module and where anyone is able to write their own Jablko Module to integrate Jablko and the real world.

The Jablko Module system is designed to make adding devices as painless as transparent as possible (See [Jablko Modules Section](#Jablko-Modules) for more info).

Windows compatibility is not considered for now as accessing devices and file paths differences unnecessarily increases complexity for this project. Adding Windows compatitbility in the future may be nice if some users have a main "media" computer. 

## Getting Started
[Documentation](docs/getting_started.md)

## Jablko Modules

[Documentation](docs/jablko_modules.md)
