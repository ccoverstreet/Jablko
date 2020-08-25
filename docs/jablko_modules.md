# Jablko Modules

Jablko Modules can be made by creating the required typescript file with external functions thay the main interface can call. Below is a list of items that **MUST** be defined in a module.

## Contents

- [Overview](#overview)
- [Generate Card](#generate-card)
- [Card Design](#card-design)

## Overview

Required components:
- "module.js": Must be located in the root of the Jablko Module's repository.
  - module.exports.permission_level: A exported int set to an integer value (0, 1, or 2). Permission level needed to use the module increase with the number. The level 2 denotes administrative level and 0 represents no required permissions.
    ```Javascript
    module.exports.permission_level = 0
    ```
  - module.exports.generate_card: If you want a card to appear on your dashboard you must provide an `async generate_card()`. function. This function is called when a user visits the server route and returns a string containing the HTML of the modules card. How you generate the string is up to each individual module. You can store the html in a separate .html file and then read in the html on load.
    ```Javascript
      module.exports.generate_card = async function generate_card() {
      return `
      <div id="mymodule_card" class="module_card">
        <div class="module_header">
                <h1>Interface Status</h1>
                <svg class="module_icon" viewBox="0 0 150 150">
                        <path d="M 20 75 H 40 L 60 30 L 90 120 L 110 75 H 130" fill="transparent" stroke="#0097e6" stroke-width="20px" stroke-linejoin="round" stroke-linecap="round"/>
                </svg>
        </div>
        <hr>
        <p>*Your Content*</p>
      </div>
      `
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

