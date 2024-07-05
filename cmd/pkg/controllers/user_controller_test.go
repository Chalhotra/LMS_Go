package controllers

import (
	"testing"
)

type payload struct {
	hashedPassword  string
	plainPassword   string
	expectedSuccess bool
}

var testcases = []payload{{
	hashedPassword:  "$2a$10$qLiCv2cA6/Reqcw0iKEdpOLlXygS38jULmX5BgRhSelSXRldms4Sq",
	plainPassword:   "pintu",
	expectedSuccess: true,
}, {
	hashedPassword:  "$2a$10$qLiCv2cA6/Reqcw0iKEdpOLlXygS38jULmX5BgRhSelSXRldms4Sq",
	plainPassword:   "somerandomstring",
	expectedSuccess: false,
}}

func TestCheckPassword(t *testing.T) {

	for _, tc := range testcases {
		hp := tc.hashedPassword
		pp := tc.plainPassword
		if checkPassword(hp, pp) != tc.expectedSuccess {
			t.Errorf("expected %v got %v", tc.expectedSuccess, checkPassword(hp, pp))
		}

	}
}
