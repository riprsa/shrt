package service_test

import (
	"testing"

	"github.com/hararudoka/shrt/service"
)

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
		err  error
	}{
		{
			name: "Empty String URL",
			err:  service.ErrEmptyURL,
		},
		{
			name: "Full HTTP URL",
			url:  "http://user:password@myhost.mydomain/path?query=value#fragment",
			want: "myhost.mydomain/path?query=value#fragment",
		},
		{
			name: "Full HTTPs URL",
			url:  "https://user:password@myhost.mydomain/path?query=value#fragment",
			want: "myhost.mydomain/path?query=value#fragment",
		},
		{
			name: "Full FTP URL",
			url:  "ftp://user:password@myhost.mydomain/path?query=value#fragment",
			want: "myhost.mydomain/path?query=value#fragment",
		},
		// {
		// 	name: "Broken HTTP",
		// 	url:  "http:/user:password@myhost.mydomain/path?query=value#fragment",
		// 	want: "myhost.mydomain/path?query=value#fragment",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			got, err := service.SanitizeURL(tt.url)
			if err != tt.err {
				t.Errorf("SanitizeURL()\nwant error: %v\n got error: %v", tt.err, err)
				return
			}
			if got != tt.want {
				t.Errorf("SanitizeURL()\nwant: %v\n got: %v", tt.want, got)
			}
		})
	}
}
