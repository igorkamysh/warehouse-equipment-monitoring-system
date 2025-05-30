#include <ArduinoJson.h>
#include <esp_wifi.h>

#include <sstream>

#include "WiFi.h"
#include "WebServer.h"

#include "Controller.hpp"

#define GREEN_LED 2
#define RELAY_PIN 26


static const char* kWifiLogin = "sw_gastr";
// static const char* kWifiPassw = "";

static const char* kServerAddr = "10.30.64.19";
static const char* kServerPath = "/register_machine";

static const uint16_t kServerPort = 8080;

unsigned long g_previousMillis    = 0;
unsigned long g_reconnectInterval = 30000;

std::string             g_chipID;
String                  g_ipAddr;
WebServer               g_server{};
StaticJsonDocument<300> g_jsonDoc;

packhouse::LED          g_greenLed{LED_BUILTIN};
packhouse::LED          g_redLed{LED_BUILTIN};
packhouse::Relay        g_relay{RELAY_PIN, HIGH};

packhouse::Controller<packhouse::Switchable> g_controller{g_greenLed, g_redLed, g_relay};

void readMacAddress(){
  uint8_t baseMac[6];
  esp_err_t ret = esp_wifi_get_mac(WIFI_IF_STA, baseMac);
  if (ret == ESP_OK) {
    Serial.printf("%02x:%02x:%02x:%02x:%02x:%02x\n",
                  baseMac[0], baseMac[1], baseMac[2],
                  baseMac[3], baseMac[4], baseMac[5]);
  } else {
    Serial.println("Failed to read MAC address");
  }
}

int getChipId()
{
  int chipId = 0;
  for (int i = 0; i < 17; i = i + 8) {
    chipId |= ((ESP.getEfuseMac() >> (40 - i)) & 0xff) << i;
  }
  return chipId;
}

void handleWebHook()
{
  String data = g_server.arg("plain");

  DeserializationError error = deserializeJson(g_jsonDoc, data);
  if (error)
  {
    Serial.print("Deserialization Error: ");
    Serial.println(error.c_str());
    return;
  }
  else
  {
    g_jsonDoc.containsKey("current_state");
    uint8_t state = g_jsonDoc["current_state"];
    if (state == static_cast<uint8_t>(packhouse::State::FREE))
    {
      // Unlock the machine
      g_controller.lock();
    }
    else if (state == static_cast<uint8_t>(packhouse::State::IN_USE))
    {
      // Lock the machine
      g_controller.unlock();
    }
    else
    {
      g_server.send(400, "text/plain", "Wrong status");
      return;
    }

    g_server.send(200, "text/plain", "success");
  }
}


void getBSSIDWebHook()
{
  WiFiClient client;
  StaticJsonDocument<100> postDoc;

  postDoc["current_bssid"] = WiFi.BSSIDstr();

  String jsonToSend;
  serializeJson(postDoc, jsonToSend);

  Serial.println(jsonToSend);

  if (client.connect(kServerAddr, kServerPort))
  {
    client.println("POST "  + String(kServerPath) + " HTTP/1.1");
    client.println("Host: " + String(kServerAddr));
    client.println("Content-Type: application/json");
    client.println("Content-Length: " + String(jsonToSend.length()));
    client.println();
    client.println(jsonToSend);

    client.stop(); // Close the connection
  }
  else
  {
    Serial.println("Connection to the server failed");
  }
}


void setup()
{
  // Set up pins
  pinMode(GREEN_LED, OUTPUT);
  // g_greenLed.off();
  pinMode(RELAY_PIN, OUTPUT);

  // Set up serial data
  Serial.begin(115200);
  Serial.println();

  // Esp chip ID from hex to string
  std::stringstream ss;
  ss << std::hex << getChipId();
  g_chipID = ss.str();

  WiFi.disconnect();
  // Connecting to Wi-Fi
  WiFi.mode(WIFI_STA);
  WiFi.begin("sw_gastr");


  // // Serial.print("Connecting");
  while (WiFi.status() != WL_CONNECTED)
  {
    delay(500);
    Serial.print(".");
    switch (WiFi.status()) {
      case WL_NO_SSID_AVAIL: Serial.println("[WiFi] SSID not found"); break;
      case WL_CONNECT_FAILED:
        Serial.print("[WiFi] Failed - WiFi not connected! Reason: ");
        break;
      case WL_CONNECTION_LOST: Serial.println("[WiFi] Connection was lost"); break;
      case WL_SCAN_COMPLETED:  Serial.println("[WiFi] Scan is completed"); break;
      case WL_DISCONNECTED:    Serial.println("[WiFi] WiFi is disconnected"); break;
      default:
        Serial.print("[WiFi] WiFi Status: ");
        Serial.println(WiFi.status());
        break;
    }
  }
  g_ipAddr = WiFi.localIP().toString();

  Serial.print("\nConnected, IP address: ");
  Serial.println(g_ipAddr.c_str());

  { ///< Send ip address to the server
    WiFiClient client;
    StaticJsonDocument<100> postDoc;

    postDoc["machine_id"] = g_chipID.c_str();
    postDoc["ip_addr"]    = g_ipAddr.c_str();

    String jsonToSend;
    serializeJson(postDoc, jsonToSend);

    Serial.println(jsonToSend);

    if (client.connect(kServerAddr, kServerPort))
    {
      // Serial.println("Connected to the server");

      client.println("POST "  + String(kServerPath) + " HTTP/1.1");
      client.println("Host: " + String(kServerAddr));
      client.println("Content-Type: application/json");
      client.println("Content-Length: " + String(jsonToSend.length()));
      client.println();
      client.println(jsonToSend);

      while (client.connected()) {
        String line = client.readStringUntil('\n'); ///< Read response from the server
        if (!line.isEmpty())
        {
          DeserializationError error = deserializeJson(g_jsonDoc, line);
          if (error)
          {
            // Serial.print("Deserialization Error: ");
            Serial.println(error.c_str());
          }
          else
          {
            uint8_t state = g_jsonDoc["current_state"];
            if (state == static_cast<uint8_t>(packhouse::State::FREE))
            {
              // Unlock the machine
              g_controller.lock();
            }
            else if (state == static_cast<uint8_t>(packhouse::State::IN_USE))
            {
              // Lock the machine
              g_controller.unlock();
            }
            else
            {
              continue;
            }
            client.stop(); // Close the connection
          }
        }
        
        // Serial.println(line);
      }
    }
    else
    {
      Serial.println("Connection to the server failed");
    }
  }

  Serial.print("Chip ID: ");
  Serial.println(g_chipID.c_str()); 

  // Create webhook '/chipID'
  std::string serverWebHook = "/" + g_chipID;
  std::string bssidWebHook = "/" + g_chipID + "/get_mac_addr";

  // Bind webhook to a function
  g_server.on(serverWebHook.c_str(), handleWebHook);
  g_server.on(bssidWebHook.c_str(), getBSSIDWebHook);

  // Start server
  g_server.begin();
  readMacAddress();
}


void loop(){
  g_server.handleClient();

  unsigned long currentMillis = millis();

  // if WiFi is down, try reconnecting
  if ((WiFi.status() != WL_CONNECTED) && (currentMillis - g_previousMillis >= g_reconnectInterval)) {
    // Serial.print(millis());
    Serial.println("Reconnecting to WiFi...");
    WiFi.disconnect();
    WiFi.reconnect();
    g_previousMillis = currentMillis;
  }
  //delay(500);
}

