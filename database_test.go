package main

import (
	"database/sql"
	"testing"
)

func TestViewsRequest_buildDateTimeFilter(t *testing.T) {
	t.Run("No DateTime filter", func(t *testing.T) {
		request := &ViewsRequest{
			from: "",
			to:   "",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotDateTimeFilter := request.buildDateTimeFilter(namedArgs);
			gotDateTimeFilter != "" ||
				len(*namedArgs) != 0 {
			t.Errorf("ViewsRequest.buildDateTimeFilter(): Wrong return string or length of namedArgs, should be empty")
		}
	})
	t.Run("From filter", func(t *testing.T) {
		request := &ViewsRequest{
			from: "2019-01-01",
			to:   "",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotDateTimeFilter := request.buildDateTimeFilter(namedArgs);
			gotDateTimeFilter != "datetime(time, 'localtime') >= :from" ||
				len(*namedArgs) != 1 ||
				(*namedArgs)[0].Name != "from" ||
				(*namedArgs)[0].Value != "2019-01-01" {
			t.Errorf("ViewsRequest.buildDateTimeFilter(): Wrong return string or namedArgs")
		}
	})
	t.Run("To filter", func(t *testing.T) {
		request := &ViewsRequest{
			from: "",
			to:   "2019-01-01",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotDateTimeFilter := request.buildDateTimeFilter(namedArgs);
			gotDateTimeFilter != "datetime(time, 'localtime') <= :to" ||
				len(*namedArgs) != 1 ||
				(*namedArgs)[0].Name != "to" ||
				(*namedArgs)[0].Value != "2019-01-01" {
			t.Errorf("ViewsRequest.buildDateTimeFilter(): Wrong return string or namedArgs")
		}
	})
	t.Run("From & To filter", func(t *testing.T) {
		request := &ViewsRequest{
			from: "2018-01-01",
			to:   "2019-01-01",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotDateTimeFilter := request.buildDateTimeFilter(namedArgs);
			gotDateTimeFilter != "datetime(time, 'localtime') between :from and :to" ||
				len(*namedArgs) != 2 ||
				(*namedArgs)[0].Name != "from" ||
				(*namedArgs)[0].Value != "2018-01-01" ||
				(*namedArgs)[1].Name != "to" ||
				(*namedArgs)[1].Value != "2019-01-01" {
			t.Errorf("ViewsRequest.buildDateTimeFilter(): Wrong return string or namedArgs")
		}
	})
}

func TestViewsRequest_buildUrlFilter(t *testing.T) {
	t.Run("No url filter", func(t *testing.T) {
		request := &ViewsRequest{
			url: "",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotUrlFilter := request.buildUrlFilter(namedArgs);
			gotUrlFilter != "" ||
				len(*namedArgs) != 0 {
			t.Errorf("ViewsRequest.buildUrlFilter(): Wrong return string or length of namedArgs, should be empty")
		}
	})
	t.Run("Url filter", func(t *testing.T) {
		request := &ViewsRequest{
			url: "google",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotUrlFilter := request.buildUrlFilter(namedArgs);
			gotUrlFilter != "url like :url" ||
				len(*namedArgs) != 1 ||
				(*namedArgs)[0].Name != "url" ||
				(*namedArgs)[0].Value != "%google%" {
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
		if gotRefFilter := request.buildRefFilter(namedArgs);
			gotRefFilter != "" ||
				len(*namedArgs) != 0 {
			t.Errorf("ViewsRequest.buildRefFilter(): Wrong return string or length of namedArgs, should be empty")
		}
	})
	t.Run("Ref filter", func(t *testing.T) {
		request := &ViewsRequest{
			ref: "google",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotRefFilter := request.buildRefFilter(namedArgs);
			gotRefFilter != "ref like :ref" ||
				len(*namedArgs) != 1 ||
				(*namedArgs)[0].Name != "ref" ||
				(*namedArgs)[0].Value != "%google%" {
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
		if gotUseragentFilter := request.buildUseragentFilter(namedArgs);
			gotUseragentFilter != "" ||
				len(*namedArgs) != 0 {
			t.Errorf("ViewsRequest.buildUseragentFilter(): Wrong return string or length of namedArgs, should be empty")
		}
	})
	t.Run("Useragent filter", func(t *testing.T) {
		request := &ViewsRequest{
			ua: "Firefox",
		}
		namedArgs := &[]sql.NamedArg{}
		if gotUseragentFilter := request.buildUseragentFilter(namedArgs);
			gotUseragentFilter != "useragent like :ua" ||
				len(*namedArgs) != 1 ||
				(*namedArgs)[0].Name != "ua" ||
				(*namedArgs)[0].Value != "%Firefox%" {
			t.Errorf("ViewsRequest.buildUseragentFilter(): Wrong return string or namedArgs")
		}
	})
}
