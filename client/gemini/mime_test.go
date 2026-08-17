package gemini

import "testing"

func TestMimeType(t *testing.T) {
	heic := []byte("\x00\x00\x00\x18ftypheic")

	cases := []struct {
		name           string
		stated, served string
		data           []byte
		want           string
	}{
		{"呼叫端說了就聽它的", "image/heic", "application/pdf", heic, "image/heic"},
		{"沒說就用伺服器給的", "", "image/heic", heic, "image/heic"},
		{"伺服器的值帶 charset 要剝掉", "", "text/plain; charset=utf-8", []byte("hi"), "text/plain"},
		{"伺服器給 octet-stream 等於沒說", "", "application/octet-stream", []byte("%PDF-1.4"), "application/pdf"},
		{"都沒有才嗅探", "", "", []byte("%PDF-1.4"), "application/pdf"},
	}
	for _, c := range cases {
		if got := mimeType(c.stated, c.served, c.data); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
