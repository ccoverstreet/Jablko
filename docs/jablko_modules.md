# Jablko Modules

Jablko Modules can be made by creating a directory or repository that has a "module.js" file and NPM package.json file. There are some required `module.exports` needed in order for the module to interface with Jablko.

## Contents

- [Overview](#overview)
- [Jablko Module Standards](#jablko-module-standards)
- [General HTML Design](#general-html-design)
- [Development Setup](#development-setup)

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

