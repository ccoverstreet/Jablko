# Getting Started

Welcome to the Jablko Smart home project. The goal of this project is to create a fully open source smart home framework that uses your local wireless network to communicate with modules through a JSON-based API. The interface is powered by Deno and the Oak web server module and uses a SQLite3 database for managing user information and authentication. The smart home uses a GroupME bot to communicate with users and any Jablko Module can be modified to send messages through the GroupMe bot. Jablko modules with a permission level of "anyone" can also be controlled with commands sent through GroupMe. You can find more in **docs/jablko_modules.md**.

## Setting up jablko_config.json

 The "jablko_config.json" file contains info neccessary to run the GroupMe messaging bot and the list of Jablko Modules in the order you want them to appear on the dashboard. Here's the general format and necessary fields:
```
{
    "http": {
        "port": 8080
    },
    "https": {
        "port": null,
        "cert_file": "your_cert.pem",
        "key_file": "your_privkey.pem"
    },
    "database": {
        "path": "./database/primary.db",
        "session_lifetime": 21600000
    },
    "GroupMe": {
        "access_token": "access_token",
        "group_name": "Group Name",
        "group_id": "Group Id",
        "bot_id": "Bot Id"
    },
    "jablko_modules": {
        "interface_status": {
            "repo_archive": "https://github.com/ccoverstreet/Jablko-Interface-Status/archive/master.zip",
            "install_dir": "./jablko_modules"
        }
    },
    "weather": {
        "key": "OWM API KEY"
    }
}
```

Let's start with the first two keys: "http" and "https". These keys contain the port number for the http and https server and the SSL certificate files. If you aren't running the HTTPS server, then you can keep the "port" value for HTTPS as `null`. If not, you must provide a port number and the relative path to your SSL certificate files.

The database section does not need to be modified unless you wish to change the session_lifetime or have your own SQLite database setup. The "session_lifetime" determines how long each user can stay logged in.

Now, let's look at the GroupMe portion. The "access_token" is not needed, but I need to fix that in the future. You can find the required information on your GroupMe developer page (https://dev.groupme.com/bots). If you don't have a bot setup, you can use the online form to create a bot in an existing groupchat. 

Next is the jablko_modules configuration. Jablko Modules are located in the "jablko_modules" directory and are automatically installed when you run the command `./jpm init` in the root of Jablko. The previous command reads this config file and will download the source of each module from its respective repository. It then copies the contents to the key for each module in the subdirectory "jablko_modules". Each Jablko Module also needs a file "module.js" that contains the functions/exports needed to interface with the main interface.  You can find out more details and how to create a Jablko Module in the [documentation for Jablko Modules](/docs/jablko_modules.md).

Now, to add modules to Jablko, you can use the JPM utility (`./jpm`) in the root of Jablko. More information in the [Install Jablko Modules](#installing-jablko-modules)

The "weather" section just contains your OpenWeatherMap API key so that Jablko can retrieve weather information.

Great! Now all that's left to setup is the database setup for user information and managing active session.

## Setting up the User Database

### Easy Method

An easy way to setup the database is to run `node database/database_maintanence.js` in the root of Jablko. This will automatically start a database creation script that will prompt you for information. Here's what you need to do:
1. Use option 1 to create a database.
2. The database name should be "primary" (Will add customizability in the future).
3. Now use option 2 to add a user.
4. Make sure to use the primary database

### Manual Method
I'll add this documentation later. You can also just look through the database directory to figure it out. #BeingLazy

## Installing Jablko Modules

The primary method for installing Jablko Modules is to use the Jablko Package Manager (JPM) utility. To use, navigate to the root of Jablko and run `./jpm`. If you want to install a module, you can run `./jpm install author/repo tag module_name` to install a Jablko Module from a GitHub repository to your desired module name. For example, `./jpm install ccoverstreet/Jablko-Interface-Status v1.1.0 interface_status` will install the Version 1.1.0 of ccoverstreet/Jablko-Interface-status to a module named "interface_status". JPM automatically will install any dependencies and update your "jablko_config.json" file. You can read more about JPM [here](/docs/jpm.md).

## Running and Using the Interface
To run Jablko, all you need to do is run the command `./jablko` in the root of Jablko. To access the dashboard, open a browser and navigate to the http or https port you specified in "jablko_config.json". You can also use port forwarding on your router to enable access when outside your network. Remember that by forwarding the port, you are opening a potential attack vector to the outside world. As long as the modules you use don't have the capability to burn your house down (*cough* smart ovens *cough*), any potential attacks won't be able to do too much harm. You will need to use port forwarding (at least HTTPS only) to be able to use the GroupMe functionality of Jablko. 

As always, think about what security risks you are willing to take.


