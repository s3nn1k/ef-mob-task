package service

import "testing"

func TestFilterVerses(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		limit   int
		offset  int
		wantRes []string
	}{
		{
			name:    "first",
			text:    "bla-bla-bla",
			limit:   1,
			offset:  0,
			wantRes: []string{"bla-bla-bla"},
		},
		{
			name:    "second",
			text:    "text1\n\ntext2\n\ntext3\n\ntext4\n\ntext5",
			limit:   2,
			offset:  0,
			wantRes: []string{"text1", "text2"},
		},
		{
			name:    "third",
			text:    "text1\n\ntext2\n\ntext3\n\ntext4\n\ntext5",
			limit:   2,
			offset:  3,
			wantRes: []string{"text4", "text5"},
		},
		{
			name:    "fourth",
			text:    "text1\n\ntext2\n\ntext3\n\ntext4\n\ntext5",
			limit:   1,
			offset:  6,
			wantRes: []string{},
		},
		{
			name:    "fifth",
			text:    "text1\n\ntext2\n\ntext3\n\ntext4\n\ntext5",
			limit:   6,
			offset:  0,
			wantRes: []string{"text1", "text2", "text3", "text4", "text5"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			verses := filterVerses(test.text, test.limit, test.offset)

			if len(verses) != len(test.wantRes) {
				t.Fatalf("error: verses slice and wantRes must have the same length")
			}

			for i := 0; i < len(verses); i++ {
				if verses[i] != test.wantRes[i] {
					t.Fatalf("error: want %s string with %v index, but got %s", test.wantRes[i], i, verses[i])
				}
			}
		})
	}
}
