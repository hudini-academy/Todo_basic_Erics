package mysql

import (
	"database/sql"
	"todo/pkg/models"
)

type TodoModel struct {
	DB *sql.DB
}

// This will insert a new task into the database.
func (m *TodoModel) Insert(title string) (int, error) {

	// Insert the task into the database.
	smt := `INSERT INTO todos(name, created, expires) VALUES (?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Execute thr query with the given parameters.
	_, err := m.DB.Exec(smt, title, 7)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// This will return a specific task based on its id.
func (m *TodoModel) GetAll() ([]*models.Todo, error) {
	// Fetch a single task from database.
	smt := `SELECT * FROM todos ORDER BY id DESC`
	rows, errFetching := m.DB.Query(smt)
	if errFetching!= nil {
        return nil, errFetching
    }

	// Initialize a pointer to a new zeroed Snippet struct.
	s := []*models.Todo{}

	// iterate through the rows and append to s slice.
	for rows.Next() {
		singleRow := &models.Todo{}
		err := rows.Scan(&singleRow.Id, &singleRow.Title, &singleRow.Created, &singleRow.Expires)

		if err!= nil {
            return nil, err
        }

        s = append(s, singleRow)
	}
	return s, nil
}

// Delete the task from database.
func (m *TodoModel) DelTaskDB(id int) (bool, error) {
	smt := `DELETE FROM todos WHERE Id = ?`
	_, err := m.DB.Exec(smt, id)

	if err!= nil {
        return false, err
    }
	return true, nil
}

// Update the task in database.

func (m *TodoModel) UpdateTaskDB(id int, title string) (bool, error) {
	smt := `UPDATE todos SET name = ? WHERE Id = ?`
	_, err := m.DB.Exec(smt, title, id)
	if err!= nil {
        return false, err
    }
	return true, nil
}
