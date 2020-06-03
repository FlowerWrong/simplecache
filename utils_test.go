package simplecache

import "testing"

func TestToBytes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	actual, err := ToBytes("2KB")
	if err != nil {
		t.Error(err)
	}

	if actual != 2048 {
		t.Error(actual)
	}

	actual, err = ToBytes("2MB")
	if err != nil {
		t.Error(err)
	}

	if actual != 1024*1024*2 {
		t.Error(actual)
	}

	actual, err = ToBytes("2GB")
	if err != nil {
		t.Error(err)
	}

	if actual != 1024*1024*1024*2 {
		t.Error(actual)
	}
}

func TestGetRealSizeOf(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	actual := GetRealSizeOf("abcd")
	if actual != 8 {
		t.Error(actual)
	}
}
