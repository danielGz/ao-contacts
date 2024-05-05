package api

import (
	"accelone-contacts/model"
	"accelone-contacts/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func RegisterCreateContactsAPI(r *mux.Router) {
	svc := service.NewInMemoryContactService()
	// register all endpoints
	r.HandleFunc("/contacts", getContacts(svc)).Methods("GET")
	r.HandleFunc("/contacts/{id}", getContact(svc)).Methods("GET")
	r.HandleFunc("/contacts", createContact(svc)).Methods("POST")
	r.HandleFunc("/contacts/{id}", updateContact(svc)).Methods("PUT")
	r.HandleFunc("/contacts/{id}", deleteContact(svc)).Methods("DELETE")
}

// Handler for creating a contact
func createContact(service service.ContactService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c model.Contact
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			log.Error().Str("endpoint", "POST /contacts").Msg("Failed to decode request body")
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		result, err := service.Create(c)
		if err != nil {
			log.Error().Str("endpoint", "POST /contacts").Msg("Failed to create contact")
			jsonResponseError(w, "Contact creation failed: "+err.Error(), http.StatusConflict)
			return
		}

		log.Info().Str("endpoint", "POST /contacts").Msg("Contact created successfully")
		json.NewEncoder(w).Encode(result)
	}
}

// Handler for retrieving a contact by ID
func getContact(service service.ContactService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		result, err := service.GetById(id)
		if err != nil {
			log.Error().Str("endpoint", "GET /contacts/{id}").Str("contact_id", id).Msg("Failed to get contact")
			jsonResponseError(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Info().Str("endpoint", "GET /contacts/{id}").Str("contact_id", id).Msg("Contact retrieved successfully")
		json.NewEncoder(w).Encode(result)
	}
}

// Handler for updating a contact
func updateContact(service service.ContactService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var c model.Contact
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			log.Error().Str("endpoint", "PUT /contacts/{id}").Str("contact_id", id).Msg("Failed to decode request body")
			jsonResponseError(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		c.Id = id // Ensure the ID is set to the URL parameter
		result, err := service.Update(c)
		if err != nil {
			log.Error().Str("endpoint", "PUT /contacts/{id}").Str("contact_id", id).Msg("Failed to update contact")
			jsonResponseError(w, "Contact update failed: "+err.Error(), http.StatusNotFound)
			return
		}

		log.Info().Str("endpoint", "PUT /contacts/{id}").Str("contact_id", id).Msg("Contact updated successfully")
		json.NewEncoder(w).Encode(result)
	}
}

// Handler for deleting a contact
func deleteContact(service service.ContactService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		_, err := service.Delete(id)
		if err != nil {
			log.Error().Str("endpoint", "DELETE /contacts/{id}").Str("contact_id", id).Msg("Failed to delete contact")
			jsonResponseError(w, "Contact deletion failed: "+err.Error(), http.StatusNotFound)
			return
		}

		log.Info().Str("endpoint", "DELETE /contacts/{id}").Str("contact_id", id).Msg("Contact deleted successfully")
		w.WriteHeader(http.StatusNoContent) // HTTP 204 if deletion was successful
	}
}

func getContacts(service service.ContactService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Default values for pagination
		page := 1
		limit := 10

		// Parse page from query parameters
		if p, ok := r.URL.Query()["page"]; ok && len(p[0]) > 0 {
			if parsedPage, err := strconv.Atoi(p[0]); err == nil {
				page = parsedPage
			}
		}
		// Parse limit from query parameters
		if l, ok := r.URL.Query()["limit"]; ok && len(l[0]) > 0 {
			if parsedLimit, err := strconv.Atoi(l[0]); err == nil {
				limit = parsedLimit
			}
		}

		contacts, err := service.Get(page, limit)
		if err != nil {
			jsonResponseError(w, "Failed to retrieve contacts", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(contacts)
	}
}
