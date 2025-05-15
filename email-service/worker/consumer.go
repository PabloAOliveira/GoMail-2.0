package worker

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"

	"email-service/database"
	"email-service/email"
	"email-service/models"
	"email-service/rabbitmq"
)

func StartConsumer() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal: %s", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email_queue", 
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao declarar fila: %s", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("Erro ao configurar QoS: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // fila
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Erro ao consumir fila: %s", err)
	}

	log.Println("Consumidor iniciado. Esperando mensagens...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Mensagem recebida: %s", d.Body)
		
			var message rabbitmq.EmailMessage
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Erro ao decodificar mensagem: %s", err)
				d.Nack(false, false) 
				continue
			}

			err = database.UpdateEmailStatus(message.ID, models.StatusSending)
			if err != nil {
				log.Printf("Erro ao atualizar status para 'sending': %s", err)
				
			}
			
			task := models.EmailTask{
				ID:      message.ID,
				To:      message.To,
				Subject: message.Subject,
				Body:    message.Body,
			}
			
			log.Printf("Enviando email ID %d para: %s, assunto: %s", task.ID, task.To, task.Subject)
			err = email.SendEmail(task)
			
			if err != nil {
				log.Printf("Erro ao enviar e-mail ID %d: %s", task.ID, err)
				database.UpdateEmailStatus(task.ID, models.StatusFailed)
			} else {
				log.Printf("E-mail ID %d enviado com sucesso para %s", task.ID, task.To)
				database.UpdateEmailStatus(task.ID, models.StatusSent)
			}
			
			d.Ack(false) 
		}
	}()

	<-forever 
}