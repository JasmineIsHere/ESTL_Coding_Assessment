package daos

import (
	"awesomeProject/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type EmployeesDAO interface {
	AddEmployee(exec boil.Executor, employee models.Employee) error
}

type employeesDAO struct{}

func NewEmployeesDAO() *employeesDAO {
	return &employeesDAO{}
}

func (dao *employeesDAO) AddEmployee(exec boil.Executor, employee models.Employee) error {

	err := employee.Upsert(exec, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}
	return nil
}
