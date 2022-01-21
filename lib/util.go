package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type (
	// Light contains properties the API expects each light to return.
	Light struct {
		On          int `json:"on`
		Brightness  int `json:"brightness,omitempty"`
		Temperature int `json:"temperature,omitempty"`
	}

	// APIResponse is the Response we expect back from the API.
	APIResponse struct {
		NumberOfLights int     `json:"numberOfLights"`
		Lights         []Light `json:"lights"`

		// Support for the /accessory-info API endpoint:

		ProductName         string   `json:"productName,omitempty"`
		HardwareBoardType   int      `json:"hardwareBoardType,omitempty"`
		FirmwareBuildNumber int      `json:"firmwareBuildNumber,omitempty"`
		FirmwareVersion     string   `json:"firmwareVersion,omitempty"`
		SerialNumber        string   `json:"serialNumber,omitempty"`
		DisplayName         string   `json:"displayName,omitempty"`
		Features            []string `json:"features,omitempty"`
	}

	// APIRequest matches our expected response.
	APIRequest APIResponse

	InfoAPIResponse struct {
	}

	// ValueProperties describes the input value.
	ValueProperties struct {
		IsNegative bool
		IsNeutral  bool
		IsPositive bool
		Value      int
	}

	client struct {
		Request    *APIRequest
		Response   *APIResponse
		Properties struct {
			IPs             []net.IP
			Port            string
			RequestModifier ValueProperties
		}
	}
)

// UnmarshalResponse will print an unmarshalled APIResponse.
func (c *client) UnmarshalResponse(body string) (APIResponse, error) {
	data := []byte(body)
	var response APIResponse
	err := json.Unmarshal(data, &response)
	c.Response = &response
	return response, err
}

// DispatchResponse will submit an API call to a light on a given IP with the input payload.
func (c *client) DispatchResponse(ip string, payload APIRequest, endpoint string, method string) (APIResponse, error) {
	c.Request = &payload
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return APIResponse{}, err
	}
	buf := bytes.NewBufferString(string(jsonData))
	url := fmt.Sprintf("http://%s:%v/elgato/%s", ip, c.Properties.Port, endpoint)
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return APIResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return APIResponse{}, err
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, err
	}

	defer resp.Body.Close()

	return c.UnmarshalResponse(string(responseData))
}

// ParseValue will parse the input value to determine specific characteristics based upon text input.
func (c *client) ParseValue(input string) error {
	if input == "" {
		return nil
	}
	if strings.HasPrefix(input, "+") {
		c.Properties.RequestModifier.IsPositive = true
		trimmed := strings.TrimLeft(input, "+")
		value, err := strconv.Atoi(trimmed)
		if err != nil {
			return err
		}
		c.Properties.RequestModifier.Value = value
	} else if strings.HasPrefix(input, "_") {
		c.Properties.RequestModifier.IsNegative = true
		trimmed := strings.TrimLeft(input, "_")
		value, err := strconv.Atoi(trimmed)
		if err != nil {
			return err
		}
		c.Properties.RequestModifier.Value = value
	} else {
		c.Properties.RequestModifier.IsNeutral = true
		value, err := strconv.Atoi(input)
		if err != nil {
			return err
		}
		c.Properties.RequestModifier.Value = value
	}
	return nil
}

// State will return the response from a 'lights' API call.
func (c *client) State(ip string) (APIResponse, error) {
	return c.DispatchResponse(ip, APIRequest{}, "lights", "GET")
}

// PrintRequest will print the request in JSON format.
func (c *client) PrintRequest() {
	result, _ := json.Marshal(c.Request)
	fmt.Println(string(result))
}

// PrintResponse will print the response in JSON format.
func (c *client) PrintResponse() {
	result, _ := json.Marshal(c.Response)
	fmt.Println(string(result))
}

// add will add two values.
func add(a int, b int) int {
	return a + b
}

// subtract will subtract a smaller number from a larger number.
func subtract(a int, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}
