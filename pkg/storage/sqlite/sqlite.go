package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/abtahihasan/students-api/pkg/config"
	"github.com/abtahihasan/students-api/pkg/types"
	_ "github.com/mattn/go-sqlite3"
)

type SQlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*SQlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			age INTEGER NOT NULL
		)
	`);

	if err != nil {
		return nil, err
	}

	return &SQlite{
		Db: db,
	}, nil
}


func (s *SQlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES(?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId();

	if err != nil {
		return 0, err
	}

	return lastId, nil
}



func (s *SQlite) GetStudentById(id int64) (types.Student, error) {

	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id = ?")

	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.ID, &student.Name, &student.Email, &student.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with id %d not found", id)
		}
		return types.Student{}, err
	}


	return student, nil

}


func (s *SQlite) GetStudentList() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students");

	if err != nil {
		return nil, err
	}
	
	defer stmt.Close();

	rows, err := stmt.Query();

	if err != nil {
		return nil, err
	}

	defer rows.Close();

	var students []types.Student;

	for rows.Next() {
		var student types.Student;

		 err := rows.Scan(&student.ID, &student.Name, &student.Email, &student.Age);
		if err != nil {
			return nil, err
		}
		slog.Info("student", slog.Any("student", student))
		students = append(students, student);

	}


	return students, nil

}