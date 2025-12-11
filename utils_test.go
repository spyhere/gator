package main

import "testing"

func Test_extractCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    command
		wantErr bool
	}{
		{
			name:    "Should return error when no commands were passed",
			args:    []string{"."},
			want:    command{},
			wantErr: true,
		},
		{
			name:    "Should return requested command with 0 arguments",
			args:    []string{".", "reset"},
			want:    command{name: "reset", args: []string{}},
			wantErr: false,
		},
		{
			name:    "Should return requested command with 1 argument",
			args:    []string{".", "login", "test_user"},
			want:    command{name: "login", args: []string{"test_user"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := extractCommand(tt.args)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("extractCommand() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("extractCommand() succeeded unexpectedly")
			}
			if got.name != tt.want.name {
				t.Errorf("extractCommand() = %v, want %v command", got.name, tt.want.name)
			}
			if len(got.args) != len(tt.want.args) {
				t.Errorf("extractCommand() = args len %v, want %v", len(got.args), len(tt.want.args))
			}
			for idx, it := range got.args {
				if it != tt.want.args[idx] {
					t.Errorf("extractCommand() = unexpected arg %v, want %v", it, tt.want.args[idx])
				}
			}
		})
	}
}
