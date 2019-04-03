package main

import (
	"os"
	"testing"
)

func Test_port(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		want   string
	}{
		{name: "default", envVar: "", want: "8080"},
		{name: "custom", envVar: "1234", want: "1234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("PORT", tt.envVar)
			if got := port(); got != tt.want {
				t.Errorf("port() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dnt(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		want   bool
	}{
		{name: "default", envVar: "", want: true},
		{envVar: "true", want: true},
		{envVar: "t", want: true},
		{envVar: "TRUE", want: true},
		{envVar: "1", want: true},
		{envVar: "false", want: false},
		{envVar: "f", want: false},
		{envVar: "0", want: false},
		{envVar: "abc", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("DNT", tt.envVar)
			if got := dnt(); got != tt.want {
				t.Errorf("dnt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dbPath(t *testing.T) {
	tests := []struct {
		name       string
		envVar     string
		wantDbPath string
	}{
		{name: "default", envVar: "", wantDbPath: "data/kis3.db"},
		{envVar: "kis3.db", wantDbPath: "kis3.db"},
		{envVar: "data.db", wantDbPath: "data.db"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("DB_PATH", tt.envVar)
			if gotDbPath := dbPath(); gotDbPath != tt.wantDbPath {
				t.Errorf("dbPath() = %v, want %v", gotDbPath, tt.wantDbPath)
			}
		})
	}
}

func Test_statsUsername(t *testing.T) {
	tests := []struct {
		name         string
		envVar       string
		wantUsername string
	}{
		{name: "default", envVar: "", wantUsername: ""},
		{envVar: "abc", wantUsername: "abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("STATS_USERNAME", tt.envVar)
			if gotUsername := statsUsername(); gotUsername != tt.wantUsername {
				t.Errorf("statsUsername() = %v, want %v", gotUsername, tt.wantUsername)
			}
		})
	}
}

func Test_statsPassword(t *testing.T) {
	tests := []struct {
		name         string
		envVar       string
		wantPassword string
	}{
		{name: "default", envVar: "", wantPassword: ""},
		{envVar: "def", wantPassword: "def"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("STATS_PASSWORD", tt.envVar)
			if gotPassword := statsPassword(); gotPassword != tt.wantPassword {
				t.Errorf("statsPassword() = %v, want %v", gotPassword, tt.wantPassword)
			}
		})
	}
}

func Test_statsAuth(t *testing.T) {
	type args struct {
		ac *config
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "default", args: struct{ ac *config }{ac: &config{}}, want: false},
		{name: "only username set", args: struct{ ac *config }{ac: &config{statsUsername: "abc"}}, want: false},
		{name: "only password set", args: struct{ ac *config }{ac: &config{statsPassword: "def"}}, want: false},
		{name: "username and password set", args: struct{ ac *config }{ac: &config{statsUsername: "abc", statsPassword: "def"}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statsAuth(tt.args.ac); got != tt.want {
				t.Errorf("statsAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
