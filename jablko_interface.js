// jablko_interface.js: Primary entrypoint for Jablko Smart Home
// Cale Overstreet
// August 19, 2020
// Contains the setup required for the NodeJS Express Web Server and initializes all Jablko Modules
// Exports: jablko_config, html_root

const fs = require("fs").promises;

const express = require("express");
const app = express();

const port = 8080;

// Predefined config and paths (with exports)
const jablko_config = require("./jablko_config.json");
module.exports.jablko_config = jablko_config;
const html_root = "public_html";
module.exports.html_root = html_root;

console.log(jablko_config);

async function jablko_modules_load() {
	for (var i = 0; i < jablko_config.jablko_modules.length; i++) {
		console.log(jablko_config.jablko_modules[i]);
	}	
}

const jablko_modules = jablko_modules_load();

//	-------------------- Middleware --------------------

app.use(require("./src/timing.js").timing_middleware);

//	-------------------- End Middleware --------------------

app.get("/", async (req, res) => {
	const dashboard_template = await fs.readFile(`${html_root}/dashboard/dashboard_template.html`, "utf8");
	res.send(dashboard_template);
});

app.listen(port, () => {
	console.log(`Started Jablko Interface on Port 8080`);
});
