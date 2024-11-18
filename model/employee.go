package model

type Employee struct {
	ID         string `json:"id,omitempty"`         // Alterado de id para ID
	Name       string `json:"name,omitempty"`       // Alterado de name para Name
	Department string `json:"department,omitempty"` // Alterado de department para Department
}
