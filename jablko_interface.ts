// Jablko Interface Server
// Cale Overstreet
// May 27, 2020
// This file is the main entrypoint for the server that runs the web interface for Jablko. This server is response for handling user interactions and sending requests to the correct Jablko component.

console.log("Starting Jablko Interface...");
import { Application, Router, send } from "https://deno.land/x/oak/mod.ts";
import { readFileStr } from "https://deno.land/std/fs/mod.ts"

const app = new Application();
const router = new Router();

const web_root = "public_html";

console.log(window);

// Keep an average for a frame of 100 requests. Once 100 is reached, clear the average and reset the size counter
export const server_start_time = new Date().getTime();

export var request_handling_times = {
	current_average: 0,
	window_size: 100,
	size_counter: 0
}

async function load_jablko_modules() {
	var loaded_modules: any = new Object();

	for await (const dirEntry of Deno.readDir("jablko_modules")) {
		loaded_modules[dirEntry.name] = await import(`./jablko_modules/${dirEntry.name}/${dirEntry.name}.ts`);
	}

	return loaded_modules;
}

var jablko_modules: any = await load_jablko_modules(); // Only bit that needs to use type any. Hopefully a future design removes this need
console.log(jablko_modules);

console.log("Creating Middleware Handlers...");

app.addEventListener("error", (evt) => {
  // Will log the thrown error to the console.
  console.log(evt.error);
});

// Timer Middleware. Logs how much time it takes to handle a request
app.use(async (context, next) => {
	const request_start_time = new Date().getTime();	
	await next();
	
	if (request_handling_times.size_counter > request_handling_times.window_size) {
		request_handling_times.size_counter = 1;
		request_handling_times.current_average = 0;
	} else {
		request_handling_times.size_counter++;
	}
	

	const size = request_handling_times.size_counter;

	request_handling_times.current_average = request_handling_times.current_average * (size - 1) / size + (new Date().getTime() - request_start_time) / (size);
}); 

router.get("/", async (context) => {
	var dashboard_string: string = await readFileStr(`${web_root}/dashboard/dashboard_template.html`);
	var module_string = "";

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

router.post('/jablko_modules/:module_name/:function_name', (context) => {
	if (context.params.module_name !== undefined && context.params.function_name !== undefined) {
		jablko_modules[context.params.module_name][context.params.function_name](context);
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

await app.listen({port: 10230});
