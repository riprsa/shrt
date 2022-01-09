package validation

import "testing"

func TestURLValidation(t *testing.T){
	v := Validator{}

	testStings := []string{"vk.com"}

	for _, e := range testStings {
		v.URLValidation(e)
	}
}