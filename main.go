package main

import (
	"apiMysql/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	http.HandleFunc("/insert", InsertData)

	// Cargar las variables de entorno desde el archivo .properties
	if _, err := os.Stat("config/.properties"); err == nil {
		godotenv.Load("config/.properties")
	}

	host := os.Getenv("AUXHOST")
	if host == "" {
		host = "0.0.0.0" // Host predeterminado
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Puerto predeterminado
	}

	//inicio servidor
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

func InsertData(w http.ResponseWriter, r *http.Request) {
	// Verificar que el método sea POST
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Método no permitido"+r.Method, http.StatusMethodNotAllowed)
	// 	return
	// }

	// Decodificar el JSON del cuerpo de la solicitud
	var requestData models.Request
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Error al decodificar JSON, err="+err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener las variables de entorno
	// dbURL := os.Getenv("MYSQL_URL")
	dbUser := os.Getenv("MYSQLUSER")
	dbPassword := os.Getenv("MYSQLPASSWORD")
	dbHost := os.Getenv("MYSQLHOST")
	dbPort := os.Getenv("MYSQLPORT")
	// dbURL=os.Getenv("MYSQL_URL")
	// Crear la cadena de conexión
	// connectionString := fmt.Sprintf("%s:%s@tcp(%s)/dbname", dbUser, dbPassword, dbURL)

	// Conectar a la base de datos MySQL

	dbURL := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/railway"
	db, err := sql.Open("mysql", dbURL)
	// log.Printf("url %v", dbURL)
	if err != nil {
		log.Printf("Error al conectar a la base de datos: %v", err)
		http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Preparar la consulta SQL
	query := "INSERT INTO registros (regi_descripcion) VALUES (?)"
	result, err := db.Exec(query, requestData.Info)
	if err != nil {
		http.Error(w, "Error al insertar datos en la base de datos", http.StatusInternalServerError)
		return
	}

	// Obtener el ID del registro insertado
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error al obtener el ID del registro insertado", http.StatusInternalServerError)
		return
	}

	// Devolver el ID del registro insertado como respuesta
	response := fmt.Sprintf("Se insertó el registro con ID %d", lastInsertID)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": response})
}
