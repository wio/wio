#ifndef COSAEXAMPLE_EXLIB_H
#define COSAEXAMPLE_EXLIB_H

#define EXLIB_INCLUDED true

#include <stdint.h>

static constexpr uint16_t get_num() {
    return 1;
}

struct DelayGenerator {
private:
    int m_delay;

public:
    explicit DelayGenerator(int delay);

    int get_delay() const;
};

#endif //COSAEXAMPLE_EXLIB_H
