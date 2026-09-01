package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Aditya-Rawat01/go-restapi/internal/config/utils"
	"github.com/Aditya-Rawat01/go-restapi/internal/storage"
	"github.com/Aditya-Rawat01/go-restapi/internal/types"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student...")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			utils.WriteJSON(w, http.StatusBadRequest, utils.GeneralError(fmt.Errorf("Request Body cannot be empty!")))
			return
		}
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, utils.GeneralError(fmt.Errorf("Bad Request or Invalid JSON")))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validationErrs := err.(validator.ValidationErrors)
			utils.WriteJSON(w, http.StatusBadRequest, utils.ValidationError(validationErrs))
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		student.Id = lastId
		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, map[string]any{
			"status":  "ok",
			"msg":     "Student created successfully",
			"student": student,
		})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(("id"))
		slog.Info("Getting info of the student with id: ", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, utils.GeneralError(err))
			return
		}
		student, err := storage.GetStudentByID(intId)

		if err != nil {
			slog.Error("error getting user", slog.String("id",id))
			utils.WriteJSON(w, http.StatusInternalServerError, utils.GeneralError(err))
			return
		}

		utils.WriteJSON(w, http.StatusOK, student)
	}
}
