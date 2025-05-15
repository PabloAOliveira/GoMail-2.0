package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"

	"email-service/database"
	"email-service/models"
)

type EmailMessage struct {
	ID      int    `json:"id"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func Publish(task models.EmailTask) (int, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Erro ao carregar o .env, usando variáveis de ambiente")
	}

	emailID, err := database.SaveEmail(task)
	if err != nil {
		log.Printf("Erro ao salvar email no banco: %s", err)
		return 0, err
	}

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Printf("Erro ao conectar ao RabbitMQ: %s", err)

		database.UpdateEmailStatus(emailID, models.StatusFailed)
		return emailID, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Erro ao criar canal: %s", err)
		database.UpdateEmailStatus(emailID, models.StatusFailed)
		return emailID, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue", // Nome da fila
		true,          // Durável
		false,         // Exclusiva
		false,         // Auto-delete
		false,         // No-wait
		nil,           // Propriedades adicionais
	)
	if err != nil {
		log.Printf("Erro ao declarar fila: %s", err)
		database.UpdateEmailStatus(emailID, models.StatusFailed)
		return emailID, err
	}

	message := EmailMessage{
		ID:      emailID,
		To:      task.To,
		Subject: task.Subject,
		Body:    task.Body,
	}

	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Erro ao serializar mensagem: %s", err)
		database.UpdateEmailStatus(emailID, models.StatusFailed)
		return emailID, err
	}

	err = ch.Publish(
		"",           // Exchange
		q.Name,       // Nome da fila
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Erro ao publicar na fila: %s", err)
		database.UpdateEmailStatus(emailID, models.StatusFailed)
		return emailID, err
	}
	
	log.Printf("Tarefa ID %d enfileirada com sucesso!aaaaa", emailID)
	return emailID, nil
}