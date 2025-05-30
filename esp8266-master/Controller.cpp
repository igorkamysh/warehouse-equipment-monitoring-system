#include "Controller.hpp"

using namespace packhouse;

void LED::on() const
{
  digitalWrite(m_pin, LOW);
}

void LED::off() const
{
  digitalWrite(m_pin, HIGH);
}

void Relay::on() const
{
  digitalWrite(m_pin, m_activeHigh ? HIGH : LOW);
}

void Relay::off() const
{
  digitalWrite(m_pin, m_activeHigh ? LOW : HIGH); 
}
