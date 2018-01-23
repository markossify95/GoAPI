package model


type Employee struct {
	ID        string `json: "id,omitempty"`
	BirthDate string `json: "birthdate,omitempty"`
	Firstname string `json: "firstname,omitempty"`
	Lastname  string `json: "lastname,omitempty"`
	Gender    string `json: "gender,omitempty"`
	Hiredate  string `json: "hiredate,omitempty"`
}

type EmpDeptStat struct {
	DeptName string `json:deptname, omitempty`
	EmpCount int `json:empcount, omitempty`
}