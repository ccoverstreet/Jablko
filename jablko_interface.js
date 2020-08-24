// jablko_interface.js: Primary entrypoint for Jablko Smart Home
// Cale Overstreet
// August 19, 2020
// Contains the setup required for the NodeJS Express Web Server and initializes all Jablko Modules
// Exports: jablko_config, html_root


// Overriding console.log
const old_console_log = console.log;
console.log = function(input) {
	old_console_log(`[${new Date().toLocaleString("sv-SE")}]:`, input);
}

console.debug = function(input) {
	if (DEBUG_ON) {
		console.log(input);
	}
}

const DEBUG_ON = (process.argv[2] == "--debug" || process.argv[2] == "-d") ? true : false;

if (DEBUG_ON) {
	console.log("Starting Jablko in DEBUG Mode...\n");
} else {
	console.log("Starting Jablko...\n");
}

async function main() {
	// Primary server calls and initialization. 
	
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
	console.log(`Loading Jablko Config from "jablko_config.json..."`);
	const jablko_config = require("./jablko_config.json");
	module.exports.jablko_config = jablko_config;
	console.debug(jablko_config);

	console.log("Setting Server Information Constants...");
	const html_root = `${__dirname}/public_html`;
	module.exports.html_root = html_root;
	module.exports.server_start_time = Date.now();

	console.log("Intializing GroupMe messaging system...");
	module.exports.messaging_system = require("./src/messaging.js");

	console.log("Opening SQLite database...");
	module.exports.user_db = await sqlite.open(jablko_config.database.path);

	console.log(`Loading OWM weather wrapper from "src/weather.js"...`);
	module.exports.weather = require("./src/weather.js");

	console.log("Loading Jablko Modules...");
	// Load and export jablko_modules
	function jablko_modules_load() {
		var loaded_modules = {};
		const modules = Object.keys(jablko_config.jablko_modules);

		for (var i = 0; i < modules.length; i++) {
			loaded_modules[modules[i]] = require(`./jablko_modules/${modules[i]}/module.js`);
		}	

		console.debug(loaded_modules);
		return loaded_modules;
	} 

	const jablko_modules = jablko_modules_load();

	// -------------------- END Module Initialization --------------------

	// -------------------- START Middleware --------------------

	console.log("Loading middleware...");
	app.use(require("./src/timing.js").timing_middleware);
	app.use(require("cookie-parser")())
	app.use(express.json());
	app.use(require("./src/user_authentication.js").user_authentication_middleware);

	// -------------------- END Middleware --------------------

	// -------------------- START End Routes --------------------

	console.log("Establishing routes...");

	app.get("/", async (req, res) => {
		var dashboard_template = await fs.readFile(`${html_root}/dashboard/dashboard_template.html`, "utf8");

		// Load Jablko Module Cards
		var module_string = "";
		const modules = Object.keys(jablko_config.jablko_modules);
		for (var i = 0; i < modules.length; i++) {
			module_string += await jablko_modules[modules[i]].generate_card();
		}

		dashboard_template = dashboard_template.replace("$JABLKO_MODULES", module_string);
		dashboard_template = dashboard_template.replace("$TOOLBAR", await fs.readFile("./public_html/toolbar/toolbar.html"));
		dashboard_template += "<style>" + await fs.readFile("./public_html/dashboard/dashboard.css") + "</style>";

		res.send(dashboard_template);
	});

	const bot_brain = require("./src/bot_brain.js");
	app.post("/bot_callback", async (req, res) => {
		console.log("Message received from GroupMe");
		await bot_brain.handle_message(req, res);
	});

	app.get("/restart", async (req, res) => { 
		process.exit(240);
	});

	app.post("/jablko_modules/:module_name/:handler", async (req, res) => {
		await jablko_modules[req.params.module_name][req.params.handler](req, res)
			.catch((error) => {
				console.log(`Requested module path "${req.params.module_name}/${req.params.handler} not found"`);
				console.log(error);
				res.json({status: "fail", message: "Module path not found"});
			});
	});

	// -------------------- END End Routes --------------------

	// -------------------- START Server Start --------------------
	// Check from config for HTTP/HTTPS configuration

	var http_server = undefined;
	if (jablko_config.http.port != null) {
		http_server = http.createServer(app);

		http_server.listen(jablko_config.http.port, () => {
			console.log(`Started Jablko Interface on Port ${jablko_config.http.port} (HTTP)`);
		});
	}

	var https_server = undefined;
	if (jablko_config.https.port != null) {
		// Read PEM files
		https_server = https.createServer({
			key: await fs.readFile(jablko_config.https.key_file, "utf8"),
			cert: await fs.readFile(jablko_config.https.cert_file, "utf8")
		}, app);

		https_server.listen(jablko_config.https.port, () => {
			console.log(`Started Jablko Interface on Port ${jablko_config.https.port} (HTTPS)`);
		});
	}

	// -------------------- END Server Start --------------------
}

main()
	.catch((error) => {
		console.log(error);
		return;
	});
