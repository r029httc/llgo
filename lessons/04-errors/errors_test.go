package errs

import (
	"errors"
	"strconv"
	"testing"
)

func TestParseAge(t *testing.T) {
	if n, err := ParseAge("30"); err != nil || n != 30 {
		t.Fatalf("ParseAge(\"30\") = %d, %v", n, err)
	}

	_, err := ParseAge("-1")
	if !errors.Is(err, ErrNegative) {
		t.Errorf("預期 ErrNegative，得到 %v", err)
	}

	_, err = ParseAge("abc")
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) { // errors.As 可以取出被包裝的特定錯誤型別
		t.Errorf("預期 *strconv.NumError，得到 %T", err)
	}
}

func TestValidateName(t *testing.T) {
	var ve *ValidationError
	if err := ValidateName(""); !errors.As(err, &ve) || ve.Field != "name" {
		t.Errorf("ValidateName(\"\") = %v", err)
	}
	if err := ValidateName("Gopher"); err != nil {
		t.Errorf("ValidateName(\"Gopher\") = %v, want nil", err)
	}
}
