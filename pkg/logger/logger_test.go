package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		logFIle string
		debug  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "debug creation success",
			args: args{"log.log", true},
			wantErr: false,
		},
		{
			name: "nondebug creation success",
			args: args{"log.log", false},
			wantErr: false,
		},
		{
			name: "creation error: invalid file",
			args: args{"log/", true},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.logFIle, tt.args.debug)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
