#include "exlib.h"

DelayGenerator::DelayGenerator(int delay)
        : m_delay(delay) {
}

int DelayGenerator::get_delay() const {
    return m_delay;
}
