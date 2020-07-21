# Getting Started

Welcome to the Jablko Smart home project. The goal of this project is to create a fully open source smart home framework that uses your local wireless network to communicate with modules through a JSON-based API. The interface is powered by Deno and the Oak web server module and uses a SQLite3 database for managing user information and authentication. The smart home uses a GroupME bot to communicate with users and any Jablko Module can be modified to send messages through the GroupMe bot. Jablko modules with a permission level of "anyone" can also be controlled with commands sent through GroupMe. You can find more in **docs/jablko_modules.md**.

## Setting up jablko_config.json

 The "jablko_config.json" file contains info neccessary to run the GroupMe messaging bot and the list of Jablko Modules in the order you want them to appear on the dashboard. Here's the general format and necessary fields:
```
{
  "GroupMe: {
    "access_token": "yourToken",
    "group_name": "YourGroupName",
    "group_id": "YourGroupId",
    "bot_id": "YourBotsId"
  },
  "jablko_modules": [
    "theFirstModule",
    "theSecondModule
  ]
}
```

Let's start with the GroupMe portion. You can find the required information on your GroupMe developer page (https://dev.groupme.com/bots). If you don't have a bot setup, you can use the online form to create a bot in an existing groupchat. 

Next is the jablko_modules configuration. Jablko Modules are located in the "jablko_modules" directory and should be in their own subdirectory (jablko_modules/mymodule/*). Each Jablko Module also needs a file "mymodule.ts" that contains the functions/exports needed to interface with the main interface. This file **MUST** be named the same as the containing directory and entry in the "jablko_modules" section of the config file. You can find out more details and how to create a Jablko Module in the [documentation for Jablko Modules](docs/jablko_modules.md).

Now, to add modules to Jablko, just add the name of the module into the "jablko_modules" array. When you restart the interface, the module will be loaded and provided the module is made correctly you should see it on your dashboard and be able to use any established routes.

Great! Now all that's left to setup is the database setup for user information and managing active session.

## Setting up the User Database

### Easy Method

An easy way to setup the database is to run `./setup_jablko.sh` in the root of Jablko. This will automatically start a database creation script that will prompt you for information. Here's what you need to do:
1. Use option 1 to create a database.
2. The database name should be "primary" (Will add customizability in the future).
3. Now use option 2 to add a user.
4. Make sure to use the primary database

### Manual Method
I'll add this documentation later. You can also just look through the database directory to figure it out.# Getting Started

## Running and Using the Interface
To run Jablko, all you need to do is run the command `./start_jablko.sh`. To access the dashboard, open a browser and navigate to http://localhost:10230. You can also use port forwarding on your router to enable access when outside your network. Remember that by forwarding the port, you are opening a potential attack vector to the outside world. As long as the modules you use don't have the capability to burn your house down (*cough* smart ovens *cough*), any potential attacks won't be able to do too much harm. 

As always, think about what security risks you are willing to take.


