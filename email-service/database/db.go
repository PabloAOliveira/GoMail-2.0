package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"email-service/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erro ao pingar o banco de dados: %v", err)
	}

	log.Println("Conectado ao banco de dados com sucesso!")
}

func SaveEmail(task models.EmailTask) (int, error) {
	var id int
	query := `INSERT INTO emails (to_email, subject, body, status) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id`
	
	err := DB.QueryRow(query, task.To, task.Subject, task.Body, "enqueued").Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("erro ao salvar email no banco: %v", err)
	}
	
	return id, nil
}

func UpdateEmailStatus(id int, status string) error {
	query := `UPDATE emails SET status = $1 WHERE id = $2`
	
	_, err := DB.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do email: %v", err)
	}
	
	return nil
}

func GetEmailByID(id int) (models.EmailTask, error) {
	var task models.EmailTask
	query := `SELECT id, to_email, subject, body, status, created_at, updated_at 
			  FROM emails WHERE id = $1`
	
	err := DB.QueryRow(query, id).Scan(
		&task.ID, &task.To, &task.Subject, &task.Body, 
		&task.Status, &task.CreatedAt, &task.UpdatedAt)
	
	if err != nil {
		return models.EmailTask{}, fmt.Errorf("erro ao buscar email: %v", err)
	}
	
	return task, nil
}

func DeleteEmail(id int) error {
    query := `DELETE FROM emails WHERE id = $1`
    _, err := DB.Exec(query, id)
    if err != nil {
        return fmt.Errorf("erro ao deletar email: %v", err)
    }
    return nil
}

func GetAllEmails() ([]models.EmailTask, error) {
	query := `SELECT id, to_email, subject, body, status, created_at, updated_at FROM emails`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar emails: %v", err)
	}
	defer rows.Close()

	var tasks []models.EmailTask
	for rows.Next() {
		var task models.EmailTask
		err := rows.Scan(&task.ID, &task.To, &task.Subject, &task.Body,
			&task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear email: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}