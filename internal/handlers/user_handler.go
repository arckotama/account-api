package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"account-api/internal/models"
	"gorm.io/gorm"
)

// UserHandler menampung dependency yang dibutuhkan handler, yaitu koneksi database.
type UserHandler struct {
	DB *gorm.DB
}

// NewUserHandler membuat UserHandler baru.
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// writeJSON menulis response JSON dengan status code tertentu
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError menulis pesan error dengan format JSON yang seragam.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

//// baru

// List menangani GET /users: mengambil semua user.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users := []models.User{}

	if err := h.DB.Order("id ASC").Find(&users).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "gagal mengambil data user")
		return
	}

	writeJSON(w, http.StatusOK, users)
}

//// baru
// Create menangani POST /users: membuat user baru.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) { // baru
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body JSON tidak valid")
		return
	}

	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name dan email wajib diisi")
		return
	}

	user := models.User{Name: req.Name, Email: req.Email}

	if err := h.DB.Create(&user).Error; err != nil {
		//println("ERROR CREATE USER:", err.Error())
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			writeError(w, http.StatusConflict, "email sudah dipakai")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal membuat user")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

//////// baru

// parseID mengambil {id} dari path dan mengubahnya menjadi angka.
func parseID(r *http.Request) (uint, error) { // baru
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GetByID menangani GET /users/{id}: mengambil satu user.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) { // baru
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka positif")
		return
	}

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "user tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mengambil user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

///////////// baru
// Update menangani PUT /users/{id}: mengganti data user.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) { // baru
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka positif")
		return
	}

	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "body JSON tidak valid")
		return
	}

	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "name dan email wajib diisi")
		return
	}

	var user models.User
	if err := h.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(w, http.StatusNotFound, "user tidak ditemukan")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mengambil user")
		return
	}

	user.Name = req.Name
	user.Email = req.Email

	if err := h.DB.Save(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			writeError(w, http.StatusConflict, "email sudah dipakai")
			return
		}
		writeError(w, http.StatusInternalServerError, "gagal mengubah user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

//////////// baru
// Delete menangani DELETE /users/{id}: menghapus user.
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) { // baru
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "id harus berupa angka positif")
		return
	}

	result := h.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		writeError(w, http.StatusInternalServerError, "gagal menghapus user")
		return
	}
	if result.RowsAffected == 0 {
		writeError(w, http.StatusNotFound, "user tidak ditemukan")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
