package repository

import (
	"context"
	"fmt"
	"rest-api/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmployeeRepository struct {
	MongoCollection *mongo.Collection
}

// Função para criar um novo repositório de funcionários
func NewEmployeeRepository(collection *mongo.Collection) *EmployeeRepository {
	return &EmployeeRepository{MongoCollection: collection}
}

// Método exportado: salva um novo funcionário no banco de dados
func (r *EmployeeRepository) Save(emp *model.Employee) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), emp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Outro exemplo de método exportado, caso queira adicionar mais operações:
func (r *EmployeeRepository) FindById(id string) (*model.Employee, error) {
	var emp model.Employee
	err := r.MongoCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&emp)
	if err != nil {
		return nil, err
	}
	return &emp, nil
}
func (r *EmployeeRepository) FindAll() ([]model.Employee, error) {
	results, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var emps []model.Employee
	err = results.All(context.Background(), &emps)

	if err != nil {
		return nil, fmt.Errorf("results decode error")
	}

	return emps, nil
}

func (r *EmployeeRepository) UpdateById(id string, entity *model.Employee) (int64, error) {
	result, err := r.MongoCollection.UpdateOne(context.Background(), bson.D{{Key: "id", Value: id}}, bson.D{{Key: "$set", Value: entity}})

	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (r *EmployeeRepository) DeleteById(id string) (int64, error) {
	result, err := r.MongoCollection.DeleteOne(context.Background(), bson.D{{Key: "id", Value: id}})

	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *EmployeeRepository) DeleteAll() (int64, error) {
	result, err := r.MongoCollection.DeleteMany(context.Background(), bson.D{}) // Alterado de DeleteOne para DeleteMany
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
