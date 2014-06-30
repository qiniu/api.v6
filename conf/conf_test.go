package conf

import (
	"strings"
	"testing"
)

func TestUA(t *testing.T) {
	err := SetUser("")
	if err != nil {
		t.Fatal(err)
	}
	err = SetUser("错误的UA")
	if err == nil {
		t.Fatal("expect an invalid ua format")
	}
}

func TestFormat(t *testing.T) {
	v := formatUserAgent("test")
	if !strings.Contains(v, "test") {
		t.Fatal("should include user")
	}
	if !strings.HasPrefix(v, "QiniuGo/"+version) {
		t.Fatal("invalid format")
	}
}
