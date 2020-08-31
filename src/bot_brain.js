// bot_brain.js: GroupMe Bot message callback handler
// Cale Overstreet
// August 23, 2020
// Describes bot behavior and decision making
// Exports: async handle_message(req, res)

const fs = require("fs").promises;
const readFileSync = require("fs").readFileSync;
const { exec } = require("child_process");

const jablko = require("../jablko_interface.js");
const jablko_config = jablko.jablko_config;

for (var i = 0; i < jablko_config.jablko_modules.length; i++) {
	if (jablko_config.jablko_modules[i].chatbot != undefined) {
		for (item of jablko_config.jablko_modules[i].chatbot) {
			console.log(item);
		}
	}
}

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

async function get_weather() {
	const available_responses = [
		"I'm getting the weather for you!",
		"Here's the weather.",
		"Coming right up."
	];

	const json_weather_data = await jablko.weather.get_current_weather();

	const temp_in_c = json_weather_data.current.temp - 273.15;
	const feels_in_c = json_weather_data.current.feels_like - 273.15;

	const weather_summary = `\nRight now it is ${(temp_in_c * 9/5 + 32).toFixed(2)} F (${temp_in_c.toFixed(2)} C) but feels like ${(feels_in_c * 9/5 + 32).toFixed(2)} (${feels_in_c.toFixed(2)} C).\nThe weather is "${json_weather_data.current.weather[0].description}" with a humidity of ${json_weather_data.current.humidity}%. \nThe wind is ${json_weather_data.current.wind_speed} m/s from ${json_weather_data.current.wind_deg} degrees from N."`;

	return available_responses[Math.floor(Math.random() * 100) % available_responses.length] + " " + weather_summary;
}

function uptime() {
	const available_responses = [
		"I've been alive for",
		"I've been running for",
		"I've been watching you for",
		"This version of me has been up for"
	]
	const raw_uptime = (new Date().getTime() - jablko.server_start_time) / 1000;	
	const hours = Math.floor(raw_uptime / 3600);
	const minutes = Math.floor((raw_uptime - 3600 * hours) / 60);
	const seconds = Math.floor(raw_uptime - 3600 * hours - 60 * minutes);
	const formatted_uptime = `${hours} h ${minutes} m ${seconds}s`;


	return available_responses[Math.floor(Math.random() * 100) % available_responses.length] + " " + formatted_uptime;
}

const actions = [
	{
		name: "Greeting",
		activation_phrases: [
			"Hi"
		],
		function: greeting
	},
	{
		name: "Test",
		activation_phrases: [
			"Send hello announcement"
		],
		function: test_fun
	},
	{
		name: "Rude",
		activation_phrases: [
			"Stupid"
		],
		function: rude
	},
	{
		name: "Get Weather",
		activation_phrases: [
			"weather"
		],
		function: get_weather
	},
	{
		name: "Uptime",
		activation_phrases: [
			"uptime",
			"running time"
		],
		function: uptime
	}
];
// -------------------- End of action definitions --------------------


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


async function parse_actions(actions) {
	for (var i = 0; i < actions.length; i++) {
		actions[i].activations = []
		for (var j = 0; j < actions[i].activation_phrases.length; j++) {
			actions[i].activations.push(await determine_intent(actions[i].activation_phrases[j]));
			for (var k = 0; k < actions[i].activations[j].length; k++) {
				if (actions[i].activations[j][k] > 1) {
					actions[i].activations[j][k] = 1;
				}
			}
		}
	}
}

parse_actions(actions);

async function create_response(message) {
	// Create intent vector
	const message_intent = await determine_intent(message);

	// Create required action list
	var action_list = []
	for (var i = 0; i < actions.length; i++) {
		for (var k = 0; k < actions[i].activations.length; k++) {

			var should_continue = false;
			var abs_diff = 0
			for (var j = 0; j < message_intent.length; j++) {
				abs_diff += Math.abs(message_intent[j] - actions[i].activations[k][j]);
				if (message_intent[j] < actions[i].activations[k][j]) {
					should_continue = true;
					break;
				}
			}

			if (abs_diff < 15 && should_continue == false) {
				action_list.push(i);
			}
		}

	}

	// Create response and call appropriate functions
	var response = "";

	for (var i = 0; i < action_list.length; i++) {
		response += await actions[action_list[i]].function() + " ";
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
		exec(`bash -c "echo ${message} >> src/unknown_phrases.txt"`);
		return confused_responses[Math.floor(Math.random() * 100) % confused_responses.length];
	}

	return response;
}

module.exports.handle_message = async (req, res) => {
	const message = await req.body.text.toLowerCase();

	if (message.includes("jablko")) {
		const generated_response = await create_response(message);
		jablko.messaging_system.send_message(generated_response);
	}

	res.send("Good");
}
