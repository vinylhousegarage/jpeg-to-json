package apierror

import (
	"errors"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	originalErr := errors.New("base error")
	internalInfo := "debug info"

	tests := []struct {
		name         string
		code         ErrorCode
		status       int
		err          error
		internalArgs []string
		wantInternal string
	}{
		{
			name:         "All arguments provided",
			code:         "TEST_CODE",
			status:       http.StatusInternalServerError,
			err:          originalErr,
			internalArgs: []string{internalInfo},
			wantInternal: internalInfo,
		},
		{
			name:         "No internal info",
			code:         "TEST_CODE_NO_INTERNAL",
			status:       http.StatusBadRequest,
			err:          originalErr,
			internalArgs: nil,
			wantInternal: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := New(tt.code, tt.status, tt.err, tt.internalArgs...)

			if got.Code != tt.code || got.HTTPStatus != tt.status || got.Err != tt.err || got.Internal != tt.wantInternal {
				t.Errorf("New() = %+v, want code=%v, status=%v, err=%v, internal=%v", 
                    got, tt.code, tt.status, tt.err, tt.wantInternal)
			}
		})
	}
}
