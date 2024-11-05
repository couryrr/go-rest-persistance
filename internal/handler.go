package internal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

type UserHandler struct{
	repo Repo
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) 

type Repo interface {
    Migrate() error
    Create(user *User) (int, error)
    Read(id int) (*User, error)
    Update(article *User) error
    //Delete(id int) error
}

type SQLiteRepo struct {
    db *sql.DB
}

func (r *SQLiteRepo) Migrate() error {
    query := `
    CREATE TABLE IF NOT EXISTS users(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        date_of_birth DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        date_created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        date_updated DATETIME
    );
    `
    _, err := r.db.Exec(query)
    return err
}

func (r *SQLiteRepo) Create(user *User) (int, error){ 
    query := `
        INSERT INTO users(name, email, date_of_birth) VALUES(?, ?, ?);
    `

    res, err := r.db.Exec(query, user.Name, user.Email, user.DateOfBirth)
    if err != nil {
        return -1, err
    }

    id, err := res.LastInsertId()
    if err != nil {
        return -1, err
    }

    return int(id), nil

}

func (r *SQLiteRepo) Read(id int) (*User, error){ 

	query := `
		SELECT id, name, email, date_of_birth FROM USERS WHERE id=?
	`

	row := r.db.QueryRow(query, id)

	var user User

	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.DateOfBirth); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLiteRepo) Update(user *User) (error){ 
	slog.Info("update called")
	query := "UPDATE users SET name = ?, email = ?, date_of_birth = ? WHERE id = ?"
	res, err := r.db.Exec(query, user.Name, user.Email, user.DateOfBirth, user.Id)
	if err != nil {
		return err
	}	

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count != 1 {
		// should probably rollback here...
	}
	return nil
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepo{
    return &SQLiteRepo{
        db: db,
    }
}

func NewUserHandler(repo Repo) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

func (h *UserHandler) GetHandler()(http.Handler){
	mux := http.NewServeMux()
	mux.HandleFunc(h.HandleAddUser())
	mux.HandleFunc(h.HandleGetUserById())
	mux.HandleFunc(h.HandleUpdateUser())
	return mux
}

func (h *UserHandler) HandleAddUser()(string, handlerFunc){
	return "POST /" , func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			slog.Error("unable to decode user", "err", err)
		}
		id, err := h.repo.Create(&user)
		if err != nil {
			slog.Error("unable to persist user", "err", err)
		}
		slog.Info("create user", "user", id)
	}
}

func (h *UserHandler) HandleGetUserById()(string, handlerFunc){
	return "GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		pv_id := r.PathValue("id")
		ctx := r.Context()
		if pv_id == ""{
			http.Error(w, "no id sent", http.StatusBadRequest)
		}

		id, err := strconv.Atoi(pv_id)

		if err != nil {
			http.Error(w, "bad id value sent", http.StatusBadRequest)
		}
		
		go func(){
			time.Sleep(3*time.Second)
			user, err := h.repo.Read(id)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows){
					slog.Error("error finding user", "err", err)
				}
				http.NotFound(w, r)
				return
			}
			json.NewEncoder(w).Encode(user)
		}()
		select{
		case <-ctx.Done():
			slog.Error("Timeout occred")
		}
	}	
}

func (h *UserHandler) HandleUpdateUser()(string, handlerFunc){
	return "PATCH /update", func(w http.ResponseWriter, r *http.Request) {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			slog.Error("unable to decode user", "err", err)
			http.Error(w, "malformed user object", http.StatusBadRequest)
			return
		}

		err = h.repo.Update(&user)	

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)		
			return
		}
	}
}

