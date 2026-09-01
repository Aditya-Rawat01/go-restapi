package storage

import "github.com/Aditya-Rawat01/go-restapi/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	UpdateStudent(data types.UpdateStudent, id int64) (types.Student, error)
	DeleteStudentById(id int64) (types.Student, error)
}
