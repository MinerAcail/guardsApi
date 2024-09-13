package resolvers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

// parseID retrieves and converts the ID parameter from the URL.
func parseID(r *http.Request) (int, error) {
	idParam := chi.URLParam(r, "id")
	return strconv.Atoi(idParam)
}

// fetchStudentByID fetches a student by their ID from the database.
func fetchStudentByID(db *gorm.DB, id int) (*models.Student, error) {
	var student models.Student
	result := db.First(&student, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

// CreateStudent handles the creation of a new student
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

// GetAllStudents handles the retrieval of all students
func GetAllStudents(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	result := db.Find(&students)
	if result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, students)
}

// GetStudentByID retrieves a student by their ID
func GetStudentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	student, err := fetchStudentByID(db, id)
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

// UpdateStudentByID handles updating a student by their ID
func UpdateStudentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	student, err := fetchStudentByID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Student not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(student); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Save(student); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, student)
}

// DeleteStudentByID handles the deletion of a student by their ID
func DeleteStudentByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid student ID")
		return
	}

	student, err := fetchStudentByID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Student not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if result := db.Delete(student); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
