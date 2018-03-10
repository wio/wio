#include "Cosa/RTT.hh"
#include "Cosa/OutputPin.hh"
#include "wlib/stl/ArrayList.h"
#include "module/module.h"

wlp::ArrayList<char> char_list(15);
OutputPin ledPin(Board::LED);

void setup() {
    RTT::begin();
    static_assert(FIND_MAX(15, 10) == 15, "FIND_MAX failed");
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

