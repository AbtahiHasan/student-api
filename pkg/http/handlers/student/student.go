package student

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/abtahihasan/students-api/pkg/storage"
	"github.com/abtahihasan/students-api/pkg/types"
	"github.com/abtahihasan/students-api/pkg/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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

		lastId, err := storage.CreateStudent(student.Name, student.Email, int(student.Age))

		if err != nil {
			response.WriteJSON(w, false, http.StatusInternalServerError, nil)
			return
		}
		

		response.WriteJSON(w, true, http.StatusCreated, map[string]int64{
			"id": lastId,})
	}
}


func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		intId,err:= strconv.ParseInt(id, 10, 64)

		if err != nil {
			response.WriteJSON(w, false, http.StatusBadRequest, nil)
			return
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			response.WriteJSON(w, false, http.StatusInternalServerError, nil)
			return
		}

		response.WriteJSON(w, true, http.StatusOK, student)
	}
}