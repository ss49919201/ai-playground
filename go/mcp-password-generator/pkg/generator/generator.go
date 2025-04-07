package generator

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

type CharacterSet struct {
	Uppercase bool
	Lowercase bool
	Digits    bool
	Special   bool
}

type Generator struct {
	MinLength int
	MaxLength int
	CharSet   CharacterSet
}

const (
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	digitChars     = "0123456789"
	specialChars   = "!@#$%^&*()-_=+[]{}|;:,.<>?/"
)

func NewGenerator(minLength, maxLength int, charSet CharacterSet) (*Generator, error) {
	if minLength < 1 {
		return nil, errors.New("minimum length must be at least 1")
	}
	
	if maxLength < minLength {
		return nil, errors.New("maximum length must be greater than or equal to minimum length")
	}
	
	if !charSet.Uppercase && !charSet.Lowercase && !charSet.Digits && !charSet.Special {
		return nil, errors.New("at least one character type must be enabled")
	}
	
	return &Generator{
		MinLength: minLength,
		MaxLength: maxLength,
		CharSet:   charSet,
	}, nil
}

func (g *Generator) Generate(length int) (string, error) {
	if length < g.MinLength || length > g.MaxLength {
		return "", errors.New("requested length is outside allowed range")
	}
	
	var charSet string
	if g.CharSet.Uppercase {
		charSet += uppercaseChars
	}
	if g.CharSet.Lowercase {
		charSet += lowercaseChars
	}
	if g.CharSet.Digits {
		charSet += digitChars
	}
	if g.CharSet.Special {
		charSet += specialChars
	}
	
	var password strings.Builder
	charSetLength := big.NewInt(int64(len(charSet)))
	
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charSetLength)
		if err != nil {
			return "", errors.New("failed to generate secure random number")
		}
		
		password.WriteByte(charSet[randomIndex.Int64()])
	}
	
	if !g.validatePassword(password.String()) {
		return g.Generate(length)
	}
	
	return password.String(), nil
}

func (g *Generator) validatePassword(password string) bool {
	hasUppercase := !g.CharSet.Uppercase
	hasLowercase := !g.CharSet.Lowercase
	hasDigit := !g.CharSet.Digits
	hasSpecial := !g.CharSet.Special
	
	for _, char := range password {
		c := string(char)
		if g.CharSet.Uppercase && strings.Contains(uppercaseChars, c) {
			hasUppercase = true
		} else if g.CharSet.Lowercase && strings.Contains(lowercaseChars, c) {
			hasLowercase = true
		} else if g.CharSet.Digits && strings.Contains(digitChars, c) {
			hasDigit = true
		} else if g.CharSet.Special && strings.Contains(specialChars, c) {
			hasSpecial = true
		}
	}
	
	return hasUppercase && hasLowercase && hasDigit && hasSpecial
}
