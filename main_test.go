package main

import "testing"

func Test_soracomcli(t *testing.T) {
	t.Setenv("SORACOM_AUTH_KEY_ID", "authKeyId")
	t.Setenv("SORACOM_AUTH_KEY", "authKey")

	type args struct {
		command string
		body    interface{}
	}
	type body struct {
		Bar string `json:"bar"`
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "with command",
			args:    args{command: "foo"},
			want:    "soracom --auth-key-id authKeyId --auth-key authKey foo",
			wantErr: false,
		},
		{
			name:    "with body",
			args:    args{command: "foo", body: body{Bar: "baz"}},
			want:    `soracom --auth-key-id authKeyId --auth-key authKey foo --body '{"bar":"baz"}'`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := soracomcli(tt.args.command, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("soracomcli() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("soracomcli() = %v, want %v", got, tt.want)
			}
		})
	}
}
