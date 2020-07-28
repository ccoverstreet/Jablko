import { readFileStr } from "https://deno.land/std/fs/mod.ts";

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

function greeting() {
	return "Hello!";
}

function dwayne() {
	return "Dwayne. Yeet.";
}

function sentience() {
	return "I am learning.";
}

const dictionary = await parse_dictionary("./dictionary.csv");
console.log(dictionary);

const actions: any = [
	{
		name: "greeting",
		activation_phrase: "Hi",
		function: greeting
	},
	{
		name: "dwayne",
		activation_phrase: "Dwayne is the singularity",
		function: dwayne
	},
	{
		name: "sentience check",
		activation_phrase: "Are you sentient",
		function: sentience
	}
]


function add_vector(v1: any, v2: any) {
	const new_vector = []
	for (var i = 0; i < v1.length; i++) {
		new_vector.push(v1[i] + v2[i]);
	}

	return new_vector;
}

async function determine_intent(phrase: string) {
	phrase = await phrase.toLowerCase();
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

async function parse_activation_phrases(actions: any) {
	for (var i = 0; i < actions.length; i++) {
		actions[i].activation = await determine_intent(actions[i].activation_phrase);
	}
}

async function determine_actions(message: string) {
	message = await message.toLowerCase();
	message = await message.replace(/[,.?:;!$#@]/, "");
	const split_message = message.split(" ");

	var message_intent = []

	for (var i = 0; i < dictionary["hi"].length; i++) {
		message_intent.push(0);
	} 

	for (var i = 0; i < split_message.length; i++) {
		if (split_message[i] != "" && dictionary[split_message[i]] != undefined) {
			message_intent = add_vector(message_intent, dictionary[split_message[i]]);
		}
	}

	var action_list: any = []
		
	for (var i = 0; i < actions.length; i++) {
		var should_continue = false;

		var abs_diff = 0;

		for (var j = 0; j < message_intent.length; j++) {
			abs_diff += Math.abs(message_intent[j] - actions[i].activation[j]);
			if (message_intent[j] < actions[i].activation[j]) {
				should_continue = true;
				break;
			}
		}

		if (should_continue) {
			continue;
		} else if (abs_diff >= 2) {
			continue;
		}

		action_list.push(i);
	}

	var response = ""

	for (var i = 0; i < action_list.length; i++) {
		response += actions[action_list[i]].function() + " ";
	}

	console.log(message);
	console.log(message_intent);
	console.log(response);
	console.log();
}

await parse_activation_phrases(actions);
console.log(actions);

await determine_actions("Dwayne is the singularity");
await determine_actions("Are you sentient?");
await determine_actions("Hi, turn on the light");
