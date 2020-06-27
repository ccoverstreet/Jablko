# Jablko Smart Home

## About

Jablko is a smart home system that uses Deno and Oak to run an interface server that communicates to other devices on the local network to control devices like lights/sensors. The goal is to have a fully customizable system where any device that can communicate on the local network through JSON requests and responses can be linked to the interface through its own Jablko Module. 

The Jablko Module system is designed to make adding devices as painless as transparent as possible (See [Jablko Modules Section](#Jablko-Modules) for more info).

Windows compatibility is not considered for now as accessing devices and file paths differences unnecessarily increases complexity for this project. Adding Windows compatitbility in the future may be nice if some users have a main "media" computer. 

## Starting Jablko

Before running Jablko for the first time, make sure to run `./setup_jablko.sh`. This script will prompt the user to create the database Jablko uses.

Run the command `./start_jablko.sh`. Make sure to make the bash script executable beforehand.

## Jablko Modules

The smart home interface communicates to other devices on the network through requests using a JSON API that I will flesh out in the future (Can't test in current place due to horrible apartment internet configuration). Inside the root directory of Jablko is a "jablko_modules" folder. Modules should be created with the following convention:

- Main TypeScript file: /jablko_modules/module_name/module_name.ts
  - Must have a generate_card() function that returns a string that contains the HTML for the module's card on the interface.
  - Routes intended for the module should be formatted as /jablko_modules/module_name/exported_function name. Jablko will identify requests formatted this way and call the requested Jablko Module Function.
    - Routes should take an Oak Context as an argument and must handle the context correctly.
- Auxiliary files: Placed in /jablko_modules/module_name

Jablko automatically loads any modules with the above structure on startup. This means that any changes to the module_name/module_name.ts file will not take effect until Jablko is restarted. 
