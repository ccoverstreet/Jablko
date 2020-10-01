# Jablko Modules

Jablko Modules can be made by creating a directory or repository that has a "module.js" file and NPM package.json file. There are some required `module.exports` needed in order for the module to interface with Jablko.

## Contents

- [Overview](#overview)
- [Jablko Module Standards](#jablko-module-standards)
- [General HTML Design](#general-html-design)
- [Development Setup](#development-setup)
- [Module Routes](#module-routes)

## Overview

Components:
- "module.js": Must be located in the root of the Jablko Module's repository.
  - `module.exports.permission_level`: A exported int set to an integer value (0, 1, or 2). Permission level needed to use the module increase with the number. The level 2 denotes administrative level and 0 represents no required permissions.
    ```Javascript
    module.exports.permission_level = 0
    ```
  - RECOMMENDED: `const module_name = path.basename(__dirname)`. Module name resolves to the name of the install directory of the module. This makes it so that every module has a unique identifier and there are no function/definition collisions between modules when they are on the dashboard.
  - RECOMMENDED: `const jablko = require(module.parent.filename)`. This is used to import any functions exposed by the main Jablko Interface.
  - RECOMMENDED: `const module_config = jablko.jablko_config.jablko_modules[module_name]`. Contains the configuration data pulled from the "jablko_config.json" file.
  - `module.exports.generate_card`: If you want a card to appear on your dashboard you must provide an `async generate_card()`. function. This function is called when a user visits the server route and returns a string containing the HTML of the modules card. How you generate the string is up to each individual module. The recommended method is to read from an html file in the module.js file's directory and replace any occurences of $MODULE_NAME with the module's installed directory name. You can see more in [HTML Standards](#html-standards)
    ```Javascript
      module.exports.generate_card = async function generate_card() {
        return (await fs.readFile(`${__dirname}/mantle_rgb.html`, "utf8")).replace(/\$MODULE_NAME/g, module_name);
      }
    ```
```Javascript
const module_name = path.basename(__dirname);
const jablko = require(module.parent.filename);
const module_config = jablko.jablko_config.jablko_modules[module_name];
```
    
## Jablko Module Standard

There are several key design principles developers should follow when creating a Jablko module:
- All ids and unique identifiers in the HTML should be prefaced with `$MODULE_NAME` (ex. `id="$MODULE_NAME_mydiv"`, which should be replaced with `module_name` when the HTML file is read in and returned.
- All JavaScript definitions should also be namespaced using a const $MODULE_NAME object and be called from that object.
    ```Javascript
    const $MODULE_NAME = {
      myfunc: async () => {
        console.log("myfunc");
      }
    }
    
    // In module <div>
    <button onclick="$MODULE_NAME.myfunc()">My Button</button>
    ```
- Any variables you replace in the html file from some form of config should be of the form `$VARIABLE_NAME` to make them easily identifiable.
- If you have config options or default config values, they **MUST** be set in the package.json in a field called `"jablko"`.
  - If you have config values that must be correct or manually set, make sure to add config validation to the very start of your module that throws an error if the config is invalid or missing.
  
**TIP** If you need examples, look at any of the official Jablko Modules listed on the [main README.md](/README.md). A simple one that uses most features is [Jablko-Interface-Status](https://github.com/ccoverstreet/Jablko-Interface-Status)

## General HTML Design

```HTML
<script>
  const $MODULE_NAME = { // Module name replaced by regex in you async generate_card function
    somefunc: async () => {
      await fetch("/jablko_modules/$MODULE_NAME/your_exposed_func", {method: "POST", body: {YOUR DATA HERE}})
        .catch((error) => {
          console.log(error);
        });
      const value_div = document.getElementById("$MODULE_NAME_somevalue");
      console.log(value_div.textContent);
      alert($SOME_CONFIG_VALUE); // Replaced by regex in your async generate_card function
    }
  }
</script>
<div id="$MODULE_NAME" class="module_card">
  <div class="module_header">
    <h1>Module</h1>
      <svg class="module_icon" viewBox="0 0 150 150">
      <path d="M 20 75 H 40 L 60 30 L 90 120 L 110 75 H 130" fill="transparent" stroke="#0097e6" stroke-width="20px" stroke-linejoin="round" stroke-linecap="round"/>
    </svg>
  </div>
  <div id="$MODULE_NAME_somevalue">100</div>
  <button onclick="$MODULE_NAME.somefunc()">Get Value</button>
</div>
```

**CSS**: Available CSS classes/presets are in [dashboard.css](/public_html/dashboard/dashboard.css)

## Development Setup

To setup up Jablko so that you can use a separate local repository for development just involves changing the `install_dir` field in the corresponding module in "jablko_config.json" and ensuring that your containing directory is valid (No spaces or dashes):

```JSON
{
  ...
  "jablko_modules": {
    "your_dev_module": {
      "repo_archive": null,
      "install_dir": "/some/absolute/path/to/your_dev_module",
      "some_config_value": true
    }
  }
}
```
**NOTE** The module name and containing directory name must be the same.

Once this is configured, Jablko will load the module from the specified directory instead of through the typical "jablko_modules" install path. 

## Module Routes

Module routes allow for Jablko to automatically pass requests to your desired module function. Your function must handle the req and send a proper res object. Module routes can be used to handle requests from both user dasboards and wifi-connected modules. If you do not wish to include authentication when having your wifi-connected module send requests to Jablko, you **MUST** prepend your path with `/module_callback`. If not, your request will be ignored. Your request will also be ignored if the request does not come from within the wifi network. 

Calling a route from the user dashboard should be done using `/jablko_modules/$MODULE_NAME/my_exported_path` where `$MODULE_NAME` is string replaced in your `async function generate_card()`. If you want a wifi-connected module to use the path and the device is on the same network as the computer running Jablko, you **MUST** prepend `/module_callback` to your request (which gives `/module_callback/jablko_modules/$MODULE_NAME/my_exported_path`). If not, Jablko will treat your request as unauthenticated. This type of request would be made if a module wants to send data to Jablko for either storage, proxying, or parsing. For example, a temperature logging device would use this type of callback to send data to Jablko where it could be logged to a file or prompt a warning to all users. 

Example:
```Javascript
// Should be async function most of the time
module.exports.my_exported_path = async function(req, res) {
  console.log(req.body); // Just print data sent for fun.
  
  // Do whatever you need for your module
  
  res.json({status: "good", message: "Did the thing"}); // Send response back to client
}

module.exports.update_config = async function(req, res) {
  // This type of function should only be used if you need to dynamically update a config value (say stored in EEPROM) of your module
  // Otherwise, you should prioritize making your modules not need dynamic config updates.
  
  await fetch("http://10.0.0.60/update_config", {
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json"
    },
    body: JSON.stringify(req.body)
  })
    .catch((error) => {
      console.log("Error updating module config of my_module"); // Pretty output when running in normal output
      console.debug(error); // Only shows full error in debug mode
    }
}
```

The `module.exports.update_config` function would be used in a case where you store the module name in the EEPROM of a microcontroller. This can be useful since Jablko Modules can be installed under any name and if the module uses the `/module_callback`, it needs to have the correct module name for the request to go through. As long as you keep usage of this function low (say 1 call per day) the EEPROM would be expected to survive a few hundred years. 


