package bbio

// #include <stddef.h>
// #include <sys/types.h>
// #include <linux/i2c-dev.h>
import "C"

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// I2C is a port of https://github.com/bivab/smbus-cffi/
type I2C struct {
	file    *os.File
	address int
}

// Connects the object to the specified SMBus.
func NewI2C(bus, address int) (*I2C, error) {
	filename := fmt.Sprintf("/dev/i2c-%d", bus)
	file, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	i2c := &I2C{file: file, address: -1}
	err = i2c.SetAddress(address)
	if err != nil {
		file.Close()
		return nil, err
	}

	return i2c, nil
}

func (i2c *I2C) Address() int {
	return i2c.address
}

func (i2c *I2C) SetAddress(address int) error {
	if address != i2c.address {
		result, _, err := syscall.Syscall(syscall.SYS_IOCTL, i2c.file.Fd(), C.I2C_SLAVE, uintptr(address))
		if result != 0 {
			return err
		}
		i2c.address = address
	}
	return nil
}

// Perform SMBus Read Byte transaction.
func (i2c *I2C) ReadUint8() (uint8, error) {
	result, err := C.i2c_smbus_read_byte(C.int(i2c.file.Fd()))
	if result == -1 {
		return 0, err
	}
	return uint8(result), nil
}

// Perform SMBus Write Byte transaction.
func (i2c *I2C) WriteUint8(value uint8) error {
	result, err := C.i2c_smbus_write_byte(C.int(i2c.file.Fd()), C.__u8(value))
	if result == -1 {
		return err
	}
	return nil
}

// Perform SMBus Read Byte transaction.
func (i2c *I2C) ReadInt8() (int8, error) {
	result, err := C.i2c_smbus_read_byte(C.int(i2c.file.Fd()))
	if result == -1 {
		return 0, err
	}
	if result >= 1<<7 {
		result -= 1 << 8
	}
	return int8(result), nil
}

// Perform SMBus Write Byte transaction.
func (i2c *I2C) WriteInt8(value int8) error {
	return i2c.WriteUint8(uint8(value))
}

// Perform SMBus Read Byte Data transaction.
func (i2c *I2C) ReadUint8Cmd(command uint8) (uint8, error) {
	result, err := C.i2c_smbus_read_byte_data(C.int(i2c.file.Fd()), C.__u8(command))
	if result == -1 {
		return 0, err
	}
	return uint8(result), nil
}

// Perform SMBus Write Byte Data transaction.
func (i2c *I2C) WriteUint8Cmd(command uint8, value uint8) error {
	result, err := C.i2c_smbus_write_byte_data(C.int(i2c.file.Fd()), C.__u8(command), C.__u8(value))
	if result == -1 {
		return err
	}
	return nil
}

// Perform SMBus Read Byte Data transaction.
func (i2c *I2C) ReadInt8Cmd(command uint8) (int8, error) {
	result, err := C.i2c_smbus_read_byte_data(C.int(i2c.file.Fd()), C.__u8(command))
	if result == -1 {
		return 0, err
	}
	if result >= 1<<7 {
		result -= 1 << 8
	}
	return int8(result), nil
}

// Perform SMBus Write Byte Data transaction.
func (i2c *I2C) WriteInt8Cmd(command uint8, value int8) error {
	return i2c.WriteUint8Cmd(command, uint8(value))
}

// Perform SMBus Read Word Data transaction.
func (i2c *I2C) ReadUint16Cmd(command uint8) (uint16, error) {
	result, err := C.i2c_smbus_read_word_data(C.int(i2c.file.Fd()), C.__u8(command))
	if result == -1 {
		return 0, err
	}
	return uint16(result), nil
}

// Perform SMBus Write Word Data transaction.
func (i2c *I2C) WriteUint16Cmd(command uint8, value uint16) error {
	result, err := C.i2c_smbus_write_word_data(C.int(i2c.file.Fd()), C.__u8(command), C.__u16(value))
	if result == -1 {
		return err
	}
	return nil
}

// Perform SMBus Read Word Data transaction.
func (i2c *I2C) ReadInt16Cmd(command uint8) (int16, error) {
	result, err := C.i2c_smbus_read_word_data(C.int(i2c.file.Fd()), C.__u8(command))
	if result == -1 {
		return 0, err
	}
	if result >= 1<<15 {
		result -= 1 << 16
	}
	return int16(result), nil
}

// Perform SMBus Write Word Data transaction.
func (i2c *I2C) WriteInt16Cmd(command uint8, value int16) error {
	return i2c.WriteUint16Cmd(command, uint16(value))
}

// Perform SMBus Quick transaction.
func (i2c *I2C) WriteQuick() error {
	result, err := C.i2c_smbus_write_quick(C.int(i2c.file.Fd()), C.I2C_SMBUS_WRITE)
	if result != 0 {
		return err
	}
	return nil
}

// Perform SMBus Process Call transaction.
//
// Note: although i2c_smbus_process_call returns a value, according to
// smbusmodule.c this method does not return a value by default.
//
// Set _compat = False on the SMBus instance to get a return value.
func (i2c *I2C) ProcessCall(command uint8, value uint16) (uint16, error) {
	result, err := C.i2c_smbus_process_call(C.int(i2c.file.Fd()), C.__u8(command), C.__u16(value))
	if result == -1 {
		return 0, err
	}
	return uint16(result), nil
}

// Perform SMBus Block Process Call transaction.
func (i2c *I2C) ProcessCallBlock(command uint8, block []byte) ([]byte, error) {
	length := len(block)
	if length == 0 || length > C.I2C_SMBUS_BLOCK_MAX {
		return nil, fmt.Errorf("Length of block is %d, but must be in the range 1 to %d", length, C.I2C_SMBUS_BLOCK_MAX)
	}
	data := make([]byte, length+1, C.I2C_SMBUS_BLOCK_MAX+2)
	data[0] = byte(length)
	copy(data[1:], block)
	result, err := C.i2c_smbus_access(C.int(i2c.file.Fd()), C.I2C_SMBUS_WRITE, C.__u8(command), C.I2C_SMBUS_BLOCK_PROC_CALL, (*C.union_i2c_smbus_data)(unsafe.Pointer(&data[0])))
	if result != 0 {
		return nil, err
	}
	return data[1 : 1+data[0]], nil
}

// Perform SMBus Read Block Data transaction.
func (i2c *I2C) ReadBlock(command uint8) ([]byte, error) {
	data := make([]byte, C.I2C_SMBUS_BLOCK_MAX+2)
	result, err := C.i2c_smbus_access(C.int(i2c.file.Fd()), C.I2C_SMBUS_READ, C.__u8(command), C.I2C_SMBUS_BLOCK_DATA, (*C.union_i2c_smbus_data)(unsafe.Pointer(&data[0])))
	if result != 0 {
		return nil, err
	}
	return data[1 : 1+data[0]], nil
}

// Perform SMBus Write Block Data transaction.
func (i2c *I2C) WriteBlock(command uint8, block []byte) error {
	length := len(block)
	if length == 0 || length > C.I2C_SMBUS_BLOCK_MAX {
		return fmt.Errorf("Length of block is %d, but must be in the range 1 to %d", length, C.I2C_SMBUS_BLOCK_MAX)
	}
	data := make([]byte, length+1)
	data[0] = byte(length)
	copy(data[1:], block)
	result, err := C.i2c_smbus_access(C.int(i2c.file.Fd()), C.I2C_SMBUS_WRITE, C.__u8(command), C.I2C_SMBUS_BLOCK_DATA, (*C.union_i2c_smbus_data)(unsafe.Pointer(&data[0])))
	if result != 0 {
		return err
	}
	return nil
}

// TODO: Perform I2C Block Read transaction.
// With if len == 32 then arg = C.I2C_SMBUS_I2C_BLOCK_BROKEN

func (i2c *I2C) Close() error {
	return i2c.file.Close()
}
