package util

import (
	"os/exec"
	"reflect"
	"testing"
)

func Test_openInBrowser(t *testing.T) {
	unsupportedOs := "this is not an os"
	supportedOs := "linux"

	type args struct {
		os  string
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"no error for supported os",
			args{supportedOs, "url"},
			false,
		},
		{
			"error for unsupported os",
			args{unsupportedOs, "url"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := openInBrowser(tt.args.os, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("OpenInBrowser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getBrowserCmd(t *testing.T) {
	type args struct {
		os  string
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    *exec.Cmd
		wantErr bool
	}{
		{
			"xdg-open cmd on linux",
			args{
				"linux",
				"url",
			},
			exec.Command("xdg-open", "url"),
			false,
		},
		{
			"rundll32 cmd on windows",
			args{
				"windows",
				"url",
			},
			exec.Command("rundll32", "url.dll,FileProtocolHandler", "url"),
			false,
		},
		{
			"open cmd on mac",
			args{
				"darwin",
				"url",
			},
			exec.Command("open", "url"),
			false,
		},
		{
			"erros on solaris",
			args{
				"solaris",
				"url",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBrowserCmd(tt.args.os, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBrowserCmd(%v, %v) = %v, want %v", tt.args.os, tt.args.url, got, tt.want)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBrowserCmd(%v, %v) = %v, want %v", tt.args.os, tt.args.url, got, tt.want)
			}
		})
	}
}
