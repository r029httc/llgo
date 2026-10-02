package errs

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

func TestParseAge(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	for _, s := range []string{"0", "30", "150"} {
		if _, err := ParseAge(s); err != nil {
			t.Errorf("ParseAge(%q) 不應出錯: %v", s, err)
		}
	}
	tests := []struct {
		in   string
		want error
	}{
		{"-1", ErrNegative},
		{"151", ErrTooOld},
		{"999", ErrTooOld},
	}
	for _, tt := range tests {
		if _, err := ParseAge(tt.in); !errors.Is(err, tt.want) {
			t.Errorf("ParseAge(%q) err = %v, want errors.Is(%v)", tt.in, err, tt.want)
		}
	}
	_, err := ParseAge("abc")
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Errorf("ParseAge(\"abc\") 應包裝 *strconv.NumError（用 %%w），得到 %v", err)
	}
}

func writeFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "app.conf") // TempDir 會在測試結束時自動刪除
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadConfig(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	path := writeFile(t, `# 這是註解
host = localhost
port=8080

  name =  我的服務  
url = http://x.com/?a=b
`)
	got, err := ReadConfig(path)
	if err != nil {
		t.Fatalf("ReadConfig: %v", err)
	}
	want := map[string]string{"host": "localhost", "port": "8080", "name": "我的服務", "url": "http://x.com/?a=b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadConfig = %v, want %v", got, want)
	}

	_, err = ReadConfig(filepath.Join(t.TempDir(), "missing.conf"))
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("檔案不存在時應 errors.Is(err, os.ErrNotExist)，得到 %v", err)
	}

	for content, wantLine := range map[string]int{
		"a=1\nthis line is bad\n": 2,
		"a=1\n\n# c\n= no key\n":  4,
	} {
		_, err = ReadConfig(writeFile(t, content))
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("格式錯誤應回傳 *ParseError，得到 %v", err)
		} else if pe.Line != wantLine {
			t.Errorf("ParseError.Line = %d, want %d", pe.Line, wantLine)
		}
	}
}

func TestValidateUser(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	if err := ValidateUser(User{"小明", "ming@example.com", 20}); err != nil {
		t.Errorf("合法使用者不應出錯: %v", err)
	}
	err := ValidateUser(User{"", "bad-email", -3})
	for _, want := range []error{ErrEmptyName, ErrInvalidEmail, ErrInvalidAge} {
		if !errors.Is(err, want) {
			t.Errorf("應包含 %v，得到 %v", want, err)
		}
	}
	err = ValidateUser(User{"小華", "hua@example.com", 200})
	if !errors.Is(err, ErrInvalidAge) || errors.Is(err, ErrEmptyName) {
		t.Errorf("只應有 ErrInvalidAge，得到 %v", err)
	}
}

func TestRetry(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	errTemp := errors.New("暫時失敗")

	calls := 0
	err := Retry(5, func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Errorf("第 3 次成功：err=%v calls=%d, want nil 3", err, calls)
	}

	calls = 0
	err = Retry(4, func() error { calls++; return errTemp })
	if !errors.Is(err, errTemp) || calls != 4 {
		t.Errorf("一直失敗：err=%v calls=%d, want 包裝 errTemp 且呼叫 4 次", err, calls)
	}

	calls = 0
	err = Retry(10, func() error {
		calls++
		return errors.Join(errors.New("db down"), ErrPermanent)
	})
	if !errors.Is(err, ErrPermanent) || calls != 1 {
		t.Errorf("永久錯誤：err=%v calls=%d, want 立即停止（1 次）", err, calls)
	}

	calls = 0
	_ = Retry(0, func() error { calls++; return nil })
	if calls != 1 {
		t.Errorf("attempts=0 時至少要呼叫 1 次，實際 %d 次", calls)
	}
}

func TestSafely(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	if err := Safely(func() {}); err != nil {
		t.Errorf("沒有 panic 時應回傳 nil，得到 %v", err)
	}
	if err := Safely(func() { panic("糟糕") }); err == nil {
		t.Error("panic(string) 應轉成 error")
	}
	errBoom := errors.New("boom")
	if err := Safely(func() { panic(errBoom) }); !errors.Is(err, errBoom) {
		t.Errorf("panic(error) 應包裝原錯誤，得到 %v", err)
	}
	err := Safely(func() {
		var m map[string]int
		m["x"] = 1 // 寫入 nil map 會 panic
	})
	if err == nil {
		t.Error("寫入 nil map 的 panic 應轉成 error")
	}
}
