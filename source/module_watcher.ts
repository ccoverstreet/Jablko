const module_watcher = Deno.watchFs("./jablko_modules");

for await (const event of module_watcher) {
	//self.postMessage("Modules Changed");
	if (event.kind == "modify") {
		// Parse out module name and pass to interface for module reloading
		var split_path = event.paths[0].split("/");
		self.postMessage(split_path[split_path.length - 2]);
	}
}

