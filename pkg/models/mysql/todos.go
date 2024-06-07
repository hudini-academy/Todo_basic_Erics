package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"todo/pkg/models"
)

type TodoModel struct {
	DB *sql.DB
}

// This will insert a new task into the database.
func (m *TodoModel) Insert(title string) (int, error) {

	// Convert the title to lower case for checking whether the title has "special:" string prefix.
	lowerTitle := strings.ToLower(strings.TrimSpace(title))
	// if has the special string prefix then insert to special table and todos table.
	if strings.HasPrefix(lowerTitle, "special:") {
		smt := `INSERT INTO special() VALUES ()`

		// Execute thr query with the given parameters.
		row, err := m.DB.Exec(smt)
		if err != nil {
			return 0, err
		}
		// Get the id of the special task and insert all the details to todos table.
		specialId, errInsertSpecial := row.LastInsertId()
		if errInsertSpecial != nil {
			log.Println(errInsertSpecial)
			return 0, errInsertSpecial
		}

		smt2 := `INSERT INTO todos(name, created, expires, spl_id) VALUES (?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`
		_, errInsertTodos := m.DB.Exec(smt2, title, 7, specialId)
		if errInsertTodos != nil {
			return 0, errInsertTodos
		}

		return 0, nil
	}
	// Insert the task into the database.
	smt := `INSERT INTO todos(name, created, expires, spl_id) VALUES (?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?)`

	// Execute thr query with the given parameters.
	_, err := m.DB.Exec(smt, title, 7, 0)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// Add tags.
func (m *TodoModel) InsertTag(tag string, id int) (int, error) {
	var tags string
	// Insert the task into the database.
	smt1 := `SELECT tags FROM todos WHERE id = ?`
	row := m.DB.QueryRow(smt1, id)
	row.Scan(&tags)

	if tags != "" {
		tags = fmt.Sprintf("%s, %s", tags, tag)
	} else {
		tags = fmt.Sprintf(",%s, %s", tags, tag)
	}
	smt := `UPDATE todos SET tags = ? WHERE id = ?`

	// Execute thr query with the given parameters.
	_, err := m.DB.Exec(smt, tags, id)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

// This will return a specific task based on its id.
func (m *TodoModel) GetAll() ([]*models.Todo, error) {
	// Fetch a single task from database.
	smt := `SELECT id, name, created, expires, tags FROM todos ORDER BY id DESC`
	rows, errFetching := m.DB.Query(smt)
	if errFetching != nil {
		return nil, errFetching
	}

	// Initialize a pointer to a new zeroed Snippet struct.
	s := []*models.Todo{}

	// iterate through the rows and append to s slice.
	for rows.Next() {		
		singleRow := &models.Todo{}
		err := rows.Scan(&singleRow.Id, &singleRow.Title, &singleRow.Created, &singleRow.Expires, &singleRow.Tags)
		tags := strings.Split(singleRow.Tags, ",")
		singleRow.Tag = tags

		log.Println(tags)
		if err != nil {
			return nil, err
		}

		s = append(s, singleRow)
	}
	return s, nil
}

// Get special tasks.
func (m *TodoModel) GetSpecial() ([]*models.Todo, error) {
	// Fetch a single task from database.
	smt := `SELECT t.id, t.name, t.created, t.expires, t.tags FROM todos t, special s WHERE t.spl_id=s.spl_id ORDER BY id DESC`
	rows, errFetching := m.DB.Query(smt)
	if errFetching != nil {
		return nil, errFetching
	}

	// Initialize a pointer to a new zeroed Snippet struct.
	s := []*models.Todo{}

	// iterate through the rows and append to s slice.
	for rows.Next() {
		singleRow := &models.Todo{}
		err := rows.Scan(&singleRow.Id, &singleRow.Title, &singleRow.Created, &singleRow.Expires, &singleRow.Tags)

		if err != nil {
			return nil, err
		}

		s = append(s, singleRow)
	}
	return s, nil
}

// Get special tasks.
func (m *TodoModel) GetTagTasks(tag string) ([]*models.Todo, error) {
	// Fetch a single task from database.
	smt := `SELECT id, name, created, expires, tags FROM todos WHERE tags LIKE '%`+tag+`%' ORDER BY id DESC`;
	rows, errFetching := m.DB.Query(smt)
	if errFetching != nil {
		return nil, errFetching
	}

	// Initialize a pointer to a new zeroed Snippet struct.
	s := []*models.Todo{}

	// iterate through the rows and append to s slice.
	for rows.Next() {
		singleRow := &models.Todo{}
		err := rows.Scan(&singleRow.Id, &singleRow.Title, &singleRow.Created, &singleRow.Expires, &singleRow.Tags)

		if err != nil {
			return nil, err
		}

		s = append(s, singleRow)

	}
	return s, nil
}

// Delete the task from database.
func (m *TodoModel) DelTaskDB(id int) (bool, error) {
	var splId int64
	// Check whether the task is a special task.
	stm1 := `SELECT spl_id FROM todos WHERE id = ?`
	row := m.DB.QueryRow(stm1, id)
	row.Scan(&splId)

	// if the task is a special task then delete.
	if splId != 0 {
		smt := `DELETE FROM special WHERE spl_id = ?`
		_, err := m.DB.Exec(smt, splId)

		if err != nil {
			return false, err
		}
	}
	// Delete the task from todos table.
	smt := `DELETE FROM todos WHERE Id = ?`
	_, err := m.DB.Exec(smt, id)

	if err != nil {
		return false, err
	}
	return true, nil
}

// Update the task in database.
func (m *TodoModel) UpdateTaskDB(id int, title string) (bool, error) {
	smt := `UPDATE todos SET name = ? WHERE Id = ?`
	_, err := m.DB.Exec(smt, title, id)
	if err != nil {
		return false, err
	}
	return true, nil
}
