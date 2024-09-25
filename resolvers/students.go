package resolvers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mineracail/guardApi/models"
	"gorm.io/gorm"
)

// respondJSON sends a JSON response with the appropriate status code.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// handleError sends an error response with a custom message and status code.
func handleError(w http.ResponseWriter, status int, errMessage string) {
	http.Error(w, errMessage, status)
}

// parseUUID retrieves and converts the UUID parameter from the URL.
func parseUUID(r *http.Request) (uuid.UUID, error) {
	idParam := chi.URLParam(r, "id")
	return uuid.Parse(idParam)
}

// fetchStudentByUUID fetches a student by their UUID from the database.
func FetchStudentByUUID(db *gorm.DB, id uuid.UUID) (*models.Student, error) {
	var student models.Student
	result := db.Where("id = ?", id).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

// CreateStudent handles the creation of a new student.
func CreateStudent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Create(&student); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusCreated, student)
}

// GetAllStudents handles the retrieval of all students.
func GetAllStudents(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	result := db.Find(&students)
	if result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, students)
}

// GetStudentByID retrieves a student by their UUID.
func GetStudentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid student UUID")
		return
	}

	student, err := FetchStudentByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Student not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, student)
}

// UpdateStudentByID handles updating a student by their UUID.
func UpdateStudentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid student UUID")
		return
	}

	student, err := FetchStudentByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Student not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Save(&student); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, student)
}

// DeleteStudentByID handles the deletion of a student by their UUID.
func DeleteStudentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid student UUID")
		return
	}

	student, err := FetchStudentByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Student not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if result := db.Delete(&student); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}