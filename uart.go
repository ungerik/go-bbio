package bbio

func UARTSetup(dt string) error {
	return LoadDeviceTree(dt)
}

func UARTCleanup() {
	UnloadDeviceTree("ADAFRUIT-UART1")
	UnloadDeviceTree("ADAFRUIT-UART2")
	UnloadDeviceTree("ADAFRUIT-UART3")
	UnloadDeviceTree("ADAFRUIT-UART4")
	UnloadDeviceTree("ADAFRUIT-UART5")
}
