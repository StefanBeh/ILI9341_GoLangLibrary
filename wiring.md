# ILI9341 Display Wiring Instructions for Raspberry Pi

To connect the ILI9341 display to your Raspberry Pi, follow these instructions:

*   **VCC** of the display to a **3.3V** pin on the Raspberry Pi.
*   **GND** of the display to a **GND** pin on the Raspberry Pi.
*   **CS** (Chip Select) of the display to **GPIO 8 (SPI0 CE0)**.
*   **RST** (Reset) of the display to **GPIO 25**.
*   **DC** (Data/Command) of the display to **GPIO 24**.
*   **MOSI** (Master Out Slave In) of the display to **GPIO 10 (SPI0 MOSI)**.
*   **SCK** (Serial Clock) of the display to **GPIO 11 (SPI0 SCLK)**.
*   **LED** (Backlight) of the display to a **3.3V** pin.
*   **MISO** (Master In Slave Out) of the display to **GPIO 9 (SPI0 MISO)**.
