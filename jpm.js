// jpm.js: Jablko Package Manager
// Cale Overstreet
// August 24, 2020
// Used to install Jablko Modules by first downloading module sources to the module_library directory and then copying contents to the correct directory in jablko_modules

const reader = require("readline-sync");
const extract = require("extract-zip");
const fetch = require("node-fetch");
const fs = require("fs");
const { execSync } = require("child_process");

const jablko_config = require("./jablko_config.json");

async function main() {
	console.log("Jablko Package Manager");
	if (process.argv[2] == "init") {
		await init();	
	}
}

async function init() {
	await execSync("mkdir -p jablko_modules");

	console.log("Installing all Jablko Modules specified in jablko_config.json");
	const module_keys = Object.keys(jablko_config.jablko_modules);
	for (module in jablko_config.jablko_modules) {
		console.log(`\t${module}`);
		await install_module(jablko_config.jablko_modules[module].repo_archive, module);
	}
}

async function install_module(repository_url, module_target_name) {
	const data = await fetch(repository_url);
	await execSync("mkdir -p module_library");
	await data.body.pipe(fs.createWriteStream(`./module_library/${module_target_name}.zip`));
	await extract(`./module_library/${module_target_name}.zip`, {dir: `${process.cwd()}/module_library`});
	console.log(repository_url);

	const split_repo_url = repository_url.split("/");
	const extracted_zip = `${split_repo_url[4]}-${split_repo_url[6].split(".zip")[0].replace("v", "")}`;

	await execSync(`mkdir -p ./jablko_modules/${module_target_name} && cp ./module_library/${extracted_zip}/* ./jablko_modules/${module_target_name}`);
}

main();
