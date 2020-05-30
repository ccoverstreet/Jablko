// Create Worker Thread for interface
var jablko_interface = new Worker("./jablko_interface.ts", {type: "module", deno: true});

jablko_interface.onmessage = (message) => {
	console.log("Here is a message");
}

// Create Worker thread for sms server
