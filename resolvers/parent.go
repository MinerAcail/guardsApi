package resolvers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mineracail/guardApi/models"
	"gorm.io/gorm"
)

// FetchParentByUUID fetches a parent by their UUID from the database.
func FetchParentByUUID(db *gorm.DB, id uuid.UUID) (*models.Parent, error) {
	var parent models.Parent
	result := db.Where("id = ?", id).First(&parent)
	if result.Error != nil {
		return nil, result.Error
	}
	return &parent, nil
}

// CreateParent handles the creation of a new parent.
func CreateParent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var parent models.Parent
	if err := json.NewDecoder(r.Body).Decode(&parent); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Create(&parent); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusCreated, parent)
}

// GetAllParents handles the retrieval of all parents.
func GetAllParents(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var parents []models.Parent
	result := db.Find(&parents)
	if result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, parents)
}

// GetParentByID retrieves a parent by their UUID.
func GetParentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid parent UUID")
		return
	}

	parent, err := FetchParentByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Parent not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, parent)
}

func GetChildByParentID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Parse UUID from the request
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid parent UUID")
		return
	}

	// Fetch the parent by UUID
	var parent models.Parent
	if err := db.First(&parent, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Parent not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// Fetch all students supervised by this parent
	var students []models.Student
	if err := db.Where("id = ANY(?)", parent.Supervise).Find(&students).Error; err != nil {
		handleError(w, http.StatusInternalServerError, "Error fetching students: "+err.Error())
		return
	}

	// Create a response with parent and supervised students
	response := map[string]interface{}{
		"parent":   parent,
		"students": students,
	}

	// Respond with the parent and supervised students
	respondJSON(w, http.StatusOK, response)
}
// UpdateParentByID handles updating a parent by their UUID.
func UpdateParentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid parent UUID")
		return
	}

	parent, err := FetchParentByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Parent not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&parent); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Save(&parent); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, parent)
}

// DeleteParentByID handles the deletion of a parent by their UUID.
func DeleteParentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid parent UUID")
		return
	}

	parent, err := FetchParentByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Parent not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if result := db.Delete(&parent); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
