package employees

import (
	"awesomeProject/daos"
	"awesomeProject/models"
	"awesomeProject/utils/db"
	"bufio"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Handler interface {
	RouterGroup(engine *gin.Engine)
}

type employeeHandler struct {
	employeesDAO daos.EmployeesDAO
}

func NewHandler(employeeDAO daos.EmployeesDAO) *employeeHandler {
	return &employeeHandler{
		employeeDAO,
	}
}

func (h *employeeHandler) RouteGroup(r *gin.Engine) {
	rg := r.Group("/users")
	rg.POST("/upload", h.uploadCSV)
	rg.GET("", h.get)
}

func (h *employeeHandler) uploadCSV(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]

	var employeesAdded int
	for _, file := range files {
		csv, err := file.Open()
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		}
		success, err := h.ProcessCSV(csv)
		employeesAdded += success

		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		}
		csv.Close()
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Number of employees inserted : %v", employeesAdded)})
}

func (h *employeeHandler) ProcessCSV(file multipart.File) (int, error) {
	br := bufio.NewReader(file)
	employeesAdded := 0

	if err := db.WithTxn(func(txn boil.Transactor) (err error) {
		for {
			s, eofErr := br.ReadString('\n')
			if eofErr != nil && !errors.Is(eofErr, io.EOF) {
				return eofErr
			}

			if len(s) > 0 && s[0] != '#' {
				cols := strings.Split(s, ",")

				if len(cols) != 4 {
					return errors.New(fmt.Sprintf("Missing employee fields: ID, login, name and salary fields are all required"))
				}

				salary, err := strconv.ParseFloat(strings.TrimSuffix(cols[3], "\n"), 64)
				if err != nil || salary < 0 {
					return errors.New(fmt.Sprintf("Invalid employee field: Salary should be a decimal that is > 0.0 for employee where id = %v", cols[0]))
				}

				if err := h.employeesDAO.AddEmployee(txn, models.Employee{
					ID:     cols[0],
					Login:  cols[1],
					Name:   cols[2],
					Salary: null.Float64From(salary),
				}); err != nil {
					return err
				}
				employeesAdded++
			}
			if errors.Is(eofErr, io.EOF) {
				break
			}
		}
		return
	}); err != nil {
		return 0, err
	}
	if employeesAdded == 0 {
		return 0, errors.New(fmt.Sprintf("Employees Added is 0 : empty file was uploaded"))
	}

	return employeesAdded, nil
}

func (h *employeeHandler) get(c *gin.Context) {
	var minSalary, maxSalary null.Float64
	var sort, order null.String
	var limit, offset int

	minSalaryString, present := c.GetQuery("minSalary")
	if present && minSalaryString != "" {
		minSalaryFloat64, err := strconv.ParseFloat(minSalaryString, 64)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		}
		minSalary = null.Float64From(minSalaryFloat64)
	}

	maxSalaryString, present := c.GetQuery("maxSalary")
	if present && maxSalaryString != "" {
		maxSalaryFloat64, err := strconv.ParseFloat(maxSalaryString, 64)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		}
		maxSalary = null.Float64From(maxSalaryFloat64)
	}

	sortString, present := c.GetQuery("sort")
	if present && sortString != "" { // TODO: assuming the value "+name" is passed in as "%2Bname" (add in README)
		symbol := sortString[:1]
		if symbol == "+" {
			order = null.StringFrom("asc")
		} else if symbol == "-" {
			order = null.StringFrom("desc")
		} else {
			c.Error(errors.New("Invalid data format: Order should be represented by %2B (+) (ascending) or - (descending)"))
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		}

		col := sortString[1:]
		if col != "id" && col != "name" && col != "login" && col != "salary" {
			// invalid sort key
			c.Error(errors.New("Invalid data format: Only columns \"id\", \"name\", \"login\" or \"salary\" can be sorted"))
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		} else {
			sort = null.StringFrom(col)
		}
	}

	var err error
	limitString, present := c.GetQuery("limit")
	if present && limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			c.Error(errors.New("Invalid data format: limit should be an integer"))
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		}
	}

	offsetString, present := c.GetQuery("offset")
	if present && offsetString != "" {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			c.Error(errors.New("Invalid data format: offset should be an integer"))
			c.JSON(http.StatusBadRequest, c.Errors.Last())
			return
		}
	}

	employeeSlice, err := h.employeesDAO.GetAll(boil.GetDB(), minSalary, maxSalary, sort, order, limit, offset)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, c.Errors.Last())
	}
	c.JSON(http.StatusOK, &employeeSlice)
}
