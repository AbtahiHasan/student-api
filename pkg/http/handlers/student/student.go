package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/abtahihasan/students-api/pkg/types"
	"github.com/abtahihasan/students-api/pkg/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, false, http.StatusBadRequest, nil)
			return
		}


		// validate request 
		if	err := validator.New().Struct(student); err != nil {
			
			response.WriteJSON(w, false, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
		return
		}
		w.Write([]byte("hello world"))
	}
}