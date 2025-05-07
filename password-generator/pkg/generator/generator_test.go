package generator

import (
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name        string
		minLength   int
		maxLength   int
		charSet     CharacterSet
		expectError bool
	}{
		{
			name:      "Valid configuration",
			minLength: 8,
			maxLength: 64,
			charSet: CharacterSet{
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
				Special:   true,
			},
			expectError: false,
		},
		{
			name:      "Invalid min length",
			minLength: 0,
			maxLength: 64,
			charSet: CharacterSet{
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
				Special:   true,
			},
			expectError: true,
		},
		{
			name:      "Max length less than min length",
			minLength: 10,
			maxLength: 8,
			charSet: CharacterSet{
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
				Special:   true,
			},
			expectError: true,
		},
		{
			name:      "No character types enabled",
			minLength: 8,
			maxLength: 64,
			charSet: CharacterSet{
				Uppercase: false,
				Lowercase: false,
				Digits:    false,
				Special:   false,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.minLength, tt.maxLength, tt.charSet)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				
				if gen == nil {
					t.Errorf("Expected generator but got nil")
				} else {
					if gen.MinLength != tt.minLength {
						t.Errorf("Expected MinLength %d but got %d", tt.minLength, gen.MinLength)
					}
					if gen.MaxLength != tt.maxLength {
						t.Errorf("Expected MaxLength %d but got %d", tt.maxLength, gen.MaxLength)
					}
					if gen.CharSet != tt.charSet {
						t.Errorf("Expected CharSet %+v but got %+v", tt.charSet, gen.CharSet)
					}
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name      string
		length    int
		charSet   CharacterSet
		minLength int
		maxLength int
	}{
		{
			name:      "All character types",
			length:    12,
			minLength: 8,
			maxLength: 64,
			charSet: CharacterSet{
				Uppercase: true,
				Lowercase: true,
				Digits:    true,
				Special:   true,
			},
		},
		{
			name:      "Only uppercase",
			length:    16,
			minLength: 8,
			maxLength: 64,
			charSet: CharacterSet{
				Uppercase: true,
				Lowercase: false,
				Digits:    false,
				Special:   false,
			},
		},
		{
			name:      "Only lowercase",
			length:    10,
			minLength: 8,
			maxLength: 64,
			charSet: CharacterSet{
				Uppercase: false,
				Lowercase: true,
				Digits:    false,
				Special:   false,
			},
		},
		{
			name:      "Only digits",
			length:    8,
			minLength: 8,
			maxLength: 64,
			charSet: CharacterSet{
				Uppercase: false,
				Lowercase: false,
				Digits:    true,
				Special:   false,
			},
		},
		{
			name:      "Only special",
			length:    20,
			minLength: 8,
			maxLength: 64,
			charSet: CharacterSet{
				Uppercase: false,
				Lowercase: false,
				Digits:    false,
				Special:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.minLength, tt.maxLength, tt.charSet)
			if err != nil {
				t.Fatalf("Failed to create generator: %v", err)
			}
			
			password, err := gen.Generate(tt.length)
			if err != nil {
				t.Fatalf("Failed to generate password: %v", err)
			}
			
			if len(password) != tt.length {
				t.Errorf("Expected password length %d but got %d", tt.length, len(password))
			}
			
			hasUppercase := false
			hasLowercase := false
			hasDigit := false
			hasSpecial := false
			
			for _, char := range password {
				c := string(char)
				if strings.Contains(uppercaseChars, c) {
					hasUppercase = true
				} else if strings.Contains(lowercaseChars, c) {
					hasLowercase = true
				} else if strings.Contains(digitChars, c) {
					hasDigit = true
				} else if strings.Contains(specialChars, c) {
					hasSpecial = true
				}
			}
			
			if tt.charSet.Uppercase && !hasUppercase {
				t.Errorf("Password should contain uppercase characters but doesn't: %s", password)
			}
			if tt.charSet.Lowercase && !hasLowercase {
				t.Errorf("Password should contain lowercase characters but doesn't: %s", password)
			}
			if tt.charSet.Digits && !hasDigit {
				t.Errorf("Password should contain digits but doesn't: %s", password)
			}
			if tt.charSet.Special && !hasSpecial {
				t.Errorf("Password should contain special characters but doesn't: %s", password)
			}
			
			if !tt.charSet.Uppercase && hasUppercase {
				t.Errorf("Password shouldn't contain uppercase characters but does: %s", password)
			}
			if !tt.charSet.Lowercase && hasLowercase {
				t.Errorf("Password shouldn't contain lowercase characters but does: %s", password)
			}
			if !tt.charSet.Digits && hasDigit {
				t.Errorf("Password shouldn't contain digits but does: %s", password)
			}
			if !tt.charSet.Special && hasSpecial {
				t.Errorf("Password shouldn't contain special characters but does: %s", password)
			}
		})
	}
}

func TestGenerateInvalidLength(t *testing.T) {
	gen, err := NewGenerator(8, 64, CharacterSet{
		Uppercase: true,
		Lowercase: true,
		Digits:    true,
		Special:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	
	_, err = gen.Generate(4)
	if err == nil {
		t.Errorf("Expected error for length below minimum but got nil")
	}
	
	_, err = gen.Generate(128)
	if err == nil {
		t.Errorf("Expected error for length above maximum but got nil")
	}
}

func TestPasswordUniqueness(t *testing.T) {
	gen, err := NewGenerator(8, 64, CharacterSet{
		Uppercase: true,
		Lowercase: true,
		Digits:    true,
		Special:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	
	passwords := make(map[string]bool)
	for i := 0; i < 10; i++ {
		password, err := gen.Generate(12)
		if err != nil {
			t.Fatalf("Failed to generate password: %v", err)
		}
		
		if passwords[password] {
			t.Errorf("Generated duplicate password: %s", password)
		}
		
		passwords[password] = true
	}
}
