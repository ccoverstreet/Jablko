// Jarmuz RGB Light ESP8266 UDP Server
// Cale Overstreet
// Jun. 26, 2021

// Used for the Jarmuz-RGB-Light JMOD

#include <ESP8266WiFi.h>
#include <WiFiUdp.h>

#include "config.h"

#define UDP_PORT 4123
WiFiUDP Udp;
char g_inBuffer[64];
const char *reply = "Received";

#define RED_PIN 13
#define BLUE_PIN 16
#define GREEN_PIN 12


void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);

  WiFi.begin(WIFI_SSID, WIFI_PASSWORD);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("Connected");
  Serial.println(WiFi.localIP());

  Udp.begin(UDP_PORT);
  Serial.println("Started UDP Server");
  analogWrite(RED_PIN, 512);
}


int lastActivity = 0;

void loop() {
  // put your main code here, to run repeatedly:
  int packetSize = Udp.parsePacket();
  Serial.printf("%d bytes received\n", packetSize);
  if (packetSize == 4) {
    lastActivity = millis();
    int nByteRead = Udp.read(g_inBuffer, 64);
    if (nByteRead > 0) {
      g_inBuffer[nByteRead] = 0;
    }

    int level = g_inBuffer[3] / 255.0 * 1023;
    int red = g_inBuffer[0] / 255.0 * level;
    int green = g_inBuffer[1] / 255.0 * level;
    int blue = g_inBuffer[2] / 255.0 * level;
    
    //Serial.printf("%d %d %d %d\n", red, green, blue, level);
    
    analogWrite(RED_PIN, red);
    analogWrite(GREEN_PIN, green);
    analogWrite(BLUE_PIN, blue);
  }
  
  if (millis() - lastActivity > 10000) {
    delay(2000);
    return;
  }
  
  delay(90);
}
