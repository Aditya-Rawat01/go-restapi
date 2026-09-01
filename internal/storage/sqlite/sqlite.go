package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Aditya-Rawat01/go-restapi/internal/config"
	"github.com/Aditya-Rawat01/go-restapi/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		age INTEGER
	)`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentByID(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)

		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (s *Sqlite) UpdateStudent(data types.UpdateStudent, id int64) (types.Student, error) {
	query := "UPDATE students SET "
	updates := []string{}
	args := []any{}
	if data.Name != nil {
		updates = append(updates, "name=? ")
		args = append(args, *data.Name)
	}
	if data.Age != nil {
		updates = append(updates, "age=? ")
		args = append(args, *data.Age)
	}
	if data.Email != nil {
		updates = append(updates, "email=? ")
		args = append(args, *data.Email)
	}

	partial := strings.Join(updates, ", ")
	query += partial
	query += "where id = ? RETURNING id, name, email, age"
	args = append(args, id)

	var student types.Student
	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(args...).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		fmt.Print("here????")
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

}

func (s *Sqlite) DeleteStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("DELETE FROM STUDENTS WHERE id = ? RETURNING id, name, email, age")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}
