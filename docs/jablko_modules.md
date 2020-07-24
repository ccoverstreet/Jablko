# Jablko Modules

Jablko Modules can be made by creating the required typescript file with external functions thay the main interface can call. Below is a list of items that **MUST** be defined in a module.

## Contents

- [Overview](#overview)

## Overview

Required components:
- TypeScript File: Must be located in the jablko_modules directory in a subdirectory with the same name (e.g. jablko_modules/mymodule/mymodule.ts).
  - A exported const info object with a permission_level member set to an integer value (0, 1, or 2). Permission level needed to use the module increase with the number. The level 2 denotes administrative level and 0 represents no required permissions.
    ```Javascript
    export const info = {
      permissions: 0
    }
    ```
  - If you want a card to appear on your dashboard you must provide an `async generate_card()`. function. This function is called when a user visits the server route and returns a string containing the HTML of the modules card. How you generate the string is up to each individual module. You can store the html in a separate .html file and then read in the html on load.
    ```Javascript
    export async function generate_card() {
      return `
      <div id="mymodule_card" class="jablko_module_card">
        <div class="card_title" background: url('/icons/interface_status_icon.svg') right; background-size: contain; background-repeat: no-repeat;">MyModule</div>
        <hr>
        <p>*Your Content*</p>
      </div>
      `
    ```
    
