package claim

import (
	"errors"
	"testing"
)

// func TestValidate(t *testing.T) {
// 	c := Claim{
// 		Subject:    "",
// 		Namespaces: []string{"default"},
// 	}
// 	err := c.Validate()
//
// 	if !errors.Is(err, ErrorRequireSubject) {
// 		t.Errorf("expecting  ErrorRequireNamespace")
// 	}
// }

func TestValidateError(t *testing.T) {
	var tests = []struct {
		name     string
		expected error
		given    Claim
	}{
		{"Expected raise ErrorRequireNamespace", ErrorRequireNamespace, Claim{Namespaces: []string{}, Subject: "subject"}},
		{"Expected raise ErrorRequireSubject", ErrorRequireSubject, Claim{Namespaces: []string{""}, Subject: ""}},
		//TODO: うまくこのテストがかけない
		// {"Expected raise ErrorRequireSubject", &ClaimValidationError{}, Claim{Namespaces: []string{"invalid@namespace"}, Subject: "subject"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.given.Validate()
			if !errors.Is(err, tt.expected) {
				t.Errorf("expected %s\n actual %v", tt.expected, err)
			}
		})
	}
}

func TestNormalizeDNS1123(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    string
	}{
		{"scott-tiger unnecessary normalize", "scott-tiger", "scott-tiger"},
		{"replace normalize target", "scott-tiger", "scott.tiger"},
		{"replace normalize target", "scott-tiger", "scott@tiger"},
		{"replace normalize target", "scott-tiger", "scott/tiger"},
		{"replace normalize target", "scott-tiger", "scott_tiger"},
		{"replace normalize target", "scotttiger", "ScottTiger"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NormalizeDNS1123(tt.given)
			if actual != tt.expected {
				t.Logf("%v", err)
				t.Errorf("given(%s): expected [%s], actual [%s]", tt.given, tt.expected, actual)
			}
		})
	}
}
