package utils

import (
	"testing"
)

func TestHashAndSalt(t *testing.T) {
	type args struct {
		pass []byte
	}
	tests := []struct {
		name  string
		args  args
		empty bool
	}{
		{
			name:  "hash successfully",
			args:  args{pass: []byte("test")},
			empty: false,
		},
		{
			name:  "hash too long",
			args:  args{pass: []byte("01234567890123456789012345678901234567890123456789012345678901234567890123456789")},
			empty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashAndSalt(tt.args.pass)
			if (got == "") != tt.empty {
				t.Errorf("HashAndSalt() = %v, args %v", got, tt.args)
			}
		})
	}
}
