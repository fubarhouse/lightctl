package lib

import "testing"

func TestOne(t *testing.T) {

	Request := APIRequest{
		NumberOfLights: 1,
		Lights: []Light{{
			On: 0,
		}},
	}

	resp, err := DispatchResponse("0.0.0.0", Request, "8000", "PUT")
	t.Log(resp, err)
}
