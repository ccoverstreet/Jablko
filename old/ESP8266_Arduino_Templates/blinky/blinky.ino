#include <ESP8266WebServer.h>
#include "config.h"

const char *ssid = WIFI_SSID;
const char *password = WIFI_PASSWORD;

ESP8266WebServer server;

int light_status = 1;

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
	pinMode(13, OUTPUT);
  delay(500);
	analogWrite(13, 0);
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

unsigned int led_brightness = 0;
void loop() {
	server.handleClient();
  delay(1);
  if (light_status == 1 && led_brightness < 1024) {
    led_brightness += 4;
    analogWrite(13, led_brightness);
  } else if (light_status == 0 && led_brightness > 0) {
    led_brightness -= 4;
    analogWrite(13, led_brightness);
  }
}

void get_status() {
	Serial.println("ONLINE");
	server.send(200, "application/json", "{\"status\": \"good\", \"message\": \"Module on\"}");
}

void toggle_light() {
	if (light_status) {
		digitalWrite(5, LOW);
		light_status = 0;
	} else {
		digitalWrite(5, HIGH);
		light_status = 1;
	}

	server.send(200, "application/json", "{\"status\": \"good\", \"message\": \"toggled light\"}");
}
