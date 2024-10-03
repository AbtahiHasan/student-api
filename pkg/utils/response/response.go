package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abtahihasan/students-api/pkg/types"
	"github.com/go-playground/validator/v10"
)



func WriteJSON(w http.ResponseWriter, success bool, status int, data interface{}) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := types.Response{
		Success: success,
		Message: status,
		Data:    data,
	}
	return json.NewEncoder(w).Encode(response)
}


func ValidationError (errs validator.ValidationErrors) string {
	var errMessages []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, err.Field()+" is required")
		default:
			errMessages = append(errMessages, err.Field()+" is invalid")
		}
	}


	message := strings.Join(errMessages, ", ")
	return message
}