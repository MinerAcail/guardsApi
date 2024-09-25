package resolvers

import (
	"encoding/json"
	"net/http"

	
	"github.com/google/uuid"
	"github.com/mineracail/guardApi/models"
	"gorm.io/gorm"
)



// fetchStaffByUUID fetches a staff by their UUID from the database.
func FetchStaffByUUID(db *gorm.DB, id uuid.UUID) (*models.Staff, error) {
	var staff models.Staff
	result := db.Where("id = ?", id).First(&staff)
	if result.Error != nil {
		return nil, result.Error
	}
	return &staff, nil
}

// CreateStaff handles the creation of a new staff.
func CreateStaff(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var staff models.Staff
	if err := json.NewDecoder(r.Body).Decode(&staff); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Create(&staff); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusCreated, staff)
}

// GetAllStaffs handles the retrieval of all staffs.
func GetAllStaffs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var staffs []models.Staff
	result := db.Find(&staffs)
	if result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, staffs)
}

// GetStaffByID retrieves a staff by their UUID.
func GetStaffByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid staff UUID")
		return
	}

	staff, err := FetchStaffByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Staff not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, staff)
}

// UpdateStaffByID handles updating a staff by their UUID.
func UpdateStaffByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid staff UUID")
		return
	}

	staff, err := FetchStaffByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Staff not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&staff); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Save(&staff); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, staff)
}

// DeleteStaffByID handles the deletion of a staff by their UUID.
func DeleteStaffByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid staff UUID")
		return
	}

	staff, err := FetchStaffByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Staff not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if result := db.Delete(&staff); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}