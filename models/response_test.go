package models

import (
	"errors"
	"testing"
)

func TestSanitizeError(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "Non-DB string error",
			input:    "request register tidak valid",
			expected: "request register tidak valid",
		},
		{
			name:     "Non-DB error object",
			input:    errors.New("some regular error"),
			expected: "some regular error",
		},
		{
			name:     "Postgres duplicate key string",
			input:    "pq: duplicate key value violates unique constraint \"users_email_key\"",
			expected: "data yang dimasukkan sudah terdaftar dalam sistem (duplikasi)",
		},
		{
			name:     "Postgres duplicate key error object",
			input:    errors.New("pq: duplicate key value violates unique constraint \"users_email_key\""),
			expected: "data yang dimasukkan sudah terdaftar dalam sistem (duplikasi)",
		},
		{
			name:     "Postgres foreign key string",
			input:    "violates foreign key constraint \"ports_country_id_fkey\"",
			expected: "referensi data tidak valid atau relasi data tidak ditemukan (violates foreign key)",
		},
		{
			name:     "DB connection refused error",
			input:    "dial tcp 127.0.0.1:5432: connect: connection refused",
			expected: "gagal terhubung ke database internal",
		},
		{
			name:     "Null value constraint error",
			input:    "null value in column \"name\" of relation \"users\" violates not-null constraint",
			expected: "kolom wajib tidak boleh kosong pada database",
		},
		{
			name:     "Value too long error",
			input:    "value too long for type character varying(10)",
			expected: "nilai input melebihi batas panjang karakter database",
		},
		{
			name:     "Generic GORM error",
			input:    "gorm: record not found",
			expected: "terjadi kesalahan internal pada sistem database",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeError(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeError() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrorResponse(t *testing.T) {
	resp := ErrorResponse("registrasi gagal", errors.New("pq: duplicate key value violates unique constraint"))
	if resp.Success {
		t.Error("Success should be false")
	}
	if resp.Message != "registrasi gagal" {
		t.Errorf("Message = %s, want %s", resp.Message, "registrasi gagal")
	}
	if resp.Errors != "data yang dimasukkan sudah terdaftar dalam sistem (duplikasi)" {
		t.Errorf("Errors = %v, want %s", resp.Errors, "data yang dimasukkan sudah terdaftar dalam sistem (duplikasi)")
	}
}
