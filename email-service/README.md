**GoMail — API de Processamento Assíncrono com Mensageria**

GoMail é uma API feita em Go para envio de e-mails de forma assíncrona, utilizando RabbitMQ como sistema de mensageria. O projeto está totalmente containerizado com Docker e possui um consumidor que escuta a fila e dispara os e-mails conforme as mensagens são processadas.

🚀 Funcionalidades

- Enfileiramento de tarefas de envio de e-mail via RabbitMQ

- Processamento assíncrono através de um worker consumidor

- Envio de e-mails reais via SMTP

- Organização em múltiplos serviços com Docker

- Estrutura limpa, escalável e simples de manter

🧰 Tecnologias Utilizadas
Go – linguagem principal da API e do worker

- Gin Gonic – framework web leve para a API

- RabbitMQ – mensageria para enfileiramento assíncrono

- Docker & Docker Compose – conteinerização dos serviços

- AMQP – biblioteca github.com/streadway/amqp para conexão com o RabbitMQ

- SMTP – envio real de e-mails

dotenv – carregamento de variáveis de ambiente (github.com/joho/godotenv)

Makefile – comandos simplificados para build e execução
