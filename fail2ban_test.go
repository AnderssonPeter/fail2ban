package fail2ban_test

import (
	"errors"
	"testing"

	plug "github.com/tommoulard/fail2ban"
)

func TestDummy(t *testing.T) {
	cfg := plug.CreateConfig()
	t.Log(cfg)
}

func TestImportIp(t *testing.T) {
	tests := []struct {
		name    string
		list    plug.List
		strWant []string
		err     error
	}{
		{
			name: "empty list",
			list: plug.List{
				Ip:    []string{},
				Files: []string{},
			},
			strWant: []string{},
			err:     nil,
		},

		{
			name: "simple import",
			list: plug.List{
				Ip:    []string{"192.168.0.0", "0.0.0.0", "255.255.255.255"},
				Files: []string{"tests/test-ipfile.txt"},
			},
			strWant: []string{"192.168.0.0", "255.0.0.0", "42.42.42.42", "13.38.70.00", "192.168.0.0", "0.0.0.0", "255.255.255.255"},
			err:     nil,
		},

		{
			name: "import only file",
			list: plug.List{
				Ip:    []string{},
				Files: []string{"tests/test-ipfile.txt"},
			},
			strWant: []string{"192.168.0.0", "255.0.0.0", "42.42.42.42", "13.38.70.00"},
			err:     nil,
		},

		{
			name: "import only ip",
			list: plug.List{
				Ip:    []string{"192.168.0.0", "0.0.0.0", "255.255.255.255"},
				Files: []string{},
			},
			strWant: []string{"192.168.0.0", "0.0.0.0", "255.255.255.255"},
			err:     nil,
		},

		{
			name: "import no file",
			list: plug.List{
				Ip:    []string{},
				Files: []string{"tests/idontexist.txt"},
			},
			strWant: []string{},
			err:     errors.New("open tests/idontexist.txt: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, e := plug.ImportIP(tt.list)
			t.Logf("%+v", got)
			if e != nil && e.Error() != tt.err.Error() {
				t.Errorf("wanted '%s' got '%s'", tt.err, e)
			}
			if len(got) != len(tt.strWant) {
				t.Errorf("wanted '%d' got '%d'", len(tt.strWant), len(got))
			}

			for i, elt := range tt.strWant {
				if got[i] != elt {
					t.Errorf("wanted '%s' got '%s'", elt, got[i])
				}
			}
		})
	}
}
