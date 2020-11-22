// module.js: Jablko Food Inventory Module
// Cale Overstreet
// August 27, 2020
// Keeps track of food using an entry in the Jablko SQLite database in the module_name table

const fs = require("fs").promises;
const path = require("path");

// Base setup
const module_name = path.basename(__dirname)
const jablko = require(module.parent.filename)
const module_config = jablko.jablko_config.jablko_modules[module_name]

module.exports.permission_level = 0

module.exports.generate_card = async function() {
	return (await fs.readFile(`${__dirname}/food_inventory.html`, "utf8")).replace(/\$MODULE_NAME/g, module_name);
}
