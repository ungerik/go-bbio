package spi

const (
	CPHA       = 0x01 /* clock phase */
	CPOL       = 0x02 /* clock polarity */
	MODE_0     = 0    /* (original MicroWire) */
	MODE_1     = CPHA
	MODE_2     = CPOL
	MODE_3     = CPOL | CPHA
	CS_HIGH    = 0x04  /* chipselect active high? */
	LSB_FIRST  = 0x08  /* per-word bits-on-wire */
	TRHEE_WIRE = 0x10  /* SI/SO signals shared */
	LOOP       = 0x20  /* loopback mode */
	NO_CS      = 0x40  /* 1 dev/bus, no chipselect */
	READY      = 0x80  /* slave pulls low to pause */
	TX_DUAL    = 0x100 /* transmit with 2 wires */
	TX_QUAD    = 0x200 /* transmit with 4 wires */
	RX_DUAL    = 0x400 /* receive with 2 wires */
	RX_QUAD    = 0x800 /* receive with 4 wires */
)
