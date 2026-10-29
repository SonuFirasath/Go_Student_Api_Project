package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/SonuFirasath/Go_Student_Api_Project/internal/types"
	"github.com/SonuFirasath/Go_Student_Api_Project/internal/utils/response"
	"github.com/go-playground/validator/v10"
	"github.com/SonuFirasath/Go_Student_Api_Project/internal/storage"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		
		slog.Info("Creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		
		if errors.Is(err,io.EOF) {
			response.WriteJson(w,http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		// request validation

		if err:=validator.New().Struct(student); err != nil {

			validateErrs := err.(validator.ValidationErrors)

			response.WriteJson(w,http.StatusBadRequest,response.ValidateError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("Student created successfully",slog.String("userId",fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err) )
			return
		}

		response.WriteJson(w,http.StatusCreated, map[string] int64 {"id":lastId})
	}

}

func GetByID(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request){

		id := r.PathValue("id")

		slog.Info("getting a student", slog.String("id",id))
		
		intId, err := strconv.ParseInt(id,10,64)
		
		if err != nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		
		student, err := storage.GetStudentByID(intId)
		
		if err != nil{
			slog.Error("error getting user", slog.String("id",id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}

}