package main

import (
	"github.com/gmhafiz/go8gen/cmd"
)

// declaring struct
type Student struct {

	// defining struct fields
	Name string
	Marks int
	Id string
}


// main function
func main() {

	cmd.Execute()

	//// defining an object of struct
	//std1 := Student{
	//	Name: "Vani",
	//	Marks: 94, Id: "20024",
	//}
	//
	//// Parsing the required html
	//// file in same directory
	//tmpl, err := template.ParseFiles("internal/tmpl/index.html")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//// standard output to print merged data
	//err = tmpl.Execute(os.Stdout, std1)
	//if err != nil {
	//	fmt.Println(err)
	//}
}