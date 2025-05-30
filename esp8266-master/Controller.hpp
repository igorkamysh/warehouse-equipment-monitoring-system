#include <sys/_stdint.h>
#pragma once

// #include "c_types.h"
#include <Arduino.h>


namespace packhouse
{


enum class State : uint8_t
{
  FREE   = 0,
  IN_USE = 1
};


class Switchable
{
public:
  ~Switchable() = default;

  virtual void on()  const = 0;
  virtual void off() const = 0;
};

class LED : public Switchable
{
private:
  uint8_t m_pin;

public:
  LED(uint8_t pin) : m_pin{pin}
  { }

  void on() const override;
  void off() const override;
};


class Relay : public Switchable
{
private:
  uint8_t m_pin;
  bool    m_activeHigh;

public:
  Relay(uint8_t pin, uint8_t activeHigh) : m_pin{pin}, m_activeHigh{activeHigh}
  { }

  void on() const override;
  void off() const override;
};


template <typename T>
class Controller
{
private:
  T & m_greenLed;
  T & m_redLed;
  T & m_relay;
  
public:
  Controller(T& greenLed, T& redLed, T& relay) :
            m_greenLed{greenLed},
            m_redLed{redLed},
            m_relay{relay}
  { }

  void unlock() 
  {
    m_relay.on();
    m_greenLed.on();
    // TODO: Logic with red led
  }
  
  void lock()
  {
    m_relay.off();
    m_greenLed.off();
    // TODO: Logic with red led
  }
};
 
} // packhouse