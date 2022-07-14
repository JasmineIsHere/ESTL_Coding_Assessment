package daos

import (
	"awesomeProject/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type EmployeesDAO interface {
	AddEmployee(exec boil.Executor, employee models.Employee) error
	GetAll(exec boil.Executor, minSalary null.Float64, maxSalary null.Float64, sort null.String, order null.String, limit int, offset int) (*models.EmployeeSlice, error)
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

func (dao *employeesDAO) GetAll(exec boil.Executor, minSalary null.Float64, maxSalary null.Float64, sort null.String, order null.String, limit int, offset int) (*models.EmployeeSlice, error) {
	var queryMods []qm.QueryMod

	if !minSalary.IsZero() {
		queryMods = append(queryMods, models.EmployeeWhere.Salary.GTE(minSalary))
	}

	if !maxSalary.IsZero() {
		queryMods = append(queryMods, models.EmployeeWhere.Salary.LTE(maxSalary))
	}

	if !sort.IsZero() && !order.IsZero() {
		queryMods = append(queryMods, qm.OrderBy(sort.String+" "+order.String))
	}

	queryMods = append(queryMods,
		qm.Limit(limit),
		qm.Offset(offset),
	)

	employeeSlice, err := models.Employees(queryMods...).All(exec)
	if err != nil {
		return nil, err
	}
	return &employeeSlice, nil
}
