package main

import (
	"database/sql"
	"testing"
)

func TestViewsRequest_buildUrlFilter(t *testing.T) {
	t.Run("No url filter", func(t *testing.T) {
		request := &ViewsRequest{
			url: "",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotUrlFilter := request.buildUrlFilter(namedArgs); gotUrlFilter != "" || len(*namedArgs) != 0 {
			t.Errorf("ViewsRequest.buildUrlFilter(): Wrong return string or length of namedArgs, should be empty")
		}
	})
	t.Run("Url filter", func(t *testing.T) {
		request := &ViewsRequest{
			url: "google",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotUrlFilter := request.buildUrlFilter(namedArgs); gotUrlFilter != "url like :url" || len(*namedArgs) != 1 || (*namedArgs)[0].Name != "url" || (*namedArgs)[0].Value != "%google%" {
			t.Errorf("ViewsRequest.buildUrlFilter(): Wrong return string or namedArgs")
		}
	})
}

func TestViewsRequest_buildRefFilter(t *testing.T) {
	t.Run("No ref filter", func(t *testing.T) {
		request := &ViewsRequest{
			ref: "",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotRefFilter := request.buildRefFilter(namedArgs); gotRefFilter != "" || len(*namedArgs) != 0 {
			t.Errorf("ViewsRequest.buildRefFilter(): Wrong return string or length of namedArgs, should be empty")
		}
	})
	t.Run("Ref filter", func(t *testing.T) {
		request := &ViewsRequest{
			ref: "google",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotRefFilter := request.buildRefFilter(namedArgs); gotRefFilter != "ref like :ref" || len(*namedArgs) != 1 || (*namedArgs)[0].Name != "ref" || (*namedArgs)[0].Value != "%google%" {
			t.Errorf("ViewsRequest.buildRefFilter(): Wrong return string or namedArgs")
		}
	})
}

func TestViewsRequest_buildUseragentFilter(t *testing.T) {
	t.Run("No useragent filter", func(t *testing.T) {
		request := &ViewsRequest{
			ua: "",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotUseragentFilter := request.buildUseragentFilter(namedArgs); gotUseragentFilter != "" || len(*namedArgs) != 0 {
			t.Errorf("ViewsRequest.buildUseragentFilter(): Wrong return string or length of namedArgs, should be empty")
		}
	})
	t.Run("Useragent filter", func(t *testing.T) {
		request := &ViewsRequest{
			ua: "Firefox",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotUseragentFilter := request.buildUseragentFilter(namedArgs); gotUseragentFilter != "useragent like :ua" || len(*namedArgs) != 1 || (*namedArgs)[0].Name != "ua" || (*namedArgs)[0].Value != "%Firefox%" {
			t.Errorf("ViewsRequest.buildUseragentFilter(): Wrong return string or namedArgs")
		}
	})
}
