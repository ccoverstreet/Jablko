#include <ESP8266WebServer.h>
#include "config.h"

/*
There should be a config.h file in the same directory with the following contents
#DEFINE NETWORK_SSID "yournetworkname"
#DEFINE NETWORK_PASSWORD "yournetworkpassword"
 */

ESP8266WebServer server; // Create server holder

void setup() {
	Serial.begin(9600);
	Serial.print("Starting up...\n");
	init_wifi(); // Initialize wifi

	// Server definitions and start
	server.on("/", [](){ server.send(200, "text/plain", "Hello World from ESP8266"); }); // Custom route
	server.on("/status", status); // Custom route

	server.begin();

}

void init_wifi() {
	// Initializes WiFi on ESP8266. Must check for failure.
	WiFi.begin(NETWORK_SSID, NETWORK_PASSWORD);
	Serial.print("Connecting to WiFi...\n");

	while (WiFi.status() != WL_CONNECTED) {
		Serial.print("ERROR: Unable to Connect to WiFi. Trying Again...\n");
		delay(5000); // delay for WiFi connection attempts
	}

	Serial.print("SUCCESS: Connected to WiFi.\n");
	Serial.println(WiFi.localIP());
}

void status() {
	Serial.println("Good!");
}

void loop() {
	server.handleClient(); // Handle network requests
	delay(1);
}
