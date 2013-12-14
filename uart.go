package bbio

import (
	"fmt"
	"io"

	"github.com/huin/goserial"
)

type UARTNr int

const (
	UART1 UARTNr = 1
	UART2 UARTNr = 2
	UART4 UARTNr = 4
	UART5 UARTNr = 5

	UART_115200_BAUD = 115200
	UART_57600_BAUD  = 57600
	UART_38400_BAUD  = 38400
	UART_19200_BAUD  = 19200
	UART_9600_BAUD   = 9600
)

type UARTByteSize goserial.ByteSize
type UARTParityMode goserial.ParityMode
type UARTStopBits goserial.StopBits

const (
	UART_ParityNone = UARTParityMode(goserial.ParityNone)
	UART_ParityEven = UARTParityMode(goserial.ParityEven)
	UART_ParityOdd  = UARTParityMode(goserial.ParityOdd)

	UART_Byte8 = UARTByteSize(goserial.Byte8)
	UART_Byte5 = UARTByteSize(goserial.Byte5)
	UART_Byte6 = UARTByteSize(goserial.Byte6)
	UART_Byte7 = UARTByteSize(goserial.Byte7)

	UART_StopBits1 = UARTStopBits(goserial.StopBits1)
	UART_StopBits2 = UARTStopBits(goserial.StopBits2)
)

// var uartTable = map[UARTName]uartInfo{
// 	UART1: {"UART1", "/dev/ttyO1", "ADAFRUIT-UART1", "P9_26", "P9_24"},
// 	UART2: {"UART2", "/dev/ttyO2", "ADAFRUIT-UART2", "P9_22", "P9_21"},
// 	UART4: {"UART4", "/dev/ttyO4", "ADAFRUIT-UART4", "P9_11", "P9_13"},
// 	UART5: {"UART5", "/dev/ttyO5", "ADAFRUIT-UART5", "P8_38", "P8_37"},
// }

// Wraps "github.com/huin/goserial"
type UART struct {
	io.ReadWriteCloser
	nr UARTNr
}

func NewUART(nr UARTNr, baud int, size UARTByteSize, parity UARTParityMode, stopBits UARTStopBits) (*UART, error) {
	dt := fmt.Sprintf("ADAFRUIT-UART%d", nr)
	err := LoadDeviceTree(dt)
	if err != nil {
		return nil, err
	}

	uart := &UART{nr: nr}

	config := &goserial.Config{
		Name:     fmt.Sprintf("/dev/ttyO%d", nr),
		Baud:     baud,
		Size:     goserial.ByteSize(size),
		Parity:   goserial.ParityMode(parity),
		StopBits: goserial.StopBits(stopBits),
	}
	uart.ReadWriteCloser, err = goserial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return uart, nil
}

func (uart *UART) Close() error {
	err := uart.ReadWriteCloser.Close()
	if err != nil {
		return err
	}
	return UnloadDeviceTree(fmt.Sprintf("ADAFRUIT-UART%d", uart.nr))
}
