package domains

type AllEmployeeResp struct {
	Results []Employee `json:"results"`
}
type Employee struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Login  string  `json:"login"`
	Salary float64 `json:"salary"`
}
