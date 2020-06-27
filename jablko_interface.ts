// Jablko Interface Server
// Cale Overstreet
// May 27, 2020
// This file is the main entrypoint for the server that runs the web interface for Jablko. This server is response for handling user interactions and sending requests to the correct Jablko component.

console.log("Starting Jablko Interface...");
import { Application, Router, send } from "https://deno.land/x/oak/mod.ts";
import { readFileStr } from "https://deno.land/std/fs/mod.ts";

const app = new Application();
const router = new Router();

// Error listener for Oak Server
app.addEventListener("error", (evt) => {
  // Will log the thrown error to the console.
  console.log(evt.error);
});

// Important Paths. Possible not needed
const web_root = "public_html";

export const server_start_time = new Date().getTime();

// Initialize SMTP system
console.log("Initializing SMTP system...")
const smtp_module = await import("./src/smtp_module.ts");
export const smtp_client = await smtp_module.Jablko_Smtp_Initialize();

async function load_jablko_modules() {
	// Creates an object containing jablko modules and all exported functions for server routing
	var loaded_modules: any = new Object();

	console.log("Loading Jablko Modules...");
	for await (const dirEntry of Deno.readDir("./jablko_modules")) {
		loaded_modules[dirEntry.name] = await import(`./jablko_modules/${dirEntry.name}/${dirEntry.name}.ts`);
	}

	return loaded_modules;
}

var jablko_modules: any = await load_jablko_modules(); // Only bit that needs to use type any. Hopefully a future design removes this need


console.log("Creating Middleware Handlers...");

// Timer Middleware
app.use((await import("./src/timing.ts")).timing_middleware);

// User authentication middleware. 
app.use((await import("./src/user_authentication.ts")).check_authentication);

// Defining Server Routes
console.log("Defining Server Routes...");

router.get("/", async (context) => {
	var dashboard_string: string = await readFileStr(`${web_root}/dashboard/dashboard_template.html`);
	var module_string = "";

	// Read in toolbar and string replace into dashboard_string
	const toolbar_string = await readFileStr(`${web_root}/toolbar/toolbar.html`);
	dashboard_string = await dashboard_string.replace("$TOOLBAR", toolbar_string);


	// Go through all modules and generate module string
	for (var module_name in jablko_modules) {
		module_string += jablko_modules[module_name]["generate_card"]();
	}

	// Replace placeholder in template file
	dashboard_string = await dashboard_string.replace("$JABLKO_MODULES", module_string);

	// Set response
	context.response.type = "html";
	context.response.body = dashboard_string;
});

// Routes requests sent by client to correct jablko module
router.post('/jablko_modules/:module_name/:function_name', async (context) => {
	if (context.params.module_name !== undefined && context.params.function_name !== undefined) {
		await jablko_modules[context.params.module_name][context.params.function_name](context);
	}
});

// Adding router to middleware
app.use(router.routes());
app.use(router.allowedMethods());

// Static Content Route. Must pass through authentication and router to get to this point
app.use(async (context) => {
	await send(context, context.request.url.pathname, {
		root: web_root
	});
});

console.log("Jablko Interface Listening on Port 10230");
await app.listen({port: 10230});
