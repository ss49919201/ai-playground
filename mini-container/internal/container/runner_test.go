package container

import (
	"testing"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name    string
		program string
		args    []string
		wantErr bool
	}{
		{
			name:    "echo command",
			program: "/bin/echo",
			args:    []string{"hello"},
			wantErr: false,
		},
		{
			name:    "true command",
			program: "/usr/bin/true",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "false command",
			program: "/usr/bin/false",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "non-existent command",
			program: "/bin/nonexistent",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RunCommand(tt.program, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunCommandOutput(t *testing.T) {
	output, err := RunCommandWithOutput("/bin/echo", []string{"test"})
	if err != nil {
		t.Errorf("RunCommandWithOutput() error = %v", err)
	}
	expected := "test\n"
	if output != expected {
		t.Errorf("RunCommandWithOutput() output = %v, want %v", output, expected)
	}
}