package uart

import "github.com/ungerik/go-bbio"

func Setup(dt string) error {
	return bbio.LoadDeviceTree(dt)
}

func Cleanup() {
	bbio.UnloadDeviceTree("ADAFRUIT-UART1")
	bbio.UnloadDeviceTree("ADAFRUIT-UART2")
	bbio.UnloadDeviceTree("ADAFRUIT-UART3")
	bbio.UnloadDeviceTree("ADAFRUIT-UART4")
	bbio.UnloadDeviceTree("ADAFRUIT-UART5")
}
