package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)


type Series struct {
	ID                 int    `json:"id"`
	Title              string `json:"title"`
	Status             string `json:"status"`
	LastEpisodeWatched int    `json:"lastEpisodeWatched"`
	TotalEpisodes      int    `json:"totalEpisodes"`
	Ranking            int    `json:"ranking"`
}

var db *sql.DB


func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}


func getSeries(w http.ResponseWriter, r *http.Request) {

	sortParam := r.URL.Query().Get("sort")  
	statusParam := r.URL.Query().Get("status") 
	searchParam := r.URL.Query().Get("search") 


	query := "SELECT id, title, status, episodes_watched, total_episodes, ranking FROM series WHERE 1=1"
	args := []interface{}{}


	if statusParam != "" {
		query += " AND status = ?"
		args = append(args, statusParam)
	}


	if searchParam != "" {
		query += " AND title LIKE ?"
		args = append(args, "%" + searchParam + "%")
	}


	if sortParam == "asc" {
		query += " ORDER BY ranking ASC"
	} else if sortParam == "desc" {
		query += " ORDER BY ranking DESC"
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var s Series
		
		if err := rows.Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		seriesList = append(seriesList, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(seriesList)
}


func getSeriesByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de serie inválido", http.StatusBadRequest)
		return
	}
	var s Series
	err = db.QueryRow("SELECT id, title, status, episodes_watched, total_episodes, ranking FROM series WHERE id = ?", id).
		Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Serie no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}


func createSeries(w http.ResponseWriter, r *http.Request) {
	var s Series
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := db.Exec("INSERT INTO series (title, status, episodes_watched, total_episodes, ranking) VALUES (?, ?, ?, ?, ?)",
		s.Title, s.Status, s.LastEpisodeWatched, s.TotalEpisodes, s.Ranking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}


func updateSeries(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de serie inválido", http.StatusBadRequest)
		return
	}
	var s Series
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.ID = id
	_, err = db.Exec("UPDATE series SET title = ?, status = ?, episodes_watched = ?, total_episodes = ?, ranking = ? WHERE id = ?",
		s.Title, s.Status, s.LastEpisodeWatched, s.TotalEpisodes, s.Ranking, s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func deleteSeries(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID de serie inválido", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("DELETE FROM series WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}


func incrementEpisode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE series SET episodes_watched = episodes_watched + 1 WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var s Series
	err = db.QueryRow("SELECT id, title, status, episodes_watched, total_episodes, ranking FROM series WHERE id = ?", id).
		Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}


func upvoteRanking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE series SET ranking = ranking + 1 WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var s Series
	err = db.QueryRow("SELECT id, title, status, episodes_watched, total_episodes, ranking FROM series WHERE id = ?", id).
		Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}


func downvoteRanking(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE series SET ranking = ranking - 1 WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var s Series
	err = db.QueryRow("SELECT id, title, status, episodes_watched, total_episodes, ranking FROM series WHERE id = ?", id).
		Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}


func updateStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	_, err = db.Exec("UPDATE series SET status = ? WHERE id = ?", body.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var s Series
	err = db.QueryRow("SELECT id, title, status, episodes_watched, total_episodes, ranking FROM series WHERE id = ?", id).
		Scan(&s.ID, &s.Title, &s.Status, &s.LastEpisodeWatched, &s.TotalEpisodes, &s.Ranking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func main() {
	// Configuración de conexión a la base de datos
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "password")
	dbHost := getEnv("DB_HOST", "mariadb")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "series_db")
	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error conectando a la base de datos: ", err)
	}
	defer db.Close()

	// Intentar la conexión (con reintentos)
	maxRetries := 15
	for i := 0; i < maxRetries; i++ {
		if err = db.Ping(); err == nil {
			log.Println("Conexión a la base de datos establecida")
			break
		}
		log.Printf("Intento %d: Base de datos no disponible, reintentando en 2 segundos...\n", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Error al hacer ping a la base de datos: ", err)
	}

	// Configuración del router
	r := mux.NewRouter()

	// Endpoints de la API existentes
	r.HandleFunc("/api/series", getSeries).Methods("GET")
	r.HandleFunc("/api/series/{id}", getSeriesByID).Methods("GET")
	r.HandleFunc("/api/series", createSeries).Methods("POST")
	r.HandleFunc("/api/series/{id}", updateSeries).Methods("PUT")
	r.HandleFunc("/api/series/{id}", deleteSeries).Methods("DELETE")

	// Nuevos endpoints PATCH
	r.HandleFunc("/api/series/{id}/episode", incrementEpisode).Methods("PATCH")
	r.HandleFunc("/api/series/{id}/upvote", upvoteRanking).Methods("PATCH")
	r.HandleFunc("/api/series/{id}/downvote", downvoteRanking).Methods("PATCH")
	r.HandleFunc("/api/series/{id}/status", updateStatus).Methods("PATCH")

	// Redirige /api a /api/series
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/series", http.StatusFound)
	}).Methods("GET")


	// Servir archivos estáticos (frontend) desde la carpeta "static"
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Println("Servidor iniciado en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
