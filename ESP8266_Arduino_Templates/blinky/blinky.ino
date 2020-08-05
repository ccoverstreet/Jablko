#include <ESP8266WebServer.h>

const char *ssid = "SSID";
const char *password = "PASSWORD";

const unsigned int timeout = 2000; // Timeout time that kills client response if taking too long. prevents hanging or soft resets.


//WiFiServer server(80); // Create and instance of WiFiServer open on Port 80.

ESP8266WebServer server;

int light_status = 0;

void setup() {
	Serial.begin(9600);
	delay(1);
	Serial.print("Starting up...\n");
	
	init_wifi();

  server.on("/", [](){server.send(200, "text/plain","Hello World");});
  server.on("/status", get_status);
  server.on("/toggle_light", toggle_light);
  server.begin();

  pinMode(5, OUTPUT);
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
  Serial.println(WiFi.localIP());
}

void loop() {
  server.handleClient();
}

void get_status() {
  Serial.println("ONLINE");
  server.send(200, "text/plain", "GOOD");
}

void toggle_light() {
  if (light_status) {
    digitalWrite(5, LOW);
    light_status = 0;
  } else {
    digitalWrite(5, HIGH);
    light_status = 1;
  }

  server.send(200, "text/plain", "Toggled Light");
}
