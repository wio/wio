generate_arduino_firmware(blink
        SRCS main.cpp
        ARDLIBS FAT16 ManchesterCodec exlib
        PORT /dev/ttyACM0
        BOARD mega2560)
