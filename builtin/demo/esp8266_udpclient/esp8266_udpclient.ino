#include <ESP8266WiFi.h>
#include <WiFiUdp.h>

#define SSID "moviecouncil"
#define WIFI_PASS "RAAE!lectron11!"
#define UDP_PORT 41000 // Port on Jablko currently being used by test mod

// Globals
WiFiUDP udp;
char packet[128];

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);

  Serial.println("Connecting to Wifi");

  WiFi.begin(SSID, WIFI_PASS);
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
  }
  
  udp.begin(UDP_PORT);
}

void loop() {
  // put your main code here, to run repeatedly:

  udp.beginPacket("10.0.0.185", 41000);
  char t_buf[32];
  String(millis()).toCharArray(t_buf, sizeof(t_buf));
  udp.write(t_buf);
  udp.endPacket();

  delay(5000);

  // For handling incoming packets
  /*
  int packet_size = udp.parsePacket();
  if (packet_size) {
    int len = udp.read(packet, 128);
    packet[len] = '\0';
    Serial.println(packet);

    // Send response
    udp.beginPacket(udp.remoteIP(), udp.remotePort());
    udp.write(reply);
    udp.endPacket();
  }
  */
}
