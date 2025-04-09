package storage

import (
	"database/sql"
	"os"
	"time"

	"gourbot/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// Storage is responsible for managing the SQLite database.
type Storage struct {
	db       *sql.DB
	filename string
}

// NewStorage initializes a new Storage instance with the given filename.
func NewStorage(filename string) *Storage {
	return &Storage{
		filename: filename,
	}
}

// Open opens the SQLite database.
func (s *Storage) Open() error {
	db, err := sql.Open("sqlite3", s.filename)
	if err != nil {
		return err
	}
	s.db = db
	return s.createTables()
}

// Close closes the SQLite database connection.
func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Delete removes the SQLite database file.
func (s *Storage) Delete() error {
	if err := s.Close(); err != nil {
		return err
	}
	return os.Remove(s.filename)
}

// createTables creates the necessary tables in the SQLite database.
func (s *Storage) createTables() error {
	if s.db == nil {
		return sql.ErrConnDone
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS example (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS tgdump (
			uid INTEGER PRIMARY KEY AUTOINCREMENT,
			out BOOLEAN NOT NULL,
			data TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS tgusers (
			id INTEGER PRIMARY KEY,
			name TEXT DEFAULT '',
			created_at INTEGER NOT NULL,
			seen_at INTEGER NOT NULL,
			permissions TEXT DEFAULT '',
			info TEXT DEFAULT ''
		);`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

// AddTgRecord adds a new record to the tgdump table.
func (s *Storage) AddTgRecord(out bool, data string) error {
	if s.db == nil {
		return sql.ErrConnDone
	}

	query := `INSERT INTO tgdump (out, data) VALUES (?, ?);`
	_, err := s.db.Exec(query, out, data)
	return err
}

// Note: The database stores timestamps as Unix time (integer), but the TgUser struct uses time.Time.
// Ensure proper conversion between Unix time and time.Time during read and write operations.

// AddTgUser adds a new user to the tgusers table.
func (s *Storage) AddTgUser(user *models.TgUser) error {
	query := `INSERT INTO tgusers (id, name, created_at, seen_at, permissions, info) VALUES (?, ?, ?, ?, ?, ?)`
	permissions := user.PermissionsToString()
	createdAtUnix := user.CreatedAt.Unix()
	seenAtUnix := user.SeenAt.Unix()
	_, err := s.db.Exec(query, user.Id, user.Name, createdAtUnix, seenAtUnix, permissions, user.Info)
	return err
}

// TgUserExists checks if a user exists in the tgusers table by ID.
func (s *Storage) TgUserExists(id int64) (bool, error) {
	query := `SELECT COUNT(*) FROM tgusers WHERE id = ?`
	var count int
	err := s.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetTgUser retrieves a user by ID from the tgusers table.
func (s *Storage) GetTgUser(id int64) (*models.TgUser, error) {
	query := `SELECT name, created_at, seen_at, permissions, info FROM tgusers WHERE id = ?`
	row := s.db.QueryRow(query, id)

	var name, info, permissions string
	var createdAtUnix, seenAtUnix int64
	err := row.Scan(&name, &createdAtUnix, &seenAtUnix, &permissions, &info)
	if err != nil {
		return nil, err
	}

	user := models.NewTgUser(id, name, info)
	user.CreatedAt = time.Unix(createdAtUnix, 0)
	user.SeenAt = time.Unix(seenAtUnix, 0)
	user.AddPermissionsFromString(permissions)

	return user, nil
}

// GetAllTgUsers retrieves all users from the tgusers table.
func (s *Storage) GetAllTgUsers() ([]*models.TgUser, error) {
	query := `SELECT id, name, created_at, seen_at, permissions, info FROM tgusers`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.TgUser
	for rows.Next() {
		var id int64
		var name, info, permissions string
		var createdAtUnix, seenAtUnix int64
		err := rows.Scan(&id, &name, &createdAtUnix, &seenAtUnix, &permissions, &info)
		if err != nil {
			return nil, err
		}

		user := models.NewTgUser(id, name, info)
		user.CreatedAt = time.Unix(createdAtUnix, 0)
		user.SeenAt = time.Unix(seenAtUnix, 0)
		user.AddPermissionsFromString(permissions)

		users = append(users, user)
	}

	return users, nil
}

// UpdateTgUser updates an existing user in the tgusers table.
func (s *Storage) UpdateTgUser(user *models.TgUser) error {
	query := `UPDATE tgusers SET name = ?, seen_at = ?, permissions = ?, info = ? WHERE id = ?`
	permissions := user.PermissionsToString()
	seenAtUnix := user.SeenAt.Unix()
	_, err := s.db.Exec(query, user.Name, seenAtUnix, permissions, user.Info, user.Id)
	return err
}
