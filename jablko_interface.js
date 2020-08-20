// jablko_interface.js: Primary entrypoint for Jablko Smart Home
// Cale Overstreet
// August 19, 2020
// Contains the setup required for the NodeJS Express Web Server and initializes all Jablko Modules
// Exports: jablko_config, html_root

async function main() {
	// Overriding console.log
	const old_console_log = console.log;
	console.log = function (input) {
		old_console_log(`[${new Date().toLocaleString("sv-SE")}]:`, input);
	}

	// -------------------- START Package Requires --------------------

	const fs = require("fs").promises;
	const readFileSync = require("fs").readFileSync;
	const http = require("http");
	const https = require("https");

	const express = require("express");
	const app = express();

	const sqlite = require("sqlite-async");

	// -------------------- END Package Requires --------------------

	// -------------------- START Module Initialization --------------------

	// Predefined config and paths (with exports)
	const jablko_config = require("./jablko_config.json");
	module.exports.jablko_config = jablko_config;
	const html_root = `${__dirname}/public_html`;
	module.exports.html_root = html_root;
	module.exports.server_start_time = Date.now();

	console.log(jablko_config);

	const user_db = await sqlite.open(jablko_config.database.path)
	module.exports.user_db = user_db

	// Load and export jablko_modules
	function jablko_modules_load() {
		var loaded_modules = {};

		for (var i = 0; i < jablko_config.jablko_modules.length; i++) {
			loaded_modules[jablko_config.jablko_modules[i]] = require(`./jablko_modules/${jablko_config.jablko_modules[i]}/${jablko_config.jablko_modules[i]}.js`);
		}	

		return loaded_modules;
	} 

	const jablko_modules = jablko_modules_load();

	// -------------------- END Module Initialization --------------------

	// -------------------- START Middleware --------------------

	app.use(require("./src/timing.js").timing_middleware);
	app.use(require("cookie-parser")())
	app.use(express.json());
	app.use(require("./src/user_authentication.js").user_authentication_middleware);

	// -------------------- END Middleware --------------------

	// -------------------- START End Routes --------------------

	app.get("/", async (req, res) => {
		var dashboard_template = await fs.readFile(`${html_root}/dashboard/dashboard_template.html`, "utf8");

		// Load Jablko Module Cards
		var module_string = "";
		for (var i = 0; i < jablko_config.jablko_modules.length; i++) {
			module_string += await jablko_modules[jablko_config.jablko_modules[i]].generate_card();
		}

		dashboard_template = dashboard_template.replace("$JABLKO_MODULES", module_string);
		dashboard_template = dashboard_template.replace("$TOOLBAR", await fs.readFile("./public_html/toolbar/toolbar.html"));
		dashboard_template += "<style>" + await fs.readFile("./public_html/dashboard/dashboard.css") + "</style>";

		res.send(dashboard_template);
	});

	app.post("/bot_callback", async (req, res) => {
		console.log("Bot Callback");
	});

	app.post("/jablko_modules/:module_name/:handler", async (req, res) => {
		await jablko_modules[req.params.module_name][req.params.handler](req, res);
	});

	// -------------------- END End Routes --------------------

	// -------------------- START Server Start --------------------
	// Check from config for HTTP/HTTPS configuration

	var http_server = undefined;
	if (jablko_config.http.port != null) {
		http_server = http.createServer(app);
		http_server.listen(jablko_config.http.port, () => {
			console.log(`Started Jablko Interface on Port ${jablko_config.http.port}`);
		});
	}

	var https_server = undefined;
	if (jablko_config.https.port != null) {
		console.log("IMPLEMENT THIS");
		app.listen(jablko_config.https.port, () => {
			console.log(`Started Jablko Interface on Port ${jablko_config.https.port}`);
		});
	}

	// -------------------- END Server Start --------------------
}

main();
