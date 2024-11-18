package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rest-api/repository"
	"rest-api/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client

func init() {
	// Carrega variáveis de ambiente
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env load error", err)
	}

	log.Println("env file loaded")

	// Cria o cliente do MongoDB
	MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal("mongo connection error", err)
	}

	err = MongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("ping failed", err)
	}

	log.Println("mongo connected")
}

func main() {
	// Fechar a conexão com o MongoDB quando a aplicação encerrar
	defer MongoClient.Disconnect(context.Background())

	// Cria uma referência à coleção do MongoDB
	coll := MongoClient.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("COLLECTION_NAME"))

	// Cria o repositório com a coleção
	employeeRepo := repository.NewEmployeeRepository(coll)

	// Cria o serviço, passando o repositório como dependência
	employeeService := service.NewEmployeeService(employeeRepo)

	// Configura o roteador
	r := mux.NewRouter()

	// Define as rotas e associa os manipuladores de requisições
	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)
	r.HandleFunc("/employee", employeeService.CreateEmployee).Methods(http.MethodPost)
	r.HandleFunc("/employee/{id}", employeeService.GetEmployeeById).Methods(http.MethodGet)
	r.HandleFunc("/employee", employeeService.GetAllEmployee).Methods(http.MethodGet)
	r.HandleFunc("/employee/{id}", employeeService.UpdateEmployeeById).Methods(http.MethodPut)
	r.HandleFunc("/employee/{id}", employeeService.DeleteEmployeeById).Methods(http.MethodDelete)
	r.HandleFunc("/employee", employeeService.DeleteAllEmployee).Methods(http.MethodDelete)

	// Inicia o servidor
	log.Println("server is running on :4444")
	http.ListenAndServe(":4444", r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("running..."))
}
