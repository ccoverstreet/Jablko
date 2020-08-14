#include <ESP8266WebServer.h>
#include <ArduinoJson.h>
#include "config.h"

const char *ssid = WIFI_SSID;
const char *password = WIFI_PASSWORD;


ESP8266WebServer server;

#define RED_PIN 13
#define BLUE_PIN 12
#define GREEN_PIN 16

void setup() {
  Serial.begin(9600);
  delay(1);
  Serial.print("Starting up...\n");

  init_wifi();


  server.on("/", [](){server.send(200, "text/plain","RGB Non-Addressable Controller");});
  //server.on("/status", get_status);
  server.on("/set_rgb", set_rgb);
  server.begin();
  

  pinMode(RED_PIN, OUTPUT);
  pinMode(BLUE_PIN, OUTPUT);
  pinMode(GREEN_PIN, OUTPUT);
  pinMode(5, OUTPUT);
  digitalWrite(5, LOW);
  digitalWrite(13, HIGH);
  //analogWrite(13, 0);
  analogWrite(BLUE_PIN, 0);
  analogWrite(GREEN_PIN, 0);
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

unsigned ascending = 1;
void loop() {
  server.handleClient();

/*
  if (ascending) {
    for (int i = 0; i < 1024; i++) {
      analogWrite(RED_PIN, i);
      analogWrite(BLUE_PIN, i);
      analogWrite(GREEN_PIN, i);
      
      delay(1);
    }
    ascending = 0;
  } else {
    for (int i = 1023; i >= 0; i--) {
      analogWrite(RED_PIN, i);
      analogWrite(BLUE_PIN, i);
      analogWrite(GREEN_PIN, i);
      
      delay(1);
    }
    ascending = 1;
  }

 */
  delay(1);
}

void set_rgb() {
  String data = server.arg("plain");
  StaticJsonDocument<200> doc;

  DeserializationError error = deserializeJson(doc, data);
  if (error) {
    Serial.println("ERROR");
    return;
  }
  int red = doc["r"];
  Serial.println(red);
  analogWrite(RED_PIN, red * 4);
  server.send(200, "text/html", "Set Red Color");
}
