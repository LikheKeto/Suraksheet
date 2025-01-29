# Suraksheet - A Safe Place for Your Documents

Suraksheet is a **document storage, management, and search system** designed to store, process, and search your documents efficiently. 

Currently, only image documents are supported but the vision is to support PDFs as well. It has **image compression** feature for minimizing size of documents. The system utilizes **Python (OCR & compression), Golang (backend & search service), Elasticsearch (indexing & search), and SvelteKit (frontend)**, all containerized using **Docker**.

## **Architecture Overview**
Suraksheet consists of four main services:

1. **Suraksheet Extractor (Python)** - Extracts text from uploaded documents using OCR (Tesseract) and compresses images.
2. **Suraksheet Backend (Golang)** - Handles authentication, document metadata storage, search functionality, and document compression requests.
3. **Suraksheet Frontend (SvelteKit)** - A user-friendly web interface for managing and searching documents.
4. **Suraksheet Search Service (ElasticSearch)** - Interacts with Elasticsearch to provide full-text search capabilities.

All services are containerized using **Docker Compose**, with **RabbitMQ for task queueing**, **MinIO for object storage**, **PostgreSQL for metadata storage**, and **Elasticsearch for full-text search**.

---

## **Installation & Setup**
### **1. Prerequisites**
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### **2. Clone the Repository**
```sh
git clone https://github.com/LikheKeto/Suraksheet.git
cd suraksheet
```

### **3. Configure Environment Variables**
Copy the example environment file and set up required variables.
```sh
cp .env.example .env
```
Modify `.env` as needed for PostgreSQL, MinIO, RabbitMQ, and Elasticsearch settings.

### **4. Run the Services**
```sh
docker-compose up --build
```
This will start all services, including the frontend, backend, OCR extraction, and databases.

---

## **Usage**
### **1. Uploading Documents**
1. Navigate to `http://localhost:3000`.
2. Login or register an account.
3. Upload a document, you can edit/compress documents in document's page. [Editing is WIP]

### **2. Searching for Documents**
1. Use the search bar to find documents by keywords.
2. Elasticsearch enables fast full-text search across stored documents.

### **3. Managing Bins**
1. Organize documents into bins for easy categorization.
2. Manage bins/accounts for different people. [WIP]

---

## **Technology Stack**
| Component         | Technology              |
| ----------------- | ----------------------- |
| Frontend          | SvelteKit               |
| Backend API       | Golang                  |
| OCR & Compression | Python (Tesseract)      |
| Database          | PostgreSQL              |
| Search Engine     | Elasticsearch           |
| Object Storage    | MinIO                   |
| Message Queue     | RabbitMQ                |
| Containerization  | Docker & Docker Compose |

---

## **Contributing**
1. Fork the repository.
2. Create a feature branch (`git checkout -b feature-branch`).
3. Commit changes (`git commit -m "Add new feature"`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

---

## **License**
This project is licensed under the MIT License.

---

## **Contact**
For questions or support, reach out via email or open an issue on GitHub.
