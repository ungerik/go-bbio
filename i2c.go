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

func SwapBytes(word uint16) uint16 {
	return word<<8 | word>>8
}

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

func (i2c *I2C) smbusAccess(readWrite, command uint8, size int, data unsafe.Pointer) (uintptr, error) {
	args := C.struct_i2c_smbus_ioctl_data{
		read_write: C.char(readWrite),
		command:    C.__u8(command),
		size:       C.int(size),
		data:       (*C.union_i2c_smbus_data)(data),
	}
	result, _, err := syscall.Syscall(syscall.SYS_IOCTL, i2c.file.Fd(), C.I2C_SMBUS, uintptr(unsafe.Pointer(&args)))
	if int(result) == -1 {
		return 0, err
	}
	return result, nil
}

// WriteQuick sends a single bit to the device, at the place of the Rd/Wr bit.
func (i2c *I2C) WriteQuick(value uint8) error {
	_, err := i2c.smbusAccess(value, 0, C.I2C_SMBUS_QUICK, nil)
	return err
}

// ReadUint8 reads a single byte from a device, without specifying a device
// register. Some devices are so simple that this interface is enough; for
// others, it is a shorthand if you want to read the same register as in
// the previous SMBus command.
func (i2c *I2C) ReadUint8() (result uint8, err error) {
	_, err = i2c.smbusAccess(C.I2C_SMBUS_READ, 0, C.I2C_SMBUS_BYTE, unsafe.Pointer(&result))
	if err != nil {
		return 0, err
	}
	return 0xFF & result, nil
}

// WriteUint8 sends a single byte to a device.
func (i2c *I2C) WriteUint8(value uint8) error {
	_, err := i2c.smbusAccess(C.I2C_SMBUS_WRITE, value, C.I2C_SMBUS_BYTE, nil)
	return err
}

// ReadInt8 reads a single byte from a device, without specifying a device
// register. Some devices are so simple that this interface is enough; for
// others, it is a shorthand if you want to read the same register as in
// the previous SMBus command.
func (i2c *I2C) ReadInt8() (int8, error) {
	result, err := i2c.ReadUint8()
	return int8(result), err
}

// WriteInt8 sends a single byte to a device.
func (i2c *I2C) WriteInt8(value int8) error {
	return i2c.WriteUint8(uint8(value))
}

// ReadUint8Cmd reads a single byte from a device, from a designated register.
// The register is specified through the command byte.
func (i2c *I2C) ReadUint8Cmd(command uint8) (result uint8, err error) {
	_, err = i2c.smbusAccess(C.I2C_SMBUS_READ, command, C.I2C_SMBUS_BYTE_DATA, unsafe.Pointer(&result))
	if err != nil {
		return 0, err
	}
	return 0xFF & result, nil
}

// WriteUint8Cmd writes a single byte to a device, to a designated register. The
// register is specified through the command byte.
func (i2c *I2C) WriteUint8Cmd(command uint8, value uint8) error {
	_, err := i2c.smbusAccess(C.I2C_SMBUS_WRITE, command, C.I2C_SMBUS_BYTE_DATA, unsafe.Pointer(&value))
	return err
}

// ReadInt8Cmd reads a single byte from a device, from a designated register.
// The register is specified through the command byte.
func (i2c *I2C) ReadInt8Cmd(command uint8) (int8, error) {
	result, err := i2c.ReadUint8Cmd(command)
	return int8(result), err
}

// WriteInt8Cmd writes a single byte to a device, to a designated register. The
// register is specified through the command byte.
func (i2c *I2C) WriteInt8Cmd(command uint8, value int8) error {
	return i2c.WriteUint8Cmd(command, uint8(value))
}

// ReadUint16Cmd is very like ReadUint8Cmd; again, data is read from a
// device, from a designated register that is specified through the command
// byte. But this time, the data is a complete word (16 bits).
func (i2c *I2C) ReadUint16Cmd(command uint8) (result uint16, err error) {
	_, err = i2c.smbusAccess(C.I2C_SMBUS_READ, command, C.I2C_SMBUS_WORD_DATA, unsafe.Pointer(&result))
	if err != nil {
		return 0, err
	}
	return 0xFFFF & result, nil
}

// WriteUint16Cmd is the opposite of the ReadUint16Cmd operation. 16 bits
// of data is written to a device, to the designated register that is
// specified through the command byte.
func (i2c *I2C) WriteUint16Cmd(command uint8, value uint16) error {
	_, err := i2c.smbusAccess(C.I2C_SMBUS_WRITE, command, C.I2C_SMBUS_WORD_DATA, unsafe.Pointer(&value))
	return err
}

// ReadUint16CmdSwapped is very like ReadUint8Cmd; again, data is read from a
// device, from a designated register that is specified through the command
// byte. But this time, the data is a complete word (16 bits).
// The bytes of the 16 bit value will be swapped.
func (i2c *I2C) ReadUint16CmdSwapped(command uint8) (result uint16, err error) {
	result, err = i2c.ReadUint16Cmd(command)
	return SwapBytes(result), err
}

// WriteUint16CmdSwapped is the opposite of the ReadUint16CmdSwapped operation. 16 bits
// of data is written to a device, to the designated register that is
// specified through the command byte.
// The bytes of the 16 bit value will be swapped.
func (i2c *I2C) WriteUint16CmdSwapped(command uint8, value uint16) error {
	return i2c.WriteUint16Cmd(command, SwapBytes(value))
}

// ReadInt16Cmd is very like ReadInt8Cmd; again, data is read from a
// device, from a designated register that is specified through the command
// byte. But this time, the data is a complete word (16 bits).
func (i2c *I2C) ReadInt16Cmd(command uint8) (int16, error) {
	result, err := i2c.ReadUint16Cmd(command)
	return int16(result), err
}

// WriteInt16Cmd is the opposite of the ReadInt16Cmd operation. 16 bits
// of data is written to a device, to the designated register that is
// specified through the command byte.
func (i2c *I2C) WriteInt16Cmd(command uint8, value int16) error {
	return i2c.WriteUint16Cmd(command, uint16(value))
}

// ReadInt16CmdSwapped is very like ReadInt8CmdSwapped; again, data is read from a
// device, from a designated register that is specified through the command
// byte. But this time, the data is a complete word (16 bits).
// The bytes of the 16 bit value will be swapped.
func (i2c *I2C) ReadInt16CmdSwapped(command uint8) (int16, error) {
	result, err := i2c.ReadUint16CmdSwapped(command)
	return int16(result), err
}

// WriteInt16CmdSwapped is the opposite of the ReadInt16CmdSwapped operation. 16 bits
// of data is written to a device, to the designated register that is
// specified through the command byte.
// The bytes of the 16 bit value will be swapped.
func (i2c *I2C) WriteInt16CmdSwapped(command uint8, value int16) error {
	return i2c.WriteUint16CmdSwapped(command, uint16(value))
}

// ProcessCall selects a device register (through the command byte), sends
// 16 bits of data to it, and reads 16 bits of data in return.
func (i2c *I2C) ProcessCall(command uint8, value uint16) (uint16, error) {
	_, err := i2c.smbusAccess(C.I2C_SMBUS_WRITE, command, C.I2C_SMBUS_PROC_CALL, unsafe.Pointer(&value))
	if err != nil {
		return 0, err
	}
	return 0xFFFF & value, nil
}

// ProcessCallSwapped selects a device register (through the command byte), sends
// 16 bits of data to it, and reads 16 bits of data in return.
// The bytes of the 16 bit value will be swapped.
func (i2c *I2C) ProcessCallSwapped(command uint8, value uint16) (uint16, error) {
	result, err := i2c.ProcessCall(command, SwapBytes(value))
	return SwapBytes(result), err
}

// ProcessCallBlock reads a block of up to 32 bytes from a device, from a
// designated register that is specified through the command byte.
func (i2c *I2C) ProcessCallBlock(command uint8, block []byte) ([]byte, error) {
	length := len(block)
	if length == 0 || length > C.I2C_SMBUS_BLOCK_MAX {
		return nil, fmt.Errorf("Length of block is %d, but must be in the range 1 to %d", length, C.I2C_SMBUS_BLOCK_MAX)
	}
	data := make([]byte, length+1, C.I2C_SMBUS_BLOCK_MAX+2)
	data[0] = byte(length)
	copy(data[1:], block)
	_, err := i2c.smbusAccess(C.I2C_SMBUS_WRITE, command, C.I2C_SMBUS_BLOCK_PROC_CALL, unsafe.Pointer(&data[0]))
	if err != nil {
		return nil, err
	}
	return data[1 : 1+data[0]], nil
}

// ReadBlock writes up to 32 bytes to a device, to a designated
// register that is specified through the command byte.
func (i2c *I2C) ReadBlock(command uint8) ([]byte, error) {
	data := make([]byte, C.I2C_SMBUS_BLOCK_MAX+2)
	_, err := i2c.smbusAccess(C.I2C_SMBUS_READ, command, C.I2C_SMBUS_BLOCK_DATA, unsafe.Pointer(&data[0]))
	if err != nil {
		return nil, err
	}
	return data[1 : 1+data[0]], nil
}

// WriteBlock selects a device register (through the command byte), sends
// 1 to 31 bytes of data to it, and reads 1 to 31 bytes of data in return.
func (i2c *I2C) WriteBlock(command uint8, block []byte) error {
	length := len(block)
	if length == 0 || length > C.I2C_SMBUS_BLOCK_MAX {
		return fmt.Errorf("Length of block is %d, but must be in the range 1 to %d", length, C.I2C_SMBUS_BLOCK_MAX)
	}
	data := make([]byte, length+1)
	data[0] = byte(length)
	copy(data[1:], block)
	_, err := i2c.smbusAccess(C.I2C_SMBUS_WRITE, command, C.I2C_SMBUS_BLOCK_DATA, unsafe.Pointer(&data[0]))
	return err
}

// TODO: Perform I2C Block Read transaction.
// With if len == 32 then arg = C.I2C_SMBUS_I2C_BLOCK_BROKEN instead of I2C_SMBUS_I2C_BLOCK_DATA ???

func (i2c *I2C) Close() error {
	return i2c.file.Close()
}
