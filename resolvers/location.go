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

// CreateHomeArrival handles the creation of a new home arrival record.
func CreateHomeArrival(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var homeArrival models.HomeArrival
	if err := json.NewDecoder(r.Body).Decode(&homeArrival); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get the current date
	currentDate := time.Now().Format("2006-01-02") // Format to YYYY-MM-DD

	// Check if a HomeArrival already exists for the same ParentID, StudentID, and current date
	var existingArrival models.HomeArrival
	if err := db.Where("parent_id = ? AND student_id = ? AND DATE(created_at) = ?", homeArrival.ParentID, homeArrival.StudentID, currentDate).First(&existingArrival).Error; err == nil {
		// Record exists, update the existing one
		existingArrival.Confirmed = homeArrival.Confirmed
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
	homeArrival.ID = uuid.New() // Ensure that a new UUID is generated
	if result := db.Create(&homeArrival); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusCreated, homeArrival)
}

func GetConfirmedArrivalsByParent(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Parse UUID from the URL path parameters
	parentID, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid parent UUID")
		return
	}

	// Get the current date
	currentDate := time.Now().Format("2006-01-02") // Format to YYYY-MM-DD

	// Query the database to retrieve all confirmed HomeArrival records for the given parent and today's date
	var confirmedArrivals []models.HomeArrival
	if err := db.Where("parent_id = ?  AND DATE(created_at) = ?", parentID.String(), currentDate).Find(&confirmedArrivals).Error; err != nil {
		handleError(w, http.StatusInternalServerError, "Error retrieving confirmed arrivals: "+err.Error())
		return
	}



	// Respond with the list of confirmed arrivals
	respondJSON(w, http.StatusOK, confirmedArrivals)
}
func GetAllConfirmedArrivals(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Get the current date
	currentDate := time.Now().Format("2006-01-02") // Format to YYYY-MM-DD

	// Query the database to retrieve all confirmed HomeArrival records for today's date
	var confirmedArrivals []models.HomeArrival
	if err := db.Where("DATE(created_at) = ?", currentDate).Find(&confirmedArrivals).Error; err != nil {
		handleError(w, http.StatusInternalServerError, "Error retrieving confirmed arrivals: "+err.Error())
		return
	}


	// Respond with the list of confirmed arrivals
	respondJSON(w, http.StatusOK, confirmedArrivals)
}

func GetAllConfirmedArrivalsStaff(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Get the current date
	currentDate := time.Now().Format("2006-01-02") // Format to YYYY-MM-DD

	// Query the database to retrieve all confirmed HomeArrival records for today's date
	var confirmedArrivals []models.SchoolArrival
	if err := db.Where("DATE(created_at) = ?", currentDate).Find(&confirmedArrivals).Error; err != nil {
		handleError(w, http.StatusInternalServerError, "Error retrieving confirmed arrivals: "+err.Error())
		return
	}

	

	// Respond with the list of confirmed arrivals
	respondJSON(w, http.StatusOK, confirmedArrivals)
}
// GetAllHomeArrivalsForThatWeek retrieves all home arrivals for the current week.
func GetAllHomeArrivalsForThatWeek(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Get the current time and calculate the start of the week (e.g., Monday)
	now := time.Now()
	startOfWeek := now.Truncate(24 * time.Hour).AddDate(0, 0, -int(now.Weekday()-1)) // Monday

	// Fetch all home arrivals for the current week
	var homeArrivals []models.HomeArrival
	if err := db.Where("created_at >= ?", startOfWeek).Find(&homeArrivals).Error; err != nil {
		handleError(w, http.StatusInternalServerError, "Error fetching home arrivals: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, homeArrivals)
}




// GetAllHomeArrivalsForThatWeekByParentId retrieves all home arrivals for a specified parent ID for the current week.
func GetAllHomeArrivalsForThatWeekByParentId(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Parse UUID from the request
	parentID, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid parent UUID")
		return
	}

	// Get the start and end of the current week
	startOfWeek := time.Now().Truncate(24 * time.Hour).AddDate(0, 0, -int(time.Now().Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	// Fetch home arrivals for the specified parent within the current week
	var homeArrivals []models.HomeArrival
	if err := db.Where("parent_id = ? AND created_at >= ? AND created_at < ?", parentID.String(), startOfWeek, endOfWeek).Find(&homeArrivals).Error; err != nil {
		handleError(w, http.StatusInternalServerError, "Error fetching home arrivals: "+err.Error())
		return
	}

	// Respond with the list of home arrivals
	respondJSON(w, http.StatusOK, homeArrivals)
}
