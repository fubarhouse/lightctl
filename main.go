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

var (
	Thanks = []string{
		"https://mensfeld.pl/2021/12/controlling-elgato-key-light-under-ubuntu-with-ruby/",
		"https://github.com/adamesch/elgato-key-light-api/blob/master/resources/lights/GET_lights.md",
	}

	app   = kingpin.New("chat", "A command-line chat application.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()

	on          = app.Command("on", "Register a new user.")
	off         = app.Command("off", "Register a new user.")
	toggle      = app.Command("toggle", "Register a new user.")
	info        = app.Command("info", "Register a new user.")
	status      = app.Command("state", "Register a new user.")
	brightness  = app.Command("brightness", "Register a new user.")
	temperature = app.Command("temperature", "Register a new user.")

	destinationValue = app.Flag("value", "value to change to. ie. 40 || 600 || _40 || +40").String()
	LIGHTS           = app.Flag("ips", "IP addresses separated by a comma").Default("10.0.0.90", "10.0.0.91").Strings()
	PORT             = app.Flag("port", "IP addresses separated by a comma").Default("9123").String()
)

type (
	Light struct {
		On          int `json:"on`
		Brightness  int `json:"brightness,omitempty"`
		Temperature int `json:"temperature,omitempty"`
	}
	APIResponse struct {
		NumberOfLights int     `json:"numberOfLights"`
		Lights         []Light `json:"lights"`
	}
	APIRequest APIResponse
)

func MarshalResponse(body string) (APIResponse, error) {
	data := []byte(body)
	var response APIResponse
	err := json.Unmarshal(data, &response)
	fmt.Println(body)
	return response, err
}

func DispatchResponse(ip string, payload APIRequest, endpoint string, method string) (APIResponse, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return APIResponse{}, err
	}
	buf := bytes.NewBufferString(string(jsonData))
	url := fmt.Sprintf("http://%s:%v/elgato/%s", ip, *PORT, endpoint)
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
	return MarshalResponse(string(responseData))
}

func state(light_ip string) (APIResponse, error) {
	return DispatchResponse(light_ip, APIRequest{}, "lights", "GET")
}

func GetInfo(light_ip string) (APIResponse, error) {
	return DispatchResponse(light_ip, APIRequest{}, "accessory-info", "GET")
}

func GetStatus(light_ip string) int {
	resp, _ := state(light_ip)
	return resp.Lights[0].On
}

func GetTemperature(light_ip string) int {
	resp, _ := state(light_ip)
	return resp.Lights[0].Brightness
}

func GetBrightness(light_ip string) int {
	resp, _ := state(light_ip)
	return resp.Lights[0].Brightness
}

func SetStatusToggle(light_ip string) (APIResponse, error) {
	currentState := GetStatus(light_ip)
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

	return DispatchResponse(light_ip, Request, "lights", "PUT")
}

func SetStatusOff(light_ip string) (APIResponse, error) {
	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On: 0,
		}},
	}
	return DispatchResponse(light_ip, Request, "lights", "PUT")
}

func SetStatusOn(light_ip string) (APIResponse, error) {
	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On: 1,
		}},
	}
	return DispatchResponse(light_ip, Request, "lights", "PUT")
}

func SetBrightnessDiff(light_ip string, diff int) (APIResponse, error) {
	brightness := GetBrightness(light_ip)
	return SetBrightness(light_ip, brightness+diff)
}

func SetTemperatureDiff(light_ip string, diff int) (APIResponse, error) {
	temperature := GetTemperature(light_ip)
	return SetTemperature(light_ip, temperature+diff)
}

func SetTemperature(light_ip string, value int) (APIResponse, error) {
	if value < 143 || value > 344 {
		return APIResponse{}, errors.New(fmt.Sprintf("Needs to be between 143 and 344, was: %v", value))
	}

	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On:          0,
			Temperature: value,
		}},
	}

	return DispatchResponse(light_ip, Request, "lights", "PUT")
}

func AddBrightness(light_ip string, value int) (APIResponse, error) {
	curr := GetBrightness(light_ip)
	return SetBrightness(light_ip, curr+value)
}

func MinusBrightness(light_ip string, value int) (APIResponse, error) {
	curr := GetBrightness(light_ip)
	return SetBrightness(light_ip, curr-value)
}

func AddTemperature(light_ip string, value int) (APIResponse, error) {
	curr := GetTemperature(light_ip)
	return SetTemperature(light_ip, curr+value)
}

func MinusTemperature(light_ip string, value int) (APIResponse, error) {
	curr := GetTemperature(light_ip)
	return SetTemperature(light_ip, curr-value)
}

func SetBrightness(light_ip string, value int) (APIResponse, error) {
	if value < 3 || value > 100 {
		return APIResponse{}, errors.New(fmt.Sprintf("Needs to be between 3 and 100, was: %v", value))
	}

	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On:         0,
			Brightness: value,
		}},
	}

	return DispatchResponse(light_ip, Request, "lights", "PUT")
}

func main() {

	app.HelpFlag.Short('h')
	LIGHTS := app.GetFlag("ips").Strings()
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case on.FullCommand():
		for _, v := range *LIGHTS {
			_, err := SetStatusOn(v)
			if err != nil {
				fmt.Println(err)
			}
		}
	case off.FullCommand():
		for _, v := range *LIGHTS {
			_, err := SetStatusOff(v)
			if err != nil {
				fmt.Println(err)
			}
		}
	case toggle.FullCommand():
		for _, v := range *LIGHTS {
			_, err := SetStatusToggle(v)
			if err != nil {
				fmt.Println(err)
			}
		}
	case info.FullCommand():
		for _, v := range *LIGHTS {
			_, err := GetInfo(v)
			if err != nil {
				fmt.Println(err)
			}
		}
	case status.FullCommand():
		for _, v := range *LIGHTS {
			_, err := state(v)
			if err != nil {
				fmt.Println(err)
			}
		}
	case brightness.FullCommand():
		for _, v := range *LIGHTS {
			if strings.HasPrefix(*destinationValue, "+") {
				newValueString := strings.TrimLeft(*destinationValue, "+")
				newValue, _ := strconv.Atoi(newValueString)
				_, err := AddBrightness(v, newValue)
				if err != nil {
					fmt.Println(err)
				}
			} else if strings.HasPrefix(*destinationValue, "_") {
				newValueString := strings.TrimLeft(*destinationValue, "_")
				newValue, _ := strconv.Atoi(newValueString)
				_, err := MinusBrightness(*destinationValue, newValue)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				newValue, _ := strconv.Atoi(*destinationValue)
				_, err := SetBrightness(v, newValue)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	case temperature.FullCommand():
		for _, v := range *LIGHTS {
			if strings.HasPrefix(*destinationValue, "+") {
				newValueString := strings.TrimLeft(*destinationValue, "+")
				newValue, _ := strconv.Atoi(newValueString)
				_, err := AddTemperature(v, newValue)
				if err != nil {
					fmt.Println(err)
				}
			} else if strings.HasPrefix(*destinationValue, "_") {
				newValueString := strings.TrimLeft(*destinationValue, "_")
				newValue, _ := strconv.Atoi(newValueString)
				_, err := MinusTemperature(*destinationValue, newValue)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				newValue, _ := strconv.Atoi(*destinationValue)
				_, err := SetTemperature(v, newValue)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

	}

}
