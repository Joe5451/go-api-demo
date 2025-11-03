package domain

import "testing"

func TestBook_Validate(t *testing.T) {
	tests := []struct {
		name    string
		b       *Book
		wantErr error
	}{
		{
			name:    "valid book",
			b:       &Book{Title: "Test Book", Author: "Test Author"},
			wantErr: nil,
		},
		{
			name:    "empty title",
			b:       &Book{Title: "", Author: "Test Author"},
			wantErr: ErrTitleRequired,
		},
		{
			name:    "empty author",
			b:       &Book{Title: "Test Book", Author: ""},
			wantErr: ErrAuthorRequired,
		},
		{
			name:    "empty title and author",
			b:       &Book{Title: "", Author: ""},
			wantErr: ErrTitleRequired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.b.Validate()
			if err != tt.wantErr {
				t.Errorf("Book.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
