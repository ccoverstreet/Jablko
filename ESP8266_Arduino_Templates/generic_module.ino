#include <ESP8266WiFi.h>

const char *ssid = "NETWORK_SSID";
const char *password = "MYPASSWORD";

const unsigned int timeout = 2000; // Timeout time that kills client response if taking too long. prevents hanging or soft resets.


WiFiServer server(80); // Create and instance of WiFiServer open on Port 80.

void setup() {
	Serial.begin(9600);
	delay(1);
	Serial.print("Starting up...\n");
	
	init_wifi();
}

void init_wifi() {
	// Initializes WiFi on ESP8266. Must check for failure.
	WiFi.begin(ssid, password);
	Serial.print("Connecting to WiFi...\n");

	while (WiFi.status() != WL_CONNECTED) {
		Serial.print("ERROR: Unable to Connect to WiFi. Trying Again...\n");
		delay(5000); // delay for WiFi connection attempts
	}

	Serial.print("SUCCESS: Connected to WiFi.\n");

}

void loop() {
	WiFiClient client = server.available(); // Listen for clients

	unsigned long int current_time = millis();
	unsigned long int previous_time = current_time;


	if (client) {
		// If client connected
		current_time = millis();
		String request = "";


	}

}
