package dtls

import (
	"reflect"
	"testing"
)

func TestChangeCipherSpecRoundTrip(t *testing.T) {
	c := changeCipherSpec{}
	raw, err := c.marshal()
	if err != nil {
		t.Error(err)
	}

	var cNew changeCipherSpec
	if err := cNew.unmarshal(raw); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(c, cNew) {
		t.Errorf("ChangeCipherSpec round trip: got %#v, want %#v", cNew, c)
	}
}

func TestChangeCipherSpecInvalid(t *testing.T) {
	c := changeCipherSpec{}
	if err := c.unmarshal([]byte{0x00}); err != errInvalidCipherSpec {
		t.Errorf("ChangeCipherSpec invalid assert: got %#v, want %#v", err, errInvalidCipherSpec)
	}

}
