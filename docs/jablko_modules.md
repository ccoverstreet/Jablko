# Jablko Modules

Jablko Modules can be made by creating the required typescript file with external functions thay the main interface can call. Below is a list of items that **MUST** be defined in a module.

## Contents

- [Overview](#overview)
- [Info](#info)
- [Generate Card](#generate-card)
- [Card Design](#card-design)

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
## Info

The info object is exported and provides the main interface or other modules information on the security level and other parameters in the future. The only critical component for now is the `permission_level` and must be included for a module to work. Permission levels: 0 all, 1 slightly elevated, 2 administrator.
```Javascript
export const info = {
  permission_level: 0
}
```

## Generate Card

This function is required if you want your module to display on the dashboard. This function must return a string that contains the html of the modules display card. If you want to look at available CSS to unify the appearance, look in the [Card Design](#card-design) section.
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

## Card Design

In Jablko there is a general CSS file that contains class definitions for certain module components. The goal of this is to unify the appearance of modules and make it easy to adjust the overall look of Jablko.

| CSS Selector | CSS Content|
| --- | --- |
| `.jablko_module_card` | display: inline-block;<br>margin: 10px;<br>border-radius: 5px;<br>width: calc(100% - 20px);<br>color: var(--font-color);<br>background-color: var(--color-card-background);|
| button | border-width: 0px;<br>border-radius: 5px;<br>padding: 5px;|

