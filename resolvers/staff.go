package resolvers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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



func CreateSchoolArrival(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var SchooArrival models.SchoolArrival
	if err := json.NewDecoder(r.Body).Decode(&SchooArrival); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get the current date
	currentDate := time.Now().Format("2006-01-02") // Format to YYYY-MM-DD

	// Check if a HomeArrival already exists for the same ParentID, StudentID, and current date
	var existingArrival models.SchoolArrival
	if err := db.Where("staff_id = ? AND student_id = ? AND DATE(created_at) = ?", SchooArrival.StaffID, SchooArrival.StudentID, currentDate).First(&existingArrival).Error; err == nil {
		// Record exists, update the existing one
		existingArrival.Confirmed = SchooArrival.Confirmed
		if result := db.Save(&existingArrival); result.Error != nil {
			handleError(w, http.StatusInternalServerError, result.Error.Error())
			return
		}
		respondJSON(w, http.StatusOK, existingArrival)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Other error occurred during the query
		handleError(w, http.StatusInternalServerError, "Error checking for existing arrival")
		return
	}

	// No record exists for today, create a new one
	SchooArrival.ID = uuid.New() // Ensure that a new UUID is generated
	if result := db.Create(&SchooArrival); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusCreated, SchooArrival)
}


func GetConfirmedArrivalsByStaff(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Parse UUID from the URL path parameters
	parentID, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid parent UUID")
		return
	}

	// Get the current date
	currentDate := time.Now().Format("2006-01-02") // Format to YYYY-MM-DD

	// Query the database to retrieve all confirmed HomeArrival records for the given parent and today's date
	var confirmedArrivals []models.SchoolArrival
	if err := db.Where("staff_id = ?  AND DATE(created_at) = ?", parentID.String(), currentDate).Find(&confirmedArrivals).Error; err != nil {
		handleError(w, http.StatusInternalServerError, "Error retrieving confirmed arrivals: "+err.Error())
		return
	}


	// Respond with the list of confirmed arrivals
	respondJSON(w, http.StatusOK, confirmedArrivals)
}