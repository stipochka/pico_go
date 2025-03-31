package picowriter

import (
	"fmt"
	"time"

	"github.com/tarm/serial"
)

type PicoProcessor interface {
	Write([]byte) error
	Read() ([]byte, error)
	Close() error
}

type UARTProcessor struct {
	port *serial.Port
}

func NewUartProcessor(filename string) (*UARTProcessor, error) {
	config := &serial.Config{
		Name:        filename,
		Baud:        115200,
		ReadTimeout: time.Second * 2,
		Size:        8,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return &UARTProcessor{
		port: port,
	}, nil
}

func (u *UARTProcessor) Write(buffer []byte) error {
	const op = "uart.Write"

	n, err := u.port.Write(buffer)
	if err != nil {
		fmt.Println(n)
		return fmt.Errorf("%s %w", op, err)
	}
	return nil
}

func (u *UARTProcessor) Read() ([]byte, error) {
	const op = "uart.Read"

	buffer := make([]byte, 133)
	n, err := u.port.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return buffer[:n], nil
}

func (u *UARTProcessor) Close() error {
	if u.port != nil {
		err := u.port.Close()
		u.port = nil
		return err
	}
	return nil
}
