package service

import (
	"encoding/json"
	"log"
	"net/http"
	"rest-api/model"
	"rest-api/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EmployeeService struct {
	EmployeeRepo *repository.EmployeeRepository // Repositório de funcionários
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// NewEmployeeService cria um novo EmployeeService com o repositório passado como dependência
func NewEmployeeService(repo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{EmployeeRepo: repo}
}

func (svc *EmployeeService) CreateEmployee(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(response).Encode(res)

	var emp model.Employee

	// Decodificando o corpo da requisição
	err := json.NewDecoder(request.Body).Decode(&emp)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid body", err)
		res.Error = err.Error()
		return
	}

	// Gerando um novo ID para o funcionário
	emp.ID = uuid.NewString()

	// Chama o método Save do repositório para salvar o novo funcionário
	createdEmployee, err := svc.EmployeeRepo.Save(&emp)
	if err != nil {
		// Caso ocorra erro ao salvar no banco
		response.WriteHeader(http.StatusInternalServerError)
		log.Println("Error saving employee:", err)
		res.Error = "Failed to create employee: " + err.Error()
		return
	}

	// Caso a criação seja bem-sucedida, retorna o funcionário criado
	res.Data = createdEmployee
	response.WriteHeader(http.StatusCreated)
}

func (svc *EmployeeService) GetEmployeeById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(response).Encode(res)

	empID := mux.Vars(request)["id"]
	log.Println("employee id", empID)

	emp, err := svc.EmployeeRepo.FindById(empID)
	if err != nil {
		// Caso ocorra erro ao buscar funcionário
		response.WriteHeader(http.StatusNotFound)
		log.Println("Error search employee by id:", err)
		res.Error = "Employee not found: " + err.Error()
		return
	}
	res.Data = emp
	response.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) GetAllEmployee(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(response).Encode(res)

	emps, err := svc.EmployeeRepo.FindAll()
	if err != nil {
		// Caso ocorra erro ao buscar todos os funcionários
		response.WriteHeader(http.StatusInternalServerError)
		log.Println("Error retrieving employees:", err)
		res.Error = "Error retrieving employees: " + err.Error()
		return
	}
	res.Data = emps
	response.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) UpdateEmployeeById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(response).Encode(res)

	empID := mux.Vars(request)["id"]
	log.Println("employee id", empID)

	if empID == "" {
		response.WriteHeader(http.StatusBadRequest)
		log.Println("Error: invalid employee ID")
		res.Error = "Invalid employee ID"
		return
	}

	var emp model.Employee

	// Decodificando o corpo da requisição
	err := json.NewDecoder(request.Body).Decode(&emp)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		log.Println("Invalid body", err)
		res.Error = err.Error()
		return
	}

	// Garantir que o ID do emp não seja alterado durante a atualização
	emp.ID = empID

	count, err := svc.EmployeeRepo.UpdateById(empID, &emp)
	if err != nil {
		// Caso ocorra erro ao atualizar o funcionário
		response.WriteHeader(http.StatusInternalServerError)
		log.Println("Error updating employee by id:", err)
		res.Error = "Error updating employee: " + err.Error()
		return
	}

	res.Data = count
	response.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) DeleteEmployeeById(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(response).Encode(res)

	empID := mux.Vars(request)["id"]
	log.Println("employee id", empID)

	count, err := svc.EmployeeRepo.DeleteById(empID)
	if err != nil {
		// Caso ocorra erro ao deletar o funcionário
		response.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting employee by id:", err)
		res.Error = "Error deleting employee: " + err.Error()
		return
	}

	res.Data = count
	response.WriteHeader(http.StatusOK)
}

func (svc *EmployeeService) DeleteAllEmployee(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("Content-Type", "application/json")

	res := &Response{}
	defer json.NewEncoder(response).Encode(res)

	count, err := svc.EmployeeRepo.DeleteAll()
	if err != nil {
		// Caso ocorra erro ao deletar todos os funcionários
		response.WriteHeader(http.StatusInternalServerError)
		log.Println("Error deleting all employees:", err)
		res.Error = "Error deleting all employees: " + err.Error()
		return
	}

	res.Data = count
	response.WriteHeader(http.StatusOK)
}
