#include <ESP8266WiFi.h>
#include "credentials.h"
#include <ESP8266HTTPClient.h>

const int windowSize = 30;
const int nWindows = 2880;

char dailyStorage[nWindows + 1] = {};
char currentWindow = '0';

int counter = 0;
int dailyCounter = 0;

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
  pinMode(12, INPUT);

  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.println("Connecting...");
  }
}

void sendStatus() {
  Serial.println(currentWindow);
  Serial.println(dailyStorage);
  
  HTTPClient http;
  http.begin("http://10.0.0.185:8080/local/hamstermonitor1/dump");
  http.addHeader("Content-Type", "application/json");

  char message[100];
  sprintf(message, "{\"Active\": %c}", currentWindow);
  Serial.println(message);
  int httpCode = http.POST(message); // "{\"Data\": [1, 2, 3]}"
  
}

void loop() {
  // put your main code here, to run repeatedly:
  Serial.printf(".");
  if (digitalRead(12)) {
    currentWindow = '1';
  }
 
  counter = counter + 1;
  
  if (counter == 30) {     
    counter = 0;
    dailyStorage[dailyCounter] = currentWindow;
    Serial.printf("\n");
    sendStatus();

    currentWindow = '0';
    dailyCounter = dailyCounter + 1;
    if (dailyCounter == nWindows) {
      dailyCounter = 0;
    }
  }
  delay(500);
}
