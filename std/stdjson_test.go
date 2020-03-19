package std

import (
	"github.com/anssihalmeaho/funl"
	"testing"
)

var jsonBlob = []byte(`[
	{"Name": "Platypus", "Order": "Monotremata"},
	{"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)

var jsonBlobFail = []byte(`[
	{"Name": Platypus", "Order": "Monotremata"},
	{"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)

var jsonX = []byte(` [true, false, ["tjaa", 123]] `)

func TestEncodeOK(t *testing.T) {
	//inValue := funl.Value{Kind: funl.StringValue, Data: "some stuff"}
	//inValue := funl.Value{Kind: funl.IntValue, Data: 123}
	//inValue := funl.Value{Kind: funl.BoolValue, Data: true}
	//inValue := funl.Value{Kind: funl.BoolValue, Data: true}
	//inValue := funl.Value{Kind: funl.OpaqueValue, Data: &OpaqueJSONnull{}}
	inValue := funl.Value{Kind: funl.FloatValue, Data: 0.5}
	ok, errText, val := encodeJSON("dumname", nil, inValue)

	t.Logf("data = %s", val.Data.(*OpaqueByteArray).data)
	if !ok {
		t.Errorf("Not ok: %s", errText)
	}
}

func TestDecodeOK(t *testing.T) {
	ok, errText, val := decodeJSON("dumname", nil, jsonBlob)
	if !ok {
		t.Logf("error text = %s", errText)
		t.Errorf("Should be ok")
	}
	if !ok {
		t.Errorf("Not ok: %s", errText)
	}
	if val.Kind != funl.ListValue {
		t.Errorf("Not list")
	}
}

func TestDecodeFail(t *testing.T) {
	ok, errText, _ := decodeJSON("dumname", nil, jsonBlobFail)
	if ok {
		t.Errorf("Should fail")
	}
	if errText == "" {
		t.Errorf("No error text")
	}
}
