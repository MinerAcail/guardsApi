package resolvers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mineracail/guardApi/models"
	"gorm.io/gorm"
)

// FetchMessageByUUID fetches a Message by their UUID from the database.
func FetchMessageByUUID(db *gorm.DB, id uuid.UUID) (*models.Message, error) {
	if _, err := uuid.Parse(id.String()); err != nil {
		return nil, err
	}
	var message models.Message
	result := db.Where("id = ?", id).First(&message)
	if result.Error != nil {
		return nil, result.Error
	}
	return &message, nil
}

// CreateMessage handles the creation of a new Message.
func CreateMessag(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var message models.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Create(&message); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusCreated, message)
}

func CreateMessage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	log.Println("Starting CreateMessage handler")
	w.Header().Set("Content-Type", "application/json")

	// Use a separate struct for the request with string receiver_id
	var messageRequest struct {
		Content  string    `json:"content"`
		SenderID uuid.UUID `json:"sender_id"`

		ReceiverID string `json:"receiver_id"` // Keep as string for JSON parsing
		Status     string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&messageRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request body",
		})
		return
	}

	// Validate required fields
	if messageRequest.Content == "" || messageRequest.ReceiverID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Content and receiver_id are required",
		})
		return
	}

	// Parse the receiver_id string into UUID
	receiverID, err := uuid.Parse(messageRequest.ReceiverID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid receiver_id format - must be a valid UUID",
		})
		return
	}

	// Create the message with the parsed UUID
	newMessage := &models.Message{
		Content:    messageRequest.Content,
		ReceiverID: receiverID,
		SenderID:   messageRequest.SenderID,
		Status:     messageRequest.Status,
		CreatedAt:  time.Now(),
	}

	result := db.Create(newMessage)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to create message",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMessage)
}

// CreateMessageToMultiple handles creating messages for multiple recipients
func CreateMessageToMultiple(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var messageRequest struct {
		Content    string      `json:"content"`
		SenderID   uuid.UUID   `json:"sender_id"`
		Recipients []uuid.UUID `json:"recipients"`
	}

	if err := json.NewDecoder(r.Body).Decode(&messageRequest); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var messages []models.Message
	for _, recipientID := range messageRequest.Recipients {
		message := models.Message{
			Content:    messageRequest.Content,
			SenderID:   messageRequest.SenderID,
			ReceiverID: recipientID,
		}
		messages = append(messages, message)
	}

	if result := db.Create(&messages); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusCreated, messages)
}

// GetAllMessages handles the retrieval of all Messages.
func GetAllMessages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var messages []models.Message
	result := db.Find(&messages)
	if result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, messages)
}

// GetMessageByID retrieves a Message by their UUID.
func GetMessageByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Message UUID")
		return
	}

	if _, err := uuid.Parse(id.String()); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Message UUID")
		return
	}

	message, err := FetchMessageByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Message not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondJSON(w, http.StatusOK, message)
}

// UpdateMessageByID handles updating a Message by their UUID.
func UpdateMessageByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Message UUID")
		return
	}

	message, err := FetchMessageByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Message not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if result := db.Save(&message); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondJSON(w, http.StatusOK, message)
}

// DeleteMessageByID handles the deletion of a Message by their UUID.
func DeleteMessageByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid Message UUID")
		return
	}

	message, err := FetchMessageByUUID(db, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, http.StatusNotFound, "Message not found")
		} else {
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if result := db.Delete(&message); result.Error != nil {
		handleError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
