package main

import (
	"testing"
)

func Test_config_statsAuth(t *testing.T) {
	type fields struct {
		StatsUsername string
		StatsPassword string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"No username nor password", fields{"", ""}, false},
		{"Only username", fields{"abc", ""}, false},
		{"Only password", fields{"", "abc"}, false},
		{"Username and password", fields{"abc", "abc"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ac := &config{
				StatsUsername: tt.fields.StatsUsername,
				StatsPassword: tt.fields.StatsPassword,
			}
			if got := ac.statsAuth(); got != tt.want {
				t.Errorf("config.statsAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
