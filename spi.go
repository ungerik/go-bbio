package bbio

// #include <linux/spi/spidev.h>
// #include <sys/ioctl.h>
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	SPI_CPHA       = 0x01  /* clock phase */
	SPI_CPOL       = 0x02  /* clock polarity */
	SPI_CS_HIGH    = 0x04  /* chipselect active high? */
	SPI_LSB_FIRST  = 0x08  /* per-word bits-on-wire */
	SPI_TRHEE_WIRE = 0x10  /* SI/SO signals shared */
	SPI_LOOP       = 0x20  /* loopback mode */
	SPI_NO_CS      = 0x40  /* 1 dev/bus, no chipselect */
	SPI_READY      = 0x80  /* slave pulls low to pause */
	SPI_TX_DUAL    = 0x100 /* transmit with 2 wires */
	SPI_TX_QUAD    = 0x200 /* transmit with 4 wires */
	SPI_RX_DUAL    = 0x400 /* receive with 2 wires */
	SPI_RX_QUAD    = 0x800 /* receive with 4 wires */
)

type SPIMode uint8

const (
	SPI_MODE_0 SPIMode = 0 /* (original MicroWire) */
	SPI_MODE_1 SPIMode = SPI_CPHA
	SPI_MODE_2 SPIMode = SPI_CPOL
	SPI_MODE_3 SPIMode = SPI_CPOL | SPI_CPHA
)

type SPI struct {
	file *os.File /* open file descriptor: /dev/spi-X.Y */
	mode uint8    /* current SPI mode */
	bpw  uint8    /* current SPI bits per word setting */
	msh  uint32   /* current SPI max speed setting in Hz */
}

func NewSPI() (*SPI, error) {
	return new(SPI), nil
}

func (spi *SPI) Read(data []byte) (n int, err error) {
	return spi.file.Read(data)
}

func (spi *SPI) Write(data []byte) (n int, err error) {
	return spi.file.Write(data)
}

func (spi *SPI) Close() error {
	return spi.file.Close()
}

func (spi *SPI) Mode() SPIMode {
	return SPIMode(spi.mode) & SPI_MODE_3
}

func (spi *SPI) SetMode(mode SPIMode) error {
	m := (spi.mode &^ (SPI_CPHA | SPI_CPOL)) | uint8(mode)
	err := spi.setMode(m)
	if err == nil {
		spi.mode = m
	}
	return err
}

func (spi *SPI) setMode(mode uint8) error {
	r, _, err := syscall.Syscall(syscall.SYS_IOCTL, spi.file.Fd(), C.SPI_IOC_WR_MODE, uintptr(unsafe.Pointer(&mode)))
	if r != 0 {
		return err
	}

	var test uint8
	r, _, err = syscall.Syscall(syscall.SYS_IOCTL, spi.file.Fd(), C.SPI_IOC_RD_MODE, uintptr(unsafe.Pointer(&test)))
	if r != 0 {
		return err
	}

	if test == mode {
		return nil
	} else {
		return fmt.Errorf("Could not set SPI mode %d", mode)
	}
}
