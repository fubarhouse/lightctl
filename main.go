package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/kingpin"
)

const (
	settingMaximumBrightness  = 100
	settingMaximumTemperature = 344
	settingMinimumBrightness  = 3
	settingMinimumTemperature = 143
)

var (
	// Thanks of the extra special kind to the following resources for making this possible.
	Thanks = []string{
		"https://mensfeld.pl/2021/12/controlling-elgato-key-light-under-ubuntu-with-ruby/",
		"https://github.com/adamesch/elgato-key-light-api",
	}

	// CLI configuration via Kingpin.
	app = kingpin.New("chat", "A command-line chat application.")

	on          = app.Command("on", "Turn your light(s) on")
	off         = app.Command("off", "Turn your light(s) off")
	toggle      = app.Command("toggle", "Toggle your light(s) on or off")
	info        = app.Command("info", "Display your light(s) information")
	status      = app.Command("state", "Display your light(s) state information")
	brightness  = app.Command("brightness", fmt.Sprintf("Set your light(s) brightness to a specific value or decrement/increment (between %v and %v)", settingMinimumBrightness, settingMaximumBrightness))
	temperature = app.Command("temperature", fmt.Sprintf("Set your light(s) temperature to a specific value or decrement/increment (between %v and %v)", settingMinimumTemperature, settingMaximumTemperature))

	destinationValue = app.Flag("value", "value to change to. ie. 40 || 600 || _40 || +40").String()
	destinationIPs   = app.Flag("ip", "IP addresses, multiple supported").Default("10.0.0.90", "10.0.0.91").IPList()
	destinationPort  = app.Flag("port", "Port to use, single value support only").Default("9123").String()

	valueProperties = ValueProperties{}
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
	}

	// APIRequest matches our expected response.
	APIRequest APIResponse

	// ValueProperties describes the input value.
	ValueProperties struct {
		IsNegative bool
		IsNeutral  bool
		IsPositive bool
		Value      int
	}
)

// ErrBrightnessOuterBounds will return a custom error message for brightness limitations based on input.
func ErrBrightnessOuterBounds(value int) error {
	return errors.New(fmt.Sprintf("Needs to be between %v and %v, was: %v", settingMinimumBrightness, settingMaximumBrightness, value))
}

// ErrTemperatureOuterBounds will return a custom error message for temperature limitations based on input.
func ErrTemperatureOuterBounds(value int) error {
	return errors.New(fmt.Sprintf("Needs to be between %v and %v, was: %v", settingMinimumTemperature, settingMaximumTemperature, value))
}

// ParseValue will parse the input value to determine specific characteristics based upon text input.
func ParseValue(value string) error {
	if value == "" {
		return nil
	}
	if strings.HasPrefix(*destinationValue, "+") {
		valueProperties.IsPositive = true
		trimmed := strings.TrimLeft(*destinationValue, "+")
		value, err := strconv.Atoi(trimmed)
		if err != nil {
			return err
		}
		valueProperties.Value = value
	} else if strings.HasPrefix(*destinationValue, "_") {
		fmt.Println(3)
		valueProperties.IsNegative = true
		trimmed := strings.TrimLeft(*destinationValue, "_")
		value, err := strconv.Atoi(trimmed)
		if err != nil {
			return err
		}
		fmt.Println(3)
		valueProperties.Value = value
	} else {
		valueProperties.IsNeutral = true
		value, err := strconv.Atoi(*destinationValue)
		if err != nil {
			return err
		}
		valueProperties.Value = value
	}
	return nil
}

// UnmarshalResponse will print an unmarshalled APIResponse.
func UnmarshalResponse(body string) (APIResponse, error) {
	data := []byte(body)
	var response APIResponse
	err := json.Unmarshal(data, &response)
	fmt.Println(body)
	return response, err
}

// DispatchResponse will submit an API call to a light on a given IP with the input payload.
func DispatchResponse(ip string, payload APIRequest, endpoint string, method string) (APIResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return APIResponse{}, err
	}
	buf := bytes.NewBufferString(string(jsonData))
	url := fmt.Sprintf("http://%s:%v/elgato/%s", ip, *destinationPort, endpoint)
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
	return UnmarshalResponse(string(responseData))
}

// state will return the response from a 'lights' API call.
func state(ip string) (APIResponse, error) {
	return DispatchResponse(ip, APIRequest{}, "lights", "GET")
}

// GetInfo will show the response from an 'accessory-info' API call.
func GetInfo(ip string) (APIResponse, error) {
	return DispatchResponse(ip, APIRequest{}, "accessory-info", "GET")
}

// GetStatus will return the state property of a light on a given IP.
func GetStatus(ip string) int {
	resp, _ := state(ip)
	return resp.Lights[0].On
}

// GetTemperature will return the temperature property of a light on a given IP.
func GetTemperature(ip string) int {
	resp, _ := state(ip)
	return resp.Lights[0].Temperature
}

// GetBrightness will return the brightness property of a light on a given IP.
func GetBrightness(ip string) int {
	resp, _ := state(ip)
	return resp.Lights[0].Brightness
}

// SetStatusToggle will toggle a light on a given IP to have the opposite value as currently set.
func SetStatusToggle(ip string) (APIResponse, error) {
	currentState := GetStatus(ip)
	Request := APIRequest{}

	if currentState == 0 {
		Request = APIRequest{
			NumberOfLights: 1,
			Lights: []Light{{
				On: 1,
			}},
		}
	}
	if currentState == 1 {
		Request = APIRequest{
			NumberOfLights: 1,
			Lights: []Light{{
				On: 0,
			}},
		}
	}

	return DispatchResponse(ip, Request, "lights", "PUT")
}

// SetStatusOff will set a light on a given IP to have an off state.
func SetStatusOff(ip string) (APIResponse, error) {
	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On: 0,
		}},
	}
	return DispatchResponse(ip, Request, "lights", "PUT")
}

// SetStatusOn will set a light on a given IP to have an on state.
func SetStatusOn(ip string) (APIResponse, error) {
	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On: 1,
		}},
	}
	return DispatchResponse(ip, Request, "lights", "PUT")
}

// SetTemperature will set a light on a given IP to a specific temperature value.
func SetTemperature(ip string, value int) (APIResponse, error) {

	if !valueProperties.IsNeutral {
		t := GetTemperature(ip)
		if valueProperties.IsPositive {
			value = t + value
		} else if valueProperties.IsNegative {
			value = t - value
		}
	}

	if value < 143 || value > 344 {
		return APIResponse{}, ErrTemperatureOuterBounds(value)
	}

	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On:          1,
			Temperature: value,
		}},
	}

	return DispatchResponse(ip, Request, "lights", "PUT")
}

// AddTemperature will increase a light on a given IP by a specific temperature value.
func AddTemperature(ip string, value int) (APIResponse, error) {
	curr := GetTemperature(ip)
	return SetTemperature(ip, curr+value)
}

// MinusTemperature will decrease a light on a given IP by a specific temperature value.
// todo not working.
func MinusTemperature(ip string, value int) (APIResponse, error) {
	curr := GetTemperature(ip)
	return SetTemperature(ip, curr-value)
}

// SetBrightness will set a light on a given IP to a specific brightness value.
func SetBrightness(ip string, value int) (APIResponse, error) {

	if !valueProperties.IsNeutral {
		b := GetBrightness(ip)
		if valueProperties.IsPositive {
			value = b + value
		} else if valueProperties.IsNegative {
			value = b - value
		}
	}

	if value < 3 || value > 100 {
		return APIResponse{}, ErrBrightnessOuterBounds(value)
	}

	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On:         1,
			Brightness: value,
		}},
	}

	return DispatchResponse(ip, Request, "lights", "PUT")
}

func main() {

	app.HelpFlag.Short('h')
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case on.FullCommand():
		for _, v := range *destinationIPs {
			_, err := SetStatusOn(v.String())
			if err != nil {
				fmt.Println(err)
			}
		}

	case off.FullCommand():
		for _, v := range *destinationIPs {
			_, err := SetStatusOff(v.String())
			if err != nil {
				fmt.Println(err)
			}
		}

	case toggle.FullCommand():
		for _, v := range *destinationIPs {
			_, err := SetStatusToggle(v.String())
			if err != nil {
				fmt.Println(err)
			}
		}

	case info.FullCommand():
		for _, v := range *destinationIPs {
			_, err := GetInfo(v.String())
			if err != nil {
				fmt.Println(err)
			}
		}

	case status.FullCommand():
		for _, v := range *destinationIPs {
			_, err := state(v.String())
			if err != nil {
				fmt.Println(err)
			}
		}

	case brightness.FullCommand():
		err := ParseValue(*destinationValue)
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range *destinationIPs {
			_, err := SetBrightness(v.String(), valueProperties.Value)
			if err != nil {
				fmt.Println(err)
			}
		}

	case temperature.FullCommand():
		err := ParseValue(*destinationValue)
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range *destinationIPs {
			_, err := SetTemperature(v.String(), valueProperties.Value)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
