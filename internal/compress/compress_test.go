package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func TestCompress(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "обычная строка",
			data:    []byte("gfdfgdfgdfgdfdgdfg"),
			wantErr: false,
		},
		{
			name:    "пустые данные",
			data:    []byte(""),
			wantErr: false,
		},
		{
			name:    "JSON",
			data:    []byte(`{"id":"Alloc","type":"gauge","value":123.45}`),
			wantErr: false,
		},
		{
			name:    "большие данные",
			data:    bytes.Repeat([]byte("abc"), 1000),
			wantErr: false,
		},
		{
			name:    "nil",
			data:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() ошибка = %v, ожидалась = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// распаковываем и сравниваем с исходными данными
			zr, err := gzip.NewReader(bytes.NewReader(got))
			if err != nil {
				t.Fatalf("Ошибка создания gzip.Reader: %v", err)
			}
			defer zr.Close()

			decompressed, err := io.ReadAll(zr)
			if err != nil {
				t.Fatalf("Ошибка распаковки: %v", err)
			}

			if !bytes.Equal(decompressed, tt.data) {
				t.Errorf("После распаковки = %q, ожидалось %q", decompressed, tt.data)
			}
		})
	}
}
