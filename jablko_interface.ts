// Jablko Interface Server
// Cale Overstreet
// May 27, 2020
// This file is the main entrypoint for the server that runs the web interface for Jablko. This server is response for handling user interactions and sending requests to the correct Jablko component.
// Exports: server_start_time, jablko_config, messaging_system, jablko_modules

self.postMessage("Starting Jablko Interface...");

import { Application, Router, send } from "https://deno.land/x/oak@v6.0.1/mod.ts";
import { DB } from "https://deno.land/x/sqlite/mod.ts" ;
import { readFileStr } from "./src/util.ts";

const app = new Application();
const router = new Router();

// Error listener for Oak Server
app.addEventListener("error", (evt) => {
	// Will log the thrown error to the console. WHY IS THIS NOT DEFAULT?
	// Have to ignore SSL certificate errors as accessing from the same network prevents standard https protocol
	console.log(evt.error);
});

export const server_start_time = new Date().getTime(); // Used for measuring server uptime

// Important Paths. Possible not needed
const web_root = "public_html";

// Load configuration file and export it so modules have access to config data
self.postMessage("Loading \"jablko_config.json\"...")
export const jablko_config = await JSON.parse(await readFileStr("./jablko_config.json"));

// Load sqlite database and export connection
export const database = new DB("./database/primary.db");

// Initialize the messaging system
self.postMessage("Initializing messaging system (GroupMe Bot)...")
export const messaging_system = await import("./src/messaging.ts");

self.postMessage("Reading \"jablko_modules.config\"...");

async function load_jablko_modules() {
	// Creates an object containing jablko modules and all exported functions for server routing
	var loaded_modules: any = new Object();

	self.postMessage("Loading Jablko Modules...");
	for (var i = 0; i < jablko_config.jablko_modules.length; i++) {
		loaded_modules[jablko_config.jablko_modules[i]] = await import(`./jablko_modules/${jablko_config.jablko_modules[i]}/${jablko_config.jablko_modules[i]}.ts`);
	}

	return loaded_modules;
}

export const jablko_modules: any = await load_jablko_modules(); // Only bit that needs to use type any. Hopefully a future design removes this need

var module_list_output: string = "";
for (var name in jablko_modules) { // Print for startup info
	module_list_output += `\n\t${name}`;
}
self.postMessage(`Loaded Modules:${module_list_output}`);

self.postMessage("Loading Middleware...");

// Timer Middleware
app.use((await import("./src/timing.ts")).timing_middleware);

// User authentication middleware. 
app.use((await import("./src/user_authentication.ts")).check_authentication);

// Defining Server Routes
self.postMessage("Defining Server Routes...");

router.get("/", async (context) => {
	var dashboard_string: string = await readFileStr(`${web_root}/dashboard/dashboard_template.html`);
	var module_string = "";

	// Read in toolbar and string replace into dashboard_string
	const toolbar_string = await readFileStr(`${web_root}/toolbar/toolbar.html`);
	dashboard_string = await dashboard_string.replace("$TOOLBAR", toolbar_string);

	// Go through all modules and generate module string
	for (var module_name in jablko_modules) {
		module_string += await jablko_modules[module_name]["generate_card"]();
	}

	// Replace placeholder in template file
	dashboard_string = await dashboard_string.replace("$JABLKO_MODULES", module_string);

	// Set response
	context.response.type = "html";
	context.response.body = dashboard_string;
});

router.get("/restart", async (context: any) => {
	if (context.user_data.permission_level > 1) {
		messaging_system.send_message("Restart message received. Restarting in 5 seconds.");
		context.response.type = "json";
		context.response.body = {status: "good", message: "Restarting server"};
		self.postMessage("restart");
	}
});

// Routes requests sent by client to correct jablko module
router.post('/jablko_modules/:module_name/:function_name', async (context: any) => {
	if (context.user_data.permission_level < jablko_modules[context.params.module_name].permission_level()) {
		context.response.type = "json";
		context.response.body = {status: "fail", message: "Insufficient permissions"};
	} else if (context.params.module_name !== undefined && context.params.function_name !== undefined) {
		await jablko_modules[context.params.module_name][context.params.function_name](context);
	}
});

router.post("/bot_callback", (await import("./src/bot_brain.ts")).handle_message);

// Adding router to application
app.use(router.routes());
app.use(router.allowedMethods());

// Static Content Route. Must pass through authentication and router to get to this point
app.use(async (context) => {
	await send(context, context.request.url.pathname, {
		root: web_root
	});
});

self.postMessage("Jablko Interface Listening on Port 80 and 443");
app.listen({port: 443, secure: true, certFile: "../cert.pem", keyFile: "../privkey.pem"});
await app.listen({port: 80});
