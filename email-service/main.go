package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"email-service/database"
	"email-service/models"
	"email-service/rabbitmq"
	"email-service/worker"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Erro ao carregar .env, usando variáveis de ambiente")
	}

	database.InitDB()

	go worker.StartConsumer()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.POST("/send-email", func(c *gin.Context) {
		var task models.EmailTask

		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		emailID, err := rabbitmq.Publish(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao publicar na fila"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{
			"message": "Tarefa enfileirada com sucesso!a",
			"id": emailID,
		})
	})

	router.DELETE("/email/:id", func(c *gin.Context) {
		id := c.Param("id")
		var emailID int
		_, err := fmt.Sscanf(id, "%d", &emailID)
		if err != nil || emailID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}
		
		err = database.DeleteEmail(emailID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar email"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Email deletado com sucesso!"})
	})
		

	router.GET("/emails", func(c *gin.Context) {
		emails, err := database.GetAllEmails()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar emails"})
			return
		}
		c.JSON(http.StatusOK, emails)
	})

	router.GET("/email/:id", func(c *gin.Context) {
		id := c.Param("id")
		var emailID int
		_, err := fmt.Sscanf(id, "%d", &emailID)
		if err != nil || emailID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
			return
		}
		
		email, err := database.GetEmailByID(emailID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email não encontrado"})
			return
		}
		
		c.JSON(http.StatusOK, email)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "2000"
	}
	router.Run(":" + port)
}