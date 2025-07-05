package container

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Command
		wantErr  bool
	}{
		{
			name:     "no arguments",
			args:     []string{},
			expected: Command{Action: "help"},
			wantErr:  false,
		},
		{
			name:     "help command",
			args:     []string{"help"},
			expected: Command{Action: "help"},
			wantErr:  false,
		},
		{
			name:     "run command",
			args:     []string{"run", "/bin/echo", "hello"},
			expected: Command{Action: "run", Program: "/bin/echo", Args: []string{"hello"}},
			wantErr:  false,
		},
		{
			name:     "invalid command",
			args:     []string{"invalid"},
			expected: Command{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Action != tt.expected.Action {
				t.Errorf("ParseArgs() Action = %v, want %v", got.Action, tt.expected.Action)
			}
			if got.Program != tt.expected.Program {
				t.Errorf("ParseArgs() Program = %v, want %v", got.Program, tt.expected.Program)
			}
		})
	}
}