// bot_brain.js: GroupMe Bot message callback handler
// Cale Overstreet
// August 23, 2020
// Describes bot behavior and decision making
// Exports: async handle_message(req, res)

const fs = require("fs").promises;
const readFileSync = require("fs").readFileSync;
const { exec } = require("child_process");

const jablko = require("../jablko_interface.js");
var jablko_config = jablko.jablko_config;
const jablko_modules = jablko.jablko_modules;


const dictionary_path = "src/dictionary.csv";

function parse_dictionary(filename) {
	const data  = readFileSync(filename, "utf8");
	const split_data = data.split("\n");

	var new_dictionary = {}
	for (var i = 0; i < split_data.length; i++) {
		if (split_data[i] != "" && !split_data[i].startsWith("Word")) {
			const line_split = split_data[i].split(",");
			var word_vector = []
			for (var j = 1; j < line_split.length; j++) {
				word_vector.push(parseInt(line_split[j]));
			}

			new_dictionary[line_split[0]] = word_vector;
		}
	}

	return new_dictionary;
}

const dictionary = parse_dictionary(dictionary_path);

// -------------------- START Module chatbot exports --------------------
// Use chat functions described in jablko_config.json

var jablko_module_functions = {};

for (module_name in jablko_config.jablko_modules) {
	jablko_module_functions[module_name] = [];
	
	for (func in jablko_config.jablko_modules[module_name].chatbot) {
		jablko_module_functions[module_name].push(jablko_config.jablko_modules[module_name].chatbot[func]);
	}
}

(async () => {
	// Load modules and initialize intent array
	for (module_name in jablko_module_functions) {
		for (var i = 0; i < jablko_module_functions[module_name].length; i++) {
			jablko_module_functions[module_name][i].activation = await determine_intent(jablko_module_functions[module_name][i].activation_phrase);
		}
	}

	console.debug(jablko_module_functions);
})();


// -------------------- END Module chatbot exports --------------------

function add_vector(v1, v2) {
	const new_vector = []
	for (var i = 0; i < v1.length; i++) {
		new_vector.push(v1[i] + v2[i]);
	}

	return new_vector;
}

async function determine_intent(phrase) {
	phrase = phrase.toLowerCase();
	phrase = phrase.replace(/[^\w\s]/g, "");
	const split_phrase = phrase.split(" ");
	var intent_vector = []

	for (var i = 0; i < dictionary["hi"].length; i++) {
		intent_vector.push(0);
	}

	for (var i = 0; i < split_phrase.length; i++) {
		if (split_phrase[i] != "") {
			if (split_phrase[i] in dictionary) {
				intent_vector = add_vector(intent_vector, dictionary[split_phrase[i]]);
			}
		}
	}

	return intent_vector;
}

async function create_response(message) {
	// Create intent vector
	const message_intent = await determine_intent(message);

	// Create required action list
	var action_list = []

	for (module_name in jablko_module_functions) {
		for (var i = 0; i < jablko_module_functions[module_name].length; i++) {

			var should_continue = false;
			var abs_diff = 0

			for (var j = 0; j < jablko_module_functions[module_name][i].activation.length; j++) {
				abs_diff += Math.abs(message_intent[j] - jablko_module_functions[module_name][i].activation[j]);
				//console.log(`${message_intent[j]}, ${jablko_module_functions[module_name][i].activation[j]}`)
				if (message_intent[j] < jablko_module_functions[module_name][i].activation[j]) {
					should_continue = true;
					break;
				}
			}	

			if (abs_diff < 15 && should_continue == false) {
				action_list.push({module_name: module_name, function: jablko_module_functions[module_name][i].function});
			}
		}
	}

	console.debug(action_list);

	// Create response and call appropriate functions
	var response = "";

	for (var i = 0; i < action_list.length; i++) {
		response += await jablko_modules[action_list[i].module_name][action_list[i].function]() + " ";
	}

	const confused_responses = [
		"Not quite sure what you said there.",
		"I don't know what you mean.",
		"I don't think I know enough to respond.",
		"Maybe you should add some words to my dictionary?",
		"Nani?",
		"#dontunderstand",
		"Well tak."
	];

	if (action_list.length == 0) {
		exec(`bash -c "echo ${message} >> log/unknown_phrases.txt"`);
		return confused_responses[Math.floor(Math.random() * 100) % confused_responses.length];
	}

	return response;
}

module.exports.parse_message = async (message) => {
	return await create_response(message);	
}
