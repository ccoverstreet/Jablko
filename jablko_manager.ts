// jablko.ts: Jablko Runtime Manager
// Cale Overstreet
// August 9, 2020
// This script creates the interface worker and restarts it when it either fails or recieves a reset command

async function jablko_log(prefix: any, input: any) {
	const prefix_length = prefix.length + 2;
	const line_split = await input.split("\n");
	var output = "";
	output += `${prefix}: ${line_split[0]}`;
	for (var i = 1; i < line_split.length; i++) {
		output += "\n" + " ".repeat(prefix_length) + `${line_split[i]}`;
	}
	console.log(output);
}

var jablko_interface: any = undefined;

function create_jablko_interface() {
	jablko_interface = new Worker(new URL("jablko_interface.ts", import.meta.url).href, {type: "module", deno: true});
	jablko_interface.addEventListener("message", function(event: any) {
		if (typeof(event.data) == "string") {
			if (event.data == "restart") {
				jablko_interface.terminate();
				Deno.exit(240);
			}
				jablko_log(`Interface [${new Date().toLocaleString("en-CA", {timeZone: "America/New_York"})}]`, event.data.toString());
		} else {
			console.log(event.data)	;
		}
	});
}

create_jablko_interface();
