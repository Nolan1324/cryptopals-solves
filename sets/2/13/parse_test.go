package main

import (
	"net/url"
	"testing"
)

func TestParseQuery(t *testing.T) {
	vals, err := url.ParseQuery("foo=bar&baz=qux&zap=zazzle")
	if err != nil {
		t.Fatal(err)
	}

	if vals.Get("foo") != "bar" {
		t.Error()
	}
	if vals.Get("baz") != "qux" {
		t.Error()
	}
	if vals.Get("zap") != "zazzle" {
		t.Error()
	}
}

func TestParseQueryDuplicates(t *testing.T) {
	vals, err := url.ParseQuery("foo=bar&foo=baz")
	if err != nil {
		t.Fatal(err)
	}

	if vals.Get("foo") != "bar" {
		t.Error()
	}
}
