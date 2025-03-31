package cwrapper

/*
#cgo CFLAGS: -I./lib/include
#cgo LDFLAGS: -L./lib -lpicolib

#include "master.h"
#include "picolib.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	picowriter "pico_go/internal/uart"
	"unsafe"
)

type PlError int

const (
	PlErrorNone                PlError = iota // PL_ERROR_NONE
	PlErrorNoPackage                          // PL_ERROR_NO_PACKAGE
	PlErrorInvalidSensorName                  // PL_ERROR_INVALID_SENSOR_NAME
	PlErrorInvalidArgument                    // PL_ERROR_INVALID_ARGUMENT
	PlErrorInvalidFunctionCode                // PL_ERROR_INVALID_FUNCTION_CODE
	PlErrorInvalidCRC                         // PL_ERROR_INVALID_CRC
	PlErrorRequestTimeout                     // PL_ERROR_REQUEST_TIMEOUT
	PlErrorTransport                          // PL_ERROR_TRANSPORT
)

// Map кодов ошибок в Go-ошибки
var plErrorMap = map[PlError]error{
	PlErrorNone:                nil,
	PlErrorNoPackage:           errors.New("no package received"),
	PlErrorInvalidSensorName:   errors.New("invalid sensor name"),
	PlErrorInvalidArgument:     errors.New("invalid argument"),
	PlErrorInvalidFunctionCode: errors.New("invalid function code"),
	PlErrorInvalidCRC:          errors.New("invalid CRC"),
	PlErrorRequestTimeout:      errors.New("request timeout"),
	PlErrorTransport:           errors.New("transport error"),
}

// Функция для конвертации C ошибки в Go ошибку
func ConvertPlError(code C.pl_error_t) error {
	if err, exists := plErrorMap[PlError(code)]; exists {
		return err
	}
	return errors.New("unknown error")
}

type CWrapper interface {
	WrapHeartbitRequest() ([]byte, error)
	WrapGetActualDataRequest(sensor_name string) ([]byte, error)
	WrapGetHistoryDataRequest(sensor_name string, num uint16) ([]byte, error)
	WrapGetSensorInfoRequest(sensor_name string) ([]byte, error)
	WrapGetMcuInfoRequest() ([]byte, error)
	WrapSetReadingPeriodRequest(sensor_name string, delay uint16) ([]byte, error)
	Close() error
}

const bufferSize = 133

/*
pl_error_t heartbit_request(uint8_t* data, uint16_t *data_size);

pl_error_t get_actual_data_request(const uint8_t* sensor_name, uint8_t* data, uint16_t *data_size);

pl_error_t get_history_data_request(const uint8_t* sensor_name, uint16_t num, uint8_t* data, uint16_t *data_size);

pl_error_t get_sensor_info_request(const uint8_t* sensor_name, uint8_t* data, uint16_t *data_size);

pl_error_t get_mcu_info_request(uint8_t* data, uint16_t *data_size);

pl_error_t set_reading_period_request(const uint8_t* sensor_name, uint16_t delay, uint8_t* data, uint16_t *data_size);
*/

func checkCRC(buff []byte) bool {
	if len(buff) < bufferSize {
		return false
	}

	ucharBuff := unsafe.Slice((*C.uchar)(unsafe.Pointer(&buff[0])), len(buff)-1)
	expectedSum := int(C.calculate_crc8(&ucharBuff[0], C.size_t(len(buff)-1)))
	checkSumResp := int(buff[len(buff)-1])

	return checkSumResp == expectedSum
}

type UARTWrapper struct {
	Processor picowriter.PicoProcessor
}

func NewGoWrapper(filename string) (*UARTWrapper, error) {
	processor, err := picowriter.NewUartProcessor(filename)
	if err != nil {
		return nil, err
	}

	return &UARTWrapper{
		Processor: processor,
	}, nil
}

func (u *UARTWrapper) WrapHeartbitRequest() ([]byte, error) {

	var buffer [bufferSize]C.char
	dataSize := C.uint16_t(len(buffer))

	errCode := C.heartbit_request((*C.uchar)(unsafe.Pointer(&buffer[0])), &dataSize)
	if err := ConvertPlError(errCode); err != nil {
		return []byte{}, err
	}

	reqBuffer := C.GoBytes(unsafe.Pointer(&buffer[0]), C.int(dataSize))

	err := u.Processor.Write(reqBuffer)
	if err != nil {
		return nil, err
	}
	resp, err := u.Processor.Read()
	if err != nil {
		return nil, err
	}

	if isValidSum := checkCRC(resp); !isValidSum {
		return nil, errors.New("invalid checksum")
	}

	return resp, nil
}

func (u *UARTWrapper) WrapGetActualDataRequest(sensor_name string) ([]byte, error) {
	var buffer [bufferSize]C.char
	dataSize := C.uint16_t(len(buffer))

	cStr := C.CString(sensor_name)
	defer C.free(unsafe.Pointer(cStr))

	errCode := C.get_actual_data_request((*C.uchar)(unsafe.Pointer(cStr)), (*C.uchar)(unsafe.Pointer(&buffer[0])), &dataSize)
	if err := ConvertPlError(errCode); err != nil {
		return []byte{}, err
	}
	reqBuffer := C.GoBytes(unsafe.Pointer(&buffer[0]), C.int(dataSize))

	err := u.Processor.Write(reqBuffer)
	if err != nil {
		return nil, err
	}
	resp, err := u.Processor.Read()
	if err != nil {
		return nil, err
	}

	if isValidSum := checkCRC(resp); !isValidSum {
		return nil, errors.New("invalid checksum")
	}

	return resp, nil
}

func (u *UARTWrapper) WrapGetHistoryDataRequest(sensor_name string, num uint16) ([]byte, error) {
	var buffer [bufferSize]C.char
	dataSize := C.uint16_t(len(buffer))

	cStr := C.CString(sensor_name)
	defer C.free(unsafe.Pointer(cStr))

	errCode := C.get_history_data_request((*C.uchar)(unsafe.Pointer(cStr)), C.uint16_t(num), (*C.uchar)(unsafe.Pointer(&buffer[0])), &dataSize)
	if err := ConvertPlError(errCode); err != nil {
		return []byte{}, err
	}
	reqBuffer := C.GoBytes(unsafe.Pointer(&buffer[0]), C.int(dataSize))

	err := u.Processor.Write(reqBuffer)
	if err != nil {
		return nil, err
	}

	resp, err := u.Processor.Read()
	if err != nil {
		return nil, err
	}

	if isValidSum := checkCRC(resp); !isValidSum {
		return nil, errors.New("invalid checksum")
	}

	return resp, nil
}

func (u *UARTWrapper) WrapGetSensorInfoRequest(sensor_name string) ([]byte, error) {
	var buffer [bufferSize]C.char
	dataSize := C.uint16_t(len(buffer))

	cStr := C.CString(sensor_name)
	defer C.free(unsafe.Pointer(cStr))

	errCode := C.get_sensor_info_request((*C.uchar)(unsafe.Pointer(cStr)), (*C.uchar)(unsafe.Pointer(&buffer[0])), &dataSize)
	if err := ConvertPlError(errCode); err != nil {
		return []byte{}, err
	}
	reqBuffer := C.GoBytes(unsafe.Pointer(&buffer[0]), C.int(dataSize))

	err := u.Processor.Write(reqBuffer)
	if err != nil {
		return nil, err
	}

	resp, err := u.Processor.Read()
	if err != nil {
		return nil, err
	}

	if isValidSum := checkCRC(resp); !isValidSum {
		return nil, errors.New("invalid checksum")
	}

	return resp, nil
}

func (u *UARTWrapper) WrapGetMcuInfoRequest() ([]byte, error) {
	var buffer [bufferSize]C.char
	dataSize := C.uint16_t(len(buffer))

	errCode := C.get_mcu_info_request((*C.uchar)(unsafe.Pointer(&buffer[0])), &dataSize)
	if err := ConvertPlError(errCode); err != nil {
		return []byte{}, err
	}
	reqBuffer := C.GoBytes(unsafe.Pointer(&buffer[0]), C.int(dataSize))

	err := u.Processor.Write(reqBuffer)
	if err != nil {
		return nil, err
	}
	resp, err := u.Processor.Read()
	if err != nil {
		return nil, err
	}

	if isValidSum := checkCRC(resp); !isValidSum {
		return nil, errors.New("invalid checksum")
	}

	return resp, nil
}

func (u *UARTWrapper) WrapSetReadingPeriodRequest(sensor_name string, delay uint16) ([]byte, error) {
	var buffer [bufferSize]C.char
	dataSize := C.uint16_t(len(buffer))

	cStr := C.CString(sensor_name)
	defer C.free(unsafe.Pointer(cStr))

	errCode := C.set_reading_period_request((*C.uchar)(unsafe.Pointer(cStr)), C.uint16_t(delay), (*C.uchar)(unsafe.Pointer(&buffer[0])), &dataSize)
	if err := ConvertPlError(errCode); err != nil {
		return []byte{}, err
	}

	reqBuffer := C.GoBytes(unsafe.Pointer(&buffer[0]), C.int(dataSize))

	err := u.Processor.Write(reqBuffer)
	if err != nil {
		return nil, err
	}

	resp, err := u.Processor.Read()
	if err != nil {
		return nil, err
	}

	if isValidSum := checkCRC(resp); !isValidSum {
		return nil, errors.New("invalid checksum")
	}

	return resp, nil
}

func (u *UARTWrapper) Close() error {
	return u.Processor.Close()
}
