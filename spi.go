package bbio

import (
	"os"
)

const (
	SPI_CPHA       = 0x01 /* clock phase */
	SPI_CPOL       = 0x02 /* clock polarity */
	SPI_MODE_0     = 0    /* (original MicroWire) */
	SPI_MODE_1     = SPI_CPHA
	SPI_MODE_2     = SPI_CPOL
	SPI_MODE_3     = SPI_CPOL | SPI_CPHA
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

type SPI struct {
	file *os.File /* open file descriptor: /dev/spi-X.Y */
	mode uint8    /* current SPI mode */
	bpw  uint8    /* current SPI bits per word setting */
	msh  uint32   /* current SPI max speed setting in Hz */
}

func NewSPI() (*SPI, error) {
	return new(SPI), nil
}

func (spi *SPI) Close() error {
	return spi.file.Close()
}
