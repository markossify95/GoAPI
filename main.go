package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"

	m "./model" 
)


func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/employees/", GetEmployeesList).Methods("GET")
	router.HandleFunc("/employees/{id:[0-9]+}/", GetEmployee).Methods("GET")
	router.HandleFunc("/stats/employees_by_depts/", GetEmployeeByDeptStats).Methods("GET")
	router.HandleFunc("/stats/gender_by_depts/{dept_name:[a-zA-Z0-9_]+}/", GetGenderByDeptStats).Methods("GET")

	//cors fix zbog ajax-a
	log.Fatal(http.ListenAndServe(":8600", router))
}

func GetEmployeesList(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@/employees")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	rows, _ := db.Query("select	* from employees LIMIT 1000")
	Response := make([]m.Employee, 0, 1000)
	for rows.Next() {
		emp := m.Employee{}
		rows.Scan(&emp.ID, &emp.BirthDate, &emp.Firstname, &emp.Lastname, &emp.Gender, &emp.Hiredate)
		Response = append(Response, emp)
	}



	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(Response)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	//db connecting
	urlParams := mux.Vars(r)
	id := urlParams["id"]
	db, err := sql.Open("mysql", "root:root@/employees")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	// var emp Employee
	emp := m.Employee{}
	err = db.QueryRow("SELECT emp_no, birth_date,first_name,last_name,gender,hire_date FROM employees WHERE emp_no = ?", id).Scan(
		&emp.ID, &emp.BirthDate, &emp.Firstname, &emp.Lastname, &emp.Gender, &emp.Hiredate)

	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(emp)
}

func GetEmployeeByDeptStats(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Pragma","no-cache")
	//ajax metoda, sredjujemo CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, err := sql.Open("mysql", "root:root@/employees")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	q := `SELECT d.dept_name, count(e.emp_no) FROM employees as e 
			JOIN dept_emp as d_e 
			ON d_e.emp_no = e.emp_no 
			JOIN departments as d 
			ON d.dept_no = d_e.dept_no
			GROUP BY d.dept_no`

	rows, _ := db.Query(q)
	Response := make([][]interface{}, 0, 100)
	for rows.Next() {
		res := m.EmpDeptStat{}
		rows.Scan(&res.DeptName, &res.EmpCount)
		final := make([]interface{}, 0, 2)
		final = append(final, res.DeptName)
		final = append(final,res.EmpCount)
		Response = append(Response, final)
	}

	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(Response)
}


func GetGenderByDeptStats(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Pragma","no-cache")
	//ajax metoda, sredjujemo CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	urlParams := mux.Vars(r)
	dept_name := urlParams["dept_name"]

	dept_name = strings.Replace(dept_name, "_", " ", -1)

	db, err := sql.Open("mysql", "root:root@/employees")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	q := `SELECT e.gender, count(e.emp_no) 
	FROM departments as d
	JOIN dept_emp as d_e ON d.dept_no = d_e.dept_no
	JOIN employees as e on e.emp_no = d_e.emp_no
	WHERE d.dept_name = ?
	GROUP BY e.gender`

	rows, _ := db.Query(q, dept_name)
	Response := make([][]interface{}, 0, 100)
	for rows.Next() {
		var gen string
		var count int
		rows.Scan(&gen, &count)
		final := make([]interface{}, 0, 2)
		final = append(final, gen)
		final = append(final, count)
		Response = append(Response, final)
	}

	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(Response)
}


//pretvara svaku struct u niz zbog glupog google chart
// SELECT e.emp_no, e.birth_date, e.first_name, e.last_name, e.gender, e.hire_date, d.dept_name
// FROM employees as e 
// JOIN dept_emp as d_e 
// ON d_e.emp_no = e.emp_no 
// JOIN departments as d 
// ON d.dept_no = d_e.dept_no
// LIMIT 1000
