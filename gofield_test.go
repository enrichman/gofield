package gofield

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Person struct {
	Name       string         `json:"name,omitempty"`
	Age        int            `json:"age,omitempty"`
	Address    *Address       `json:"addr,omitempty"`
	AddressNop Address        `json:"addrNop,omitempty"`
	Email      []*MailAddress `json:"email,omitempty"`
}

type Address struct {
	City     string `json:"city_name,omitempty"`
	PostCode string `json:"code,omitempty"`
}

type MailAddress struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func Test_Split(t *testing.T) {
	tests := []struct {
		name           string
		fields         string
		expectedFields []string
	}{
		{
			"Empty split",
			"foo",
			[]string{"foo"},
		},
		{
			"Simple split",
			"foo,bar,foobar",
			[]string{"foo", "bar", "foobar"},
		},
		{
			"Combined split",
			"foo,bar{innerbar},foobar",
			[]string{"foo", "bar{innerbar}", "foobar"},
		},
		{
			"Combined multiple split",
			"foo,bar{innerbar,innerfoo,innerfoobar},foobar",
			[]string{"foo", "bar{innerbar,innerfoo,innerfoobar}", "foobar"},
		},
		{
			"Combined multiple split multiple values",
			"foo,bar{innerbar,innerfoo,innerfoobar},foobar{innerbar2}",
			[]string{"foo", "bar{innerbar,innerfoo,innerfoobar}", "foobar{innerbar2}"},
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		resultObj := Split(tt.fields, ",")
		assert.Equal(t, tt.expectedFields, resultObj, "Unexpected fields")
	}
}

func Test_SimpleObject(t *testing.T) {

	tests := []struct {
		name        string
		obj         interface{}
		expectedObj interface{}
		fields      string
	}{
		{
			"No field",
			Person{
				Name: "Enrico",
				Age:  27,
			},
			Person{
				Name: "Enrico",
				Age:  27,
			},
			"",
		},
		{
			"Missing field",
			Person{
				Name: "Enrico",
				Age:  27,
			},
			Person{},
			"foobar",
		},
		{
			"One field",
			Person{
				Name: "Enrico",
				Age:  27,
			},
			Person{Name: "Enrico"},
			"name",
		},
		{
			"Two fields, one missing",
			Person{
				Name: "Enrico",
				Age:  27,
			},
			Person{Name: "Enrico"},
			"name,foobar",
		},
		{
			"Two fields",
			Person{
				Name: "Enrico",
				Age:  27,
			},
			Person{
				Name: "Enrico",
				Age:  27,
			},
			"name,age",
		},
		{
			"Inner field",
			Person{
				Name: "Enrico",
				Age:  27,
				Address: &Address{
					City:     "Rome",
					PostCode: "00155",
				},
				AddressNop: Address{
					City:     "RomeNop",
					PostCode: "00155",
				},
			},
			Person{
				Name: "Enrico",
				Address: &Address{
					City: "Rome",
				},
				AddressNop: Address{
					City:     "RomeNop",
					PostCode: "00155",
				},
			},
			"name,addr{city_name},addrNop",
		},
		{
			"Inner multiple field",
			Person{
				Name: "Enrico",
				Age:  27,
				Address: &Address{
					City:     "Rome",
					PostCode: "00155",
				},
				AddressNop: Address{
					City:     "RomeNop",
					PostCode: "00155",
				},
			},
			Person{
				Name: "Enrico",
				Address: &Address{
					City:     "Rome",
					PostCode: "00155",
				},
				AddressNop: Address{
					City:     "RomeNop",
					PostCode: "00155",
				},
			},
			"name,addr{city_name,code},addrNop",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		resultObj := Reduce(tt.obj, tt.fields)

		b1, _ := json.Marshal(resultObj)
		b2, _ := json.Marshal(tt.expectedObj)
		t.Log("Unwinded:", string(b1))
		t.Log("Expected:", string(b2))

		var decoded Person
		_ = json.Unmarshal(b1, &decoded)

		assert.True(t, reflect.DeepEqual(decoded, tt.expectedObj))
	}
}

func TestArr(t *testing.T) {
	tests := []struct {
		name        string
		obj         interface{}
		expectedObj interface{}
		fields      string
	}{
		{
			"No field",
			[]*MailAddress{
				&MailAddress{"personal", "foo@email.com"},
				&MailAddress{"work", "bar@email.com"},
			},
			[]*MailAddress{
				&MailAddress{"personal", "foo@email.com"},
				&MailAddress{"work", "bar@email.com"},
			},
			"",
		},
		{
			"Missing field",
			[]*MailAddress{
				&MailAddress{"personal", "foo@email.com"},
				&MailAddress{"work", "bar@email.com"},
			},
			[]*MailAddress{
				&MailAddress{},
				&MailAddress{},
			},
			"foobar",
		},
		{
			"One field",
			[]*MailAddress{
				&MailAddress{"personal", "foo@email.com"},
				&MailAddress{"work", "bar@email.com"},
			},
			[]*MailAddress{
				&MailAddress{Name: "personal"},
				&MailAddress{Name: "work"},
			},
			"name",
		},
		{
			"Two fields",
			[]*MailAddress{
				&MailAddress{"personal", "foo@email.com"},
				&MailAddress{"work", "bar@email.com"},
			},
			[]*MailAddress{
				&MailAddress{"personal", "foo@email.com"},
				&MailAddress{"work", "bar@email.com"},
			},
			"name,email",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		resultObj := Reduce(tt.obj, tt.fields)

		b1, _ := json.Marshal(resultObj)
		b2, _ := json.Marshal(tt.expectedObj)
		t.Log("Unwinded:", string(b1))
		t.Log("Expected:", string(b2))

		var decoded1 interface{}
		_ = json.Unmarshal(b1, &decoded1)
		var decoded2 interface{}
		_ = json.Unmarshal(b2, &decoded2)

		assert.True(t, reflect.DeepEqual(decoded1, decoded2))
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name        string
		obj         map[string]interface{}
		expectedObj map[string]interface{}
		fields      string
	}{
		{
			"No field",
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
			},
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
			},
			"",
		},
		{
			"Missing field",
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
			},
			map[string]interface{}{},
			"foobar",
		},
		{
			"One field",
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
			},
			map[string]interface{}{
				"name": "Enrico",
			},
			"name",
		},
		{
			"Two fields",
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
			},
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
			},
			"name,age",
		},
		{
			"Inner fields",
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
				"address": map[string]interface{}{
					"city":      "Rome",
					"post_code": "00123",
				},
			},
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
				"address": map[string]interface{}{
					"city": "Rome",
				},
			},
			"name,age,address{city}",
		},
		{
			"Inner multiple fields",
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
				"address": map[string]interface{}{
					"city":      "Rome",
					"post_code": "00123",
				},
			},
			map[string]interface{}{
				"name": "Enrico",
				"age":  27,
				"address": map[string]interface{}{
					"city":      "Rome",
					"post_code": "00123",
				},
			},
			"name,age,address{city,post_code}",
		},
	}

	for _, tt := range tests {
		t.Log(tt.name)
		resultObj := Reduce(tt.obj, tt.fields)

		b1, _ := json.Marshal(resultObj)
		b2, _ := json.Marshal(tt.expectedObj)
		t.Log("Unwinded:", string(b1))
		t.Log("Expected:", string(b2))

		var decoded1 interface{}
		_ = json.Unmarshal(b1, &decoded1)
		var decoded2 interface{}
		_ = json.Unmarshal(b2, &decoded2)

		assert.True(t, reflect.DeepEqual(decoded1, decoded2))
	}
}
