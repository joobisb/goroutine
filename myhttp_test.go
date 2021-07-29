package main

import (
	"reflect"
	"testing"
	"time"
)

const googleURI = "https://google.com"
const facebookURI = "https://facebook.com"

var urlResponseExpectedMap = map[string]string{
	googleURI:   "123",
	facebookURI: "456",
}

type MockHTTPClient struct{}

func (MockHTTPClient) GetResponseHash(uri string, timeout time.Duration) (string, error) {
	switch uri {
	case googleURI:
		return urlResponseExpectedMap[googleURI], nil
	case facebookURI:
		return urlResponseExpectedMap[facebookURI], nil
	}
	return "", nil
}

func TestExecuteSuccessCase(t *testing.T) {
	testURIs := []string{
		googleURI,
		facebookURI,
	}

	actualReponseMap, err := execute(MockHTTPClient{}, testURIs, 1)
	if err != nil {
		t.Errorf("Error while execute: %v", err)
	}

	if !reflect.DeepEqual(urlResponseExpectedMap, actualReponseMap) {
		t.Errorf("Unexpected result expected: %v, actual: %v", urlResponseExpectedMap, actualReponseMap)
	}
}

func TestExecuteFailureCase(t *testing.T) {
	testURIs := []string{
		"fakeURL",
	}

	actualReponseMap, err := execute(MockHTTPClient{}, testURIs, 1)
	if err != nil {
		t.Errorf("Error while execute: %v", err)
	}

	if reflect.DeepEqual(urlResponseExpectedMap, actualReponseMap) {
		t.Errorf("Unexpected result expected: %v, actual: %v", urlResponseExpectedMap, actualReponseMap)
	}
}
