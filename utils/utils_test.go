package utils

import (
	"bytes"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old // restoring the real stdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestPrint2DArrayAsTable(t *testing.T) {
	tests := []struct {
		name    string
		mxc     int
		headers []string
		data    [][]string
		want    string
	}{
		{
			name:    "With Headers",
			mxc:     10,
			headers: []string{"ID", "Name"},
			data:    [][]string{{"1", "Alice"}, {"2", "Bob"}},
			want:    "|====|====|\n| ID | Name |\n|====|====|\n| 1  | Alice |\n| 2  | Bob |\n",
		},
		{
			name: "Without Headers",
			mxc:  10,
			data: [][]string{{"ID", "Name"}, {"1", "Alice"}, {"2", "Bob"}},
			want: "|====|====|\n| ID | Name |\n|====|====|\n| 1  | Alice |\n| 2  | Bob |\n",
		},
		{
			name:    "With Long Strings",
			mxc:     5,
			headers: []string{"ID", "Description"},
			data:    [][]string{{"1", "A very long description that needs clipping"}},
			want:    "|====|====|\n| ID | Descr |\n|====|====|\n| 1  | A ver |\n",
		},
		{
			name:    "Empty Data",
			mxc:     10,
			headers: []string{"ID", "Name"},
			data:    [][]string{},
			want:    "|====|====|\n| ID | Name |\n|====|====|\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureOutput(func() {
				Print2DArrayAsTable(tt.mxc, tt.headers, tt.data)
			})

			if got != tt.want {
				t.Errorf("Print2DArrayAsTable() got = %v, want %v", got, tt.want)
			}
		})
	}
}
