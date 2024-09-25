package api

import (
	"errors"
	"testing"
)

func TestFindLocation(t *testing.T) {
	var cases = []struct {
		IP   string
		resp Response
		err  error
	}{
		{"1", Response{}, ErrInvalidQuery},
		{"192.168.1.1", Response{}, ErrPrivateRange},
		{"0.0.0.0", Response{}, ErrReservedRange},
		{"1.1.1.1", Response{"success", "", "Australia", "Queensland", "South Brisbane"}, nil},
		{"8.8.8.8", Response{"success", "", "United States", "Virginia", "Ashburn"}, nil},
	}

	for i, c := range cases {
		resp, err := FindLocation(c.IP)
		t.Log(err)
		if !errors.Is(err, c.err) {
			t.Errorf("[%d] excepted: %v, got: %v", i+1, c.err, err)
		}
		if resp != c.resp {
			t.Errorf("[%d] excepted: %+v, got: %+v", i+1, resp, c.resp)
		}
	}
}
