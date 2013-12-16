package bbio

// #include <stddef.h>
// #include <sys/types.h>
// #include <linux/i2c-dev.h>
import "C"

import (
	"fmt"
	"os"
	"syscall"
)

// #include <linux/i2c.h>

type I2C struct {
	address int
	file    *os.File
}

func NewI2C(address int) (*I2C, error) {
	filename := fmt.Sprintf("/dev/i2c-%d", address)
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	_, _, err = syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), C.I2C_SLAVE, uintptr(address))
	if err != nil {
		return nil, err
	}

	return &I2C{address, file}, nil
}

func (i2c *I2C) errMsg() error {
	return fmt.Errorf("Error accessing 0x%02X: Check your I2C address", i2c.address)
}

func (i2c *I2C) ReadUint8() (uint8, error) {
	result := C.i2c_smbus_read_byte(C.int(i2c.file.Fd()))
	if result == -1 {
		return 0, i2c.errMsg()
	}
	return uint8(result), nil
}

func (i2c *I2C) WriteUint8(value uint8) error {
	result := C.i2c_smbus_write_byte(C.int(i2c.file.Fd()), C.__u8(value))
	if result == -1 {
		return i2c.errMsg()
	}
	return nil
}

func (i2c *I2C) ReadInt8() (int8, error) {
	result := C.i2c_smbus_read_byte(C.int(i2c.file.Fd()))
	if result == -1 {
		return 0, i2c.errMsg()
	}
	if result >= 1<<7 {
		result -= 1 << 8
	}
	return int8(result), nil
}

func (i2c *I2C) WriteInt8(value int8) error {
	result := C.i2c_smbus_write_byte(C.int(i2c.file.Fd()), C.__u8(value))
	if result == -1 {
		return i2c.errMsg()
	}
	return nil
}

func (i2c *I2C) ReadUint8Cmd(command uint8) (uint8, error) {
	result := C.i2c_smbus_read_byte_data(C.int(i2c.file.Fd()), C.__u8(command))
	if result == -1 {
		return 0, i2c.errMsg()
	}
	return uint8(result), nil
}

func (i2c *I2C) WriteUint8Cmd(command uint8, value uint8) error {
	result := C.i2c_smbus_write_byte_data(C.int(i2c.file.Fd()), C.__u8(command), C.__u8(value))
	if result == -1 {
		return i2c.errMsg()
	}
	return nil
}

func (i2c *I2C) ReadInt8Cmd(command uint8) (int8, error) {
	result := C.i2c_smbus_read_byte_data(C.int(i2c.file.Fd()), C.__u8(command))
	if result == -1 {
		return 0, i2c.errMsg()
	}
	if result >= 1<<7 {
		result -= 1 << 8
	}
	return int8(result), nil
}

func (i2c *I2C) WriteInt8Cmd(command uint8, value int8) error {
	result := C.i2c_smbus_write_byte_data(C.int(i2c.file.Fd()), C.__u8(command), C.__u8(value))
	if result == -1 {
		return i2c.errMsg()
	}
	return nil
}

func (i2c *I2C) ReadUint16Cmd(command uint8) (uint16, error) {
	result := C.i2c_smbus_read_word_data(C.int(i2c.file.Fd()), C.__u8(command))
	if result == -1 {
		return 0, i2c.errMsg()
	}
	return uint16(result), nil
}

func (i2c *I2C) WriteUint16Cmd(command uint8, value uint16) error {
	result := C.i2c_smbus_write_word_data(C.int(i2c.file.Fd()), C.__u8(command), C.__u16(value))
	if result == -1 {
		return i2c.errMsg()
	}
	return nil
}

func (i2c *I2C) ReadInt16Cmd(command uint8) (int16, error) {
	result := C.i2c_smbus_read_word_data(C.int(i2c.file.Fd()), C.__u8(command))
	if result == -1 {
		return 0, i2c.errMsg()
	}
	if result >= 1<<15 {
		result -= 1 << 16
	}
	return int16(result), nil
}

func (i2c *I2C) WriteInt16Cmd(command uint8, value int16) error {
	result := C.i2c_smbus_write_word_data(C.int(i2c.file.Fd()), C.__u8(command), C.__u16(value))
	if result == -1 {
		return i2c.errMsg()
	}
	return nil
}

func (i2c *I2C) WriteQuick() error {
	result := C.i2c_smbus_write_quick(C.int(i2c.file.Fd()), C.I2C_SMBUS_WRITE)
	if result != 0 {
		return i2c.errMsg()
	}
	return nil
}

func (i2c *I2C) ProcessCall(command uint8, value uint16) error {
	result := C.i2c_smbus_process_call(C.int(i2c.file.Fd()), C.__u8(command), C.__u16(value))
	if result == -1 {
		return i2c.errMsg()
	}
	return nil
}

func (i2c *I2C) Read(data []byte) (n int, err error) {
	return
}

func (i2c *I2C) Write(data []byte) (n int, err error) {
	return
}

func (i2c *I2C) Close() error {
	return i2c.file.Close()
}
