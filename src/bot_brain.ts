// bot_brain.ts: GroupMe Bot message callback handler
// Cale Overstreet
// July 24, 2020
// This file describes the behavior of the GroupMe bot. In the future it could use language comprehension to determine actions, but I'll just use some general rules for now.
// Exports: handle_message

import { readFileStr } from "https://deno.land/std/fs/mod.ts";

const dictionary_path = "src/dictionary.csv";

// Parses designed dictionary file creates a dictionary object
async function parse_dictionary(filename: string) {
	const data  = await readFileStr(filename);
	const split_data = await data.split("\n");

	var new_dictionary: any = {}
	for (var i = 0; i < split_data.length; i++) {
		if (split_data[i] != "" && !split_data[i].startsWith("Word")) {
			const line_split = await split_data[i].split(",");
			var word_vector = []
			for (var j = 1; j < line_split.length; j++) {
				word_vector.push(parseInt(line_split[j]));
			}
			
			new_dictionary[line_split[0]] = word_vector;
		}
	}


	return new_dictionary;
}

const dictionary = await parse_dictionary(dictionary_path);
const messaging_system = (await import("../jablko_interface.ts")).messaging_system;

// -------------------- Begining of action definitions -------------------- 
function greeting() {
	return "Hello";
}

function test_fun() {
	messaging_system.send_message("Testing parser");
	return "Ran test function";
}

const actions: any = [
	{
		name: "Greeting",
		activation_phrase: "Hi",
		function: greeting
	},
	{
		name: "Test",
		activation_phrase: "Send hello announcement",
		function: test_fun
	}
];
// -------------------- End of action definitions --------------------

function add_vector(v1: any, v2: any) {
	const new_vector = []
	for (var i = 0; i < v1.length; i++) {
		new_vector.push(v1[i] + v2[i]);
	}

	return new_vector;
}

async function determine_intent(phrase: string) {
	phrase = await phrase.toLowerCase();
	phrase = await phrase.replace(/[,.?:;!$#@]/, "");
	const split_phrase = await phrase.split(" ");
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

	return intent_vector
}


async function parse_actions(actions: any) {
	for (var i = 0; i < actions.length; i++) {
		actions[i].activation = await determine_intent(actions[i].activation_phrase)		
	}
}

await parse_actions(actions);

console.log(actions)

async function create_response(message: string) {
	
}


export async function handle_message(context: any) {
	context.json_content.text = context.json_content.text.toLowerCase();
	console.log(context.json_content.text);

	context.response.type = "html"
	context.response.body = "";
}

