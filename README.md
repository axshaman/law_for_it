# GDPR & Cookie Compliance Automation

## Overview

This project automates the generation of GDPR and Cookie compliance documents for websites, surveys, and other interactive activities. Users interact through a **Telegram bot**, **WhatsApp bot**, or **web interface**, sequentially answering questions and selecting checklist items to generate compliance documents in **PDF** or **DOCX** formats.

The backend is built in **Go (Golang)**, the frontend uses **HTML, JavaScript, and CSS3**, and all services run in **Docker** with environment configuration stored in `.env` and `config.yaml`.

This system can also be used for rapid development of any document generation system.

---

## Features

- **Automated GDPR & Cookie policy generation**
- **Interactive Telegram bot, WhatsApp bot, and Web interface for input collection**
- **Document export in PDF & DOCX formats**
- **Simple web interface for customization**
- **Modular and scalable architecture**
- **Dockerized deployment**
- **Flexible document templating system**
- **Payment system for premium document downloads**

---

## Tech Stack

### Backend:
- **Go (Golang)** – Core application logic
- **Fiber** – Web framework for API endpoints
- **Gorm** – ORM for database interactions
- **MySQL** – Database for storing user sessions and document history
- **go-telegram-bot-api** – Telegram bot integration
- **go-whatsapp** – WhatsApp bot integration
- **go-pdf** & **unioffice** – PDF and DOCX generation

### Frontend:
- **HTML, JavaScript, CSS3** – Lightweight UI for customization
- **Alpine.js** – Minimal JavaScript framework for interactivity

### Infrastructure:
- **Docker & Docker Compose** – Containerized deployment
- **.env & config.yaml** – Environment-based configuration management
- **NGINX (optional)** – Reverse proxy support

---

## Installation & Setup

### 1. Clone the Repository
```sh
git clone https://github.com/axshaman/law_for_it.git
cd law_for_it
```

### 2. Configure Environment Variables
Create a `.env` file in the root directory:
```ini
TELEGRAM_BOT_TOKEN=your-bot-token
WHATSAPP_API_KEY=your-whatsapp-api-key
DB_TYPE=mysql  # or "sqlite"
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=gdpr_db
DOCUMENT_FORMATS=pdf,docx  # Allowed export formats
```

### 3. Define Additional Configurations
Edit `config.yaml`:
```yaml
server:
  port: 8080
  host: "0.0.0.0"

database:
  type: "mysql"
  host: "localhost"
  port: 3306
  user: "your_user"
  password: "your_password"
  name: "gdpr_db"

document:
  formats: ["pdf", "docx"]

bot:
  telegram_token: "your-telegram-bot-token"
  whatsapp_api_key: "your-whatsapp-api-key"
```

### 4. Configure Payment System
Edit `billing.cfg`:
```ini
PAYMENT_API_URL=https://payment-gateway.com/api
PAYMENT_API_KEY=your-api-key
PAYMENT_PROVIDER=stripe  # or "paypal"
```

### 5. Build & Run with Docker Compose
```sh
docker-compose up --build -d
```

---

## Document Templating System

### Concept
In the `autodocs` folder, document templates are stored in three parts, and **all three files must have the same filename**:
1. **TXT files** – Contain the document structure with placeholders (`XXXX`, `YYYY`, etc.).
2. **CFG files** – Define input fields for each placeholder (e.g., `XXXX=input1`, `YYYY=checkbox1`).
3. **HTML templates** – Contain the final document layout.

### Workflow
1. The system reads TXT, CFG, and HTML files from `autodocs/`, ensuring that all three files have the same name.
2. User inputs are collected via one of the following:
   - **Telegram bot**
   - **WhatsApp bot**
   - **Web interface**
3. The script combines the inputs with the template and generates PDF/DOCX files.
4. Free PDF downloads include a watermark.
5. Users can pay to download watermark-free PDFs or DOCX files.

---

## Payment System

- Free users can download **PDF files with a watermark**.
- **Paid users** can download:
  - Watermark-free PDF files.
  - Editable DOCX files.
- Payments are processed through **Stripe** or **PayPal**, configured in `billing.cfg`.
- Users can view their **account balance** and **transaction history** in the web interface.
- Balance top-ups are handled via an external backend API.

---

## Usage

### Telegram Bot Workflow
1. Start the bot in Telegram.
2. Answer a sequence of GDPR-related questions.
3. Select compliance checklist items.
4. Choose the output format (PDF or DOCX).
5. Receive a **free** PDF with a watermark or pay for a premium version.

### WhatsApp Bot Workflow
1. Start a conversation with the WhatsApp bot.
2. Follow the guided prompts to enter required document data.
3. Choose the output format (PDF or DOCX).
4. Receive the generated document via WhatsApp.
5. If a premium document is requested, complete the payment process.

### Web Interface
1. Open `index.html` in a browser.
2. Fill out the GDPR/Cookie compliance form.
3. Click "Generate" to receive a downloadable document.
4. If a premium document is requested, complete the payment process.

---

## API Endpoints
| Method | Endpoint         | Description                    |
|--------|-----------------|--------------------------------|
| GET    | `/status`       | Health check                   |
| POST   | `/generate`     | Generate a GDPR/Cookie policy  |
| GET    | `/downloadf/:id`| Download generated document    |
| GET    | `/downloadp/:id`| Download full document         |
| GET    | `/balance`      | Get user balance               |
| GET    | `/payhistory`   | Get billing's history          |
| POST   | `/topup`        | Top up user balance            |

---

## Docker Configuration

### Dockerfile
```dockerfile
FROM golang:1.19-alpine
WORKDIR /app
COPY . .
RUN go mod download && go build -o gdpr-server main.go
CMD ["./gdpr-server"]
EXPOSE 8080
```

### docker-compose.yml
```yaml
version: "3.8"
services:
  app:
    build: .
    ports:
      - "8080:8080"
    env_file:
      - .env
    volumes:
      - ./data:/app/data
    restart: always
  mysql:
    image: mysql:latest
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: gdpr_db
      MYSQL_USER: your_user
      MYSQL_PASSWORD: your_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
volumes:
  mysql_data:
```

---

## License
This project is licensed under the MIT License.
