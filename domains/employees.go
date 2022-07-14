package domains

type (
	AllEmployeeResp struct {
		Results []Employee `json:"results"`
	}

	Employee struct {
		ID     string  `json:"id"`
		Name   string  `json:"name"`
		Login  string  `json:"login"`
		Salary float64 `json:"salary"`
	}
)

type EmployeeReqResp struct {
	Name   string  `json:"name"`
	Login  string  `json:"login"`
	Salary float64 `json:"salary"`
}
