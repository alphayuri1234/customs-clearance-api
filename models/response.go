package models

import "strings"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func sanitizeError(err interface{}) interface{} {
	if err == nil {
		return nil
	}

	var errStr string
	switch e := err.(type) {
	case error:
		errStr = e.Error()
	case string:
		errStr = e
	default:
		return err
	}

	lowerErr := strings.ToLower(errStr)

	// Saring dan ubah error terkait database/driver agar tidak tampil sebagai exception langsung
	if strings.Contains(lowerErr, "duplicate key") || strings.Contains(lowerErr, "violates unique constraint") || strings.Contains(lowerErr, "unique") {
		return "data yang dimasukkan sudah terdaftar dalam sistem (duplikasi)"
	}
	if strings.Contains(lowerErr, "violates foreign key constraint") || strings.Contains(lowerErr, "foreign key") {
		return "referensi data tidak valid atau relasi data tidak ditemukan (violates foreign key)"
	}
	if strings.Contains(lowerErr, "dial tcp") || strings.Contains(lowerErr, "connection refused") || strings.Contains(lowerErr, "gorm: database") || strings.Contains(lowerErr, "sql: database") {
		return "gagal terhubung ke database internal"
	}
	if strings.Contains(lowerErr, "value too long for type") {
		return "nilai input melebihi batas panjang karakter database"
	}
	if strings.Contains(lowerErr, "violates not-null constraint") || strings.Contains(lowerErr, "null value in column") {
		return "kolom wajib tidak boleh kosong pada database"
	}
	if strings.Contains(lowerErr, "invalid input syntax") {
		return "format input tidak valid untuk tipe data database"
	}
	if strings.Contains(lowerErr, "pq:") || strings.Contains(lowerErr, "gorm:") || strings.Contains(lowerErr, "sql:") || strings.Contains(lowerErr, "postgres") {
		return "terjadi kesalahan internal pada sistem database"
	}

	return errStr
}

func ErrorResponse(message string, errors interface{}) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Errors:  sanitizeError(errors),
	}
}
