package gohtmltopdf

import (
	"encoding/json"
	"time"
)

type ErrorProcess struct {
	Msg string `json:"msg"`
}

func (e ErrorProcess) Error() string {
	return e.Msg
}

type requestHTML struct {
	// TODO agregar propiedades como: tamaño de página, número de páginas, etc

	// Data must be a string with HTML format.
	Data string `json:"data"`
}

type requestDIANForm220 struct {
	Data DIANForms220Relation `json:"data"`
}

type DIANForm220 struct {
	ID            uint            `json:"id"`
	EmployerID    uint            `json:"employer_id"`
	Year          uint            `json:"year"`
	Sequence      uint            `json:"sequence"`
	ContractID    uint            `json:"contract_id"`
	BeginsAt      time.Time       `json:"begins_at"`
	EndsAt        time.Time       `json:"ends_at"`
	Records       json.RawMessage `json:"rows"`
	AverageSalary float64         `json:"average_salary"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`

	// Fields only for logic
	RowsMap map[string]float64
}

type DIANForm220Relation struct {
	DIANForm220

	// Employer
	Nit              string
	Dv               string
	BusinessName     string
	DepartmentCode   string
	MunicipalityCode string
	Place            string

	// Employee
	IdentificationTypeID   uint
	IdentificationTypeCode uint
	IdentificationNumber   string
	FirstName              string
	MiddleName             string
	LastName               string
	Surname                string
}

type DIANForms220Relation []DIANForm220Relation
