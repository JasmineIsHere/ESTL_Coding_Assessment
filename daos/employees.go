package daos

import (
	"awesomeProject/domains"
	"awesomeProject/models"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type EmployeesDAO interface {
	AddEmployee(exec boil.Executor, employee models.Employee) error
	DeleteEmployee(exec boil.Executor, empID string) error
	GetAll(exec boil.Executor, minSalary null.Float64, maxSalary null.Float64, sort null.String, order null.String, limit int, offset int) (*models.EmployeeSlice, error)
	GetByID(exec boil.Executor, empID string) (*models.Employee, error)
	UpdateEmployee(exec boil.Executor, employee domains.EmployeeReqResp, empID string) error
	UpsertEmployee(exec boil.Executor, employee models.Employee) error
}

type employeesDAO struct{}

func NewEmployeesDAO() *employeesDAO {
	return &employeesDAO{}
}
func (dao *employeesDAO) AddEmployee(exec boil.Executor, employee models.Employee) error {
	err := employee.Insert(exec, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (dao *employeesDAO) DeleteEmployee(exec boil.Executor, empID string) error {
	employeeInDB, err := dao.GetByID(exec, empID)
	if err != nil {
		return err
	}
	_, err = employeeInDB.Delete(exec)
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

func (dao *employeesDAO) GetByID(exec boil.Executor, empID string) (*models.Employee, error) {
	employee, err := models.Employees(models.EmployeeWhere.ID.EQ(empID)).One(exec)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (dao *employeesDAO) UpdateEmployee(exec boil.Executor, employee domains.EmployeeReqResp, empID string) error {
	employeeInDB, err := dao.GetByID(exec, empID)
	if err != nil {
		return err
	}
	// assumption: all fields will be passed in the body, regardless of whether value was modified
	employeeInDB.Login = employee.Login
	employeeInDB.Name = employee.Name
	employeeInDB.Salary = null.Float64From(employee.Salary)

	_, err = employeeInDB.Update(exec, boil.Infer())
	if err != nil {
		return err
	}
	return nil
}

func (dao *employeesDAO) UpsertEmployee(exec boil.Executor, employee models.Employee) error {
	err := employee.Upsert(exec, boil.Infer(), boil.Infer())
	if err != nil {
		return err
	}
	return nil
}
