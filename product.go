package main

import(
	"fmt"
	"html/template"
	"net/http"
	"database/sql"
)

func products(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	tplValues := map[string]interface{}{"Header": "Products", "Copyright": "Mohit"}
	db, err := sql.Open("sqlite3", "file:./db/app.db?foreign_keys=true")
	if err != nil {
		serveError(w, err)
		return
	}
	defer db.Close()

	sql := "select title, description, price, quantity, filename from products order by title"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Printf("%q: %s\n", err, sql)
		serveError(w, err)
		return
	}
	defer rows.Close()

	levels := []map[string]string{}
	var title, description, price, quantity, filename string
	for rows.Next() {
		rows.Scan(&title, &description, &price, &quantity, &filename)
		levels = append(levels, map[string]string{"title": title, "description": description, "price": price, "quantity": quantity, "filename": filename})
	}
	tplValues["levels"] = levels

	rows.Close()

	pageTemplate, err := template.ParseFiles("tpl/products.html", "tpl/header.html", "tpl/bar.html", "tpl/footer.html")
	if err != nil {
		serveError(w, err)
	}

	if i, ok := session.Values["login"]; ok {
		tplValues["login"] = i
	}

	err = pageTemplate.Execute(w, tplValues)
	if err != nil {
		serveError(w, err)
	}
}
