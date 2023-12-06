package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID          int
	Name        string
	Price       float32
	Description string
}

var tpl *template.Template

var db *sql.DB

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")
	var err error

	db, err = sql.Open("mysql", "root:AngelinaNava$1@tcp(localhost:3306)/test222")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/browse", browseHandler)
	http.HandleFunc("/update/", updateHandler)
	http.HandleFunc("/updateresult/", updateResultHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/", homePageHandler)
	http.ListenAndServe(":8790", nil)
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("^^^Browse Handler^^^")
	stmt := "SELECT * FROM products"

	rows, err := db.Query(stmt)
	if err != nil {
		fmt.Println("blaaaadd")
	}

	defer rows.Close()
	var products []Product

	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description)
		if err != nil {
			panic(err)
		}

		products = append(products, p)
	}
	tpl.ExecuteTemplate(w, "select.html", products)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("^^^Insert Handler^^^")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "insert.html", nil)
		return
	}
	r.ParseForm()

	name := r.FormValue("nameName")
	price := r.FormValue("priceName")
	descr := r.FormValue("descrName")
	var err error
	if name == "" || price == "" || descr == "" {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "insert.html", "Error inserting DATA, please chech all fields")
		return
	}
	var ins *sql.Stmt

	ins, err = db.Prepare("INSERT INTO `test222`.`products` (`name`,`price`, `description`) VALUES(?, ?, ?);")
	if err != nil {
		panic(err)
	}
	defer ins.Close()

	res, err := ins.Exec(name, price, descr)

	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("Error inserting row:", err)
		tpl.ExecuteTemplate(w, "insert.html", "Error inserting data, please check all fields")
		return
	}

	lastInserted, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("ID of last row insertes:", lastInserted)
	fmt.Println("Number of affected rows: ", rowsAffected)
	tpl.ExecuteTemplate(w, "insert.html", "Product Successfully Inserted")

}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("^^^updateHandler running^^^")
	r.ParseForm()
	id := r.FormValue("idproducts")
	row := db.QueryRow("SELECT * FROM test222.products WHERE idproducts = ?;", id)
	var p Product

	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Description)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/browse", 307)
		return
	}
	tpl.ExecuteTemplate(w, "update.html", p)
}

func updateResultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("^^^updateResultHandler running^^^")
	r.ParseForm()
	id := r.FormValue("idproducts")
	name := r.FormValue("nameName")
	price := r.FormValue("priceName")
	description := r.FormValue("descrName")
	upStmt := "UPDATE `test222`.`products` SET `name` = ?, `price` = ?, `description` = ? WHERE (`idproducts` = ?);"

	stmt, err := db.Prepare(upStmt)
	if err != nil {
		fmt.Println("error preparing err:", err)
		panic(err)
	}
	fmt.Println("db.Prepare err:", err)
	fmt.Println("db.Prepare stmt:", stmt)
	defer stmt.Close()
	var res sql.Result

	res, err = stmt.Exec(name, price, description, id)
	rowsAff, _ := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "result.html", "There was a problem updating the product")
		return
	}
	tpl.ExecuteTemplate(w, "result.html", "Product was Successfully Updated")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("^^^deleteHandler running^^^")
	r.ParseForm()
	id := r.FormValue("idproducts")

	del, err := db.Prepare("DELETE FROM `test222`.`products` WHERE (`idproducts` = ?);")
	if err != nil {
		panic(err)
	}
	defer del.Close()
	var res sql.Result
	res, err = del.Exec(id)
	rowsAff, _ := res.RowsAffected()
	fmt.Println("rowsAff:", rowsAff)

	if err != nil || rowsAff != 1 {
		fmt.Println(w, "Error delating product")
		return
	}

	fmt.Println("err:", err)
	tpl.ExecuteTemplate(w, "result.html", "product was successfully Deleted")

}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/browse", 307)
}
