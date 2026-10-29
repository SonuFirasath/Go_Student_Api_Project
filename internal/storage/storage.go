package storage
import "github.com/SonuFirasath/Go_Student_Api_Project/internal/types"
type Storage interface {
	CreateStudent(name string,email string, age int) (int64,error)
	GetStudentByID(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
}