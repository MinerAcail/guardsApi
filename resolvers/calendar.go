package resolvers

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/mineracail/guardApi/models"
	"gorm.io/gorm"
)


// fetchCalendarByUUID fetches a Calendar by their UUID from the database.
func FetchCalendarByUUID(db *gorm.DB, id uuid.UUID) (*models.Calendar, error) {
	var Calendar models.Calendar
	result := db.Where("id = ?", id).First(&Calendar)
	if result.Error != nil {
		return nil, result.Error
	}
	return &Calendar, nil
}

// CreateCalendar handles the creation of a new Calendar.
func CreateCalendar(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var Calendar models.Calendar
	if err := json.NewDecoder(r.Body).Decode(&Calendar); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Create(&Calendar); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusCreated, Calendar)
}

// GetAllCalendars handles the retrieval of all Calendars.
func GetAllCalendars(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var Calendars []models.Calendar
	result := db.Find(&Calendars)
	if result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, Calendars)
}

// GetCalendarByID retrieves a Calendar by their UUID.
func GetCalendarByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Calendar UUID")
		return
	}

	Calendar, err := FetchCalendarByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Calendar not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, Calendar)
}

// UpdateCalendarByID handles updating a Calendar by their UUID.
func UpdateCalendarByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Calendar UUID")
		return
	}

	Calendar, err := FetchCalendarByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Calendar not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&Calendar); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Save(&Calendar); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, Calendar)
}

// DeleteCalendarByID handles the deletion of a Calendar by their UUID.
func DeleteCalendarByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Calendar UUID")
		return
	}

	Calendar, err := FetchCalendarByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Calendar not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if result := db.Delete(&Calendar); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}