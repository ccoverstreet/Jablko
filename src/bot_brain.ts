// bot_brain.ts: GroupMe Bot message callback handler
// Cale Overstreet
// July 24, 2020
// This file describes the behavior of the GroupMe bot. In the future it could use language comprehension to determine actions, but I'll just use some general rules for now.
// Exports: handle_message

import { readFileStr } from "https://deno.land/std/fs/mod.ts";
import { exec } from "https://deno.land/x/exec/mod.ts";

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
	const available_responses = [
		"Hi!",
		"Hello!",
		"What's up my homie?",
		"Dzien Dobry!",
		"Top of the morning to ya!",
		"How's it hanging?"
	];

	return available_responses[Math.floor(Math.random() * 100) % available_responses.length];
}

function test_fun() {
	messaging_system.send_message("Testing parser");
	return "Ran test function.";
}

function rude() {
	const available_responses = [
		"Well ok then.",
		"Your mom.",
		"Do I sense a small pp?",
		"Hoe.",
		"Hoe-bitch."
	];

	return available_responses[Math.floor(Math.random() * 100) % available_responses.length];
}

function get_weather() {
	const available_responses = [
		"I'm getting the weather for you!",
		"Here's the weather.",
		"Coming right up."
	];


	return available_responses[Math.floor(Math.random() * 100) % available_responses.length];
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
	},
	{
		name: "Rude",
		activation_phrase: "Stupid",
		function: rude
	},
	{
		name: "Get Weather",
		activation_phrase: "What is the weather?",
		function: get_weather
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
	phrase = await phrase.replace(/[^\w\s]/g, "");
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

	return intent_vector;
}


async function parse_actions(actions: any) {
	for (var i = 0; i < actions.length; i++) {
		actions[i].activation = await determine_intent(actions[i].activation_phrase)		
		for (var j = 0; j < actions[i].activation.length; j++) {
			if (actions[i].activation[j] > 1) {
				actions[i].activation[j] = 1;
			}
		}
	}
}

await parse_actions(actions);

async function create_response(message: string) {
	// Create intent vector
	const message_intent = await determine_intent(message);

	// Create required action list
	var action_list = []
	for (var i = 0; i < actions.length; i++) {
		var should_continue = false;
		var abs_diff = 0
		for (var j = 0; j < message_intent.length; j++) {
			abs_diff += Math.abs(message_intent[j] - actions[i].activation[j]);
			if (message_intent[j] < actions[i].activation[j]) {
				should_continue = true;
				break;
			}
		}

		if (abs_diff < 15 && should_continue == false) {
			action_list.push(i);
		}
	}

	// Create response and call appropriate functions
	var response = "";

	for (var i = 0; i < action_list.length; i++) {
		response += await actions[action_list[i]].function() + " ";
	}

	const confused_responses = [
		"Not quite sure what you said there.",
		"I don't think I know enough to respond.",
		"Maybe you should add some words to my dictionary?",
		"Nani?",
		":/",
		"#dontunderstand",
		"Well tak."
	];

	if (action_list.length == 0) {
		exec(`bash -c "echo ${message} >> src/unknown_phrases.txt"`);
		return confused_responses[Math.floor(Math.random() * 100) % confused_responses.length];
	}

	return response;
}



export async function handle_message(context: any) {
	var message = await context.json_content.text.toLowerCase();

	if (message.includes("jablko")) {
		const generated_response = await create_response(context.json_content.text);
		messaging_system.send_message(generated_response);
	}

	context.response.type = "html"
	context.response.body = "";
}

