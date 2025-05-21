📬 GoMail — API de Processamento Assíncrono com Mensageria e Interface Web
GoMail é uma aplicação completa para envio assíncrono de e-mails, desenvolvida em Go, utilizando RabbitMQ como sistema de mensageria e PostgreSQL como banco de dados para persistência. Todo o sistema é containerizado com Docker, e inclui tanto a API backend quanto uma interface frontend feita com Vue 3 e Vuetify 3.

🚀 Funcionalidades

✅ Enfileiramento de tarefas de envio de e-mails via RabbitMQ

⚙️ Processamento assíncrono através de um worker que escuta a fila

📬 Envio de e-mails reais via SMTP

🧾 Persistência no PostgreSQL dos e-mails enviados e seus respectivos status (sucesso, erro, pendente etc.)

🌐 Interface web com Vue 3 + Vuetify 3 para visualização e controle dos e-mails

🔎 Rotas disponíveis na API:

- POST /send – Enfileira novo e-mail

- GET /get-id/:id – Retorna o status de um e-mail específico

- GET /get-all – Lista todos os e-mails registrados

📦 Tecnologias Utilizadas
- Go – linguagem principal da API e do worker

- Gin Gonic – framework web para criação da API

- RabbitMQ – mensageria para enfileiramento assíncrono

- PostgreSQL – banco de dados para persistência dos e-mails

- Docker & Docker Compose – conteinerização de todos os serviços

- AMQP – biblioteca github.com/streadway/amqp para RabbitMQ

- SMTP – envio real de e-mails

- dotenv – carregamento de variáveis de ambiente

- Makefile – automação de comandos de build e execução

- Vue 3 + Vuetify 3 – interface web moderna e responsiva
