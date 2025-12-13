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

func Test_createStringifiedTable(t *testing.T) {
	tests := []struct {
		name    string
		title   []string
		content [][]string
		want    string
		wantErr bool
	}{
		{
			name:    "should return formatted table",
			title:   []string{"ID", "First Name"},
			content: [][]string{{"1", "Jarvis"}},
			want:    "|ID||First Name|\n1     Jarvis   \n",
			wantErr: false,
		},
		{
			name:    "should correctly space columns with padding when content is bigger that title",
			title:   []string{"ID", "report"},
			content: [][]string{{"uuid-1234-5678", "feeling good and alive"}},
			want:    "|------ID------||--------report--------|\nuuid-1234-5678  feeling good and alive \n",
			wantErr: false,
		},
		{
			name:    "should correctly detect the longest content and space columns with padding",
			title:   []string{"ID", "report"},
			content: [][]string{{"1", "status is nominal"}, {"uuid-1234", "gud"}},
			want:    "|---ID----||-----report------|\n    1      status is nominal \n uuid-1234         gud        \n",
			wantErr: false,
		},
		{
			name:    "should give an error when title has less columns than content",
			title:   []string{"ID"},
			content: [][]string{{"1", "another column"}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "should return empty string when no content were given",
			title:   []string{"ID"},
			content: [][]string{},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := createStringifiedTable(tt.title, tt.content)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("formatContentWithTitle() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("formatContentWithTitle() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("formatContentWithTitle() = \n%v \nwant \n%v", got, tt.want)
			}
		})
	}
}

func Test_assembleFormattedString(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		length  int
		border  string
		padding string
		want    string
	}{
		{
			name:    "should assemle string with border and padding",
			str:     "test",
			length:  12,
			border:  "$",
			padding: ".",
			want:    "$...test...$",
		},
		{
			name:    "should return original string when given length is not enough to assemble",
			str:     "test",
			length:  5,
			border:  "$",
			padding: ".",
			want:    "test",
		},
		{
			name:    "should return original string when given more than 1 rune padding",
			str:     "test",
			length:  13,
			border:  "$",
			padding: ".-.",
			want:    "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assembleFormattedString(tt.str, tt.length, tt.border, tt.padding)
			if got != tt.want {
				t.Errorf("assembleFormattedString() = \n%v\n want\n%v\n", got, tt.want)
			}
		})
	}
}
