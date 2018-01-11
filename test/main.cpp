#include "Cosa/RTT.hh"
#include "Cosa/OutputPin.hh"
#include "stl/ArrayList.h"

wlp::ArrayList<char> char_list(15);
OutputPin ledPin(Board::LED);

void setup() {
    RTT::begin();
}

void loop() {
    ledPin.on();
    delay(50);
    ledPin.off();
    delay(500);
    char_list.push_back('a');
    if (char_list.size() == char_list.capacity()) {
        char_list.clear();
    }
}

