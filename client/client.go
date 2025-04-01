package picoclient

import (
	"encoding/binary"

	"github.com/stipochka/pico_go/models"

	cwrapper "github.com/stipochka/pico_go/internal/wrapper"
)

type Client struct {
	PicoClient cwrapper.CWrapper
}

func NewClient(filename string) (*Client, error) {
	picoClient, err := cwrapper.NewGoWrapper(filename)
	if err != nil {
		return nil, err
	}
	return &Client{
		PicoClient: picoClient,
	}, nil
}

func convertToResponseModel(buff []byte) *models.McuResponse {
	var resp models.McuResponse
	resp.FuncCode = uint8(buff[0])
	resp.ErrCode = uint8(buff[1])
	resp.Args = binary.LittleEndian.Uint16(buff[2:4])
	resp.Buffer = buff[4 : len(buff)-2]
	return &resp
}

func (c *Client) HeatbitRequest() (*models.McuResponse, error) {
	resp, err := c.PicoClient.WrapHeartbitRequest()
	if err != nil {
		return nil, err
	}
	return convertToResponseModel(resp), nil
}

func (c *Client) GetActualDataRequest(sensorName string) (*models.McuResponse, error) {
	resp, err := c.PicoClient.WrapGetActualDataRequest(sensorName)
	if err != nil {
		return nil, err
	}
	return convertToResponseModel(resp), nil
}

func (c *Client) GetHistoryDataRequest(sensorName string, num uint16) (*models.McuResponse, error) {
	resp, err := c.PicoClient.WrapGetHistoryDataRequest(sensorName, num)
	if err != nil {
		return nil, err
	}
	return convertToResponseModel(resp), nil
}

func (c *Client) GetSensorInfoRequest(sensorName string) (*models.McuResponse, error) {
	resp, err := c.PicoClient.WrapGetSensorInfoRequest(sensorName)
	if err != nil {
		return nil, err
	}
	return convertToResponseModel(resp), nil
}

func (c *Client) GetMcuInfoRequest() (*models.McuResponse, error) {
	resp, err := c.PicoClient.WrapGetMcuInfoRequest()
	if err != nil {
		return nil, err
	}
	return convertToResponseModel(resp), nil
}

func (c *Client) SetReadingPeriodRequest(sensorName string, delay uint16) (*models.McuResponse, error) {
	resp, err := c.PicoClient.WrapSetReadingPeriodRequest(sensorName, delay)
	if err != nil {
		return nil, err
	}
	return convertToResponseModel(resp), nil
}

func (c *Client) Close() error {
	return c.PicoClient.Close()
}
