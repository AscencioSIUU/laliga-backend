package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	MatchDate string `json:"matchDate"` // o usa time.Time si prefieres
	ScoreA    int    `json:"score_a"`
	ScoreB    int    `json:"score_b"`
}

var db *sql.DB

func main() {
	// Cargar variables de entorno (opcional)
	if err := godotenv.Load(); err != nil {
		log.Println("No se pudo cargar .env, se usarán variables de entorno existentes")
	}

	// Leer credenciales
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construir DSN: user:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	log.Println("Conectado a la base de datos exitosamente")

	// Inicializar Gin y habilitar CORS
	router := gin.Default()
	router.Use(cors.Default())

	// Servir archivos estáticos
	router.Static("/static", "./public")
	router.StaticFile("/", "./public/LaLigaTracker.html")

	// Rutas
	api := router.Group("/api")
	{
		matches := api.Group("/matches")
		{
			// Endpoints principales
			matches.GET("", getAllMatches)
			matches.GET("/:id", getMatchByID)
			matches.POST("", createMatch)
			matches.PUT("/:id", updateMatch)
			matches.DELETE("/:id", deleteMatch)

			// Endpoints PATCH que incrementan en 1
			matches.PATCH("/:id/goals", updateGoals)
			matches.PATCH("/:id/yellowcards", updateYellowCards)
			matches.PATCH("/:id/redcards", updateRedCards)
			matches.PATCH("/:id/extratime", updateExtraTime)
		}
	}

	router.Run(":8080")
}

// -------------------------------------------------------------------
// Handlers principales
// -------------------------------------------------------------------

func getAllMatches(c *gin.Context) {
	rows, err := db.Query("SELECT id, team_a, team_b, match_date, score_a, score_b FROM matches")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener partidos"})
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate, &m.ScoreA, &m.ScoreB); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer partidos"})
			return
		}
		matches = append(matches, m)
	}
	c.JSON(http.StatusOK, gin.H{"matches": matches})
}

func getMatchByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar el partido en la base de datos
	var m Match
	query := "SELECT id, team_a, team_b, match_date, score_a, score_b FROM matches WHERE id = ?"
	if err := db.QueryRow(query, id).Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate, &m.ScoreA, &m.ScoreB); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el partido"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"match": m})
}

func createMatch(c *gin.Context) {
	var newMatch Match
	if err := c.ShouldBindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO matches (team_a, team_b, match_date, score_a, score_b) VALUES (?, ?, ?, ?, ?)"
	result, err := db.Exec(query, newMatch.HomeTeam, newMatch.AwayTeam, newMatch.MatchDate, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el partido"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el ID"})
		return
	}
	newMatch.ID = int(id)
	c.JSON(http.StatusCreated, gin.H{"match": newMatch})
}

func updateMatch(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var updatedMatch Match
	if err := c.ShouldBindJSON(&updatedMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar en la BD (ejemplo simple, sin validations)
	query := "UPDATE matches SET team_a = ?, team_b = ?, match_date = ? WHERE id = ?"
	if _, err := db.Exec(query, updatedMatch.HomeTeam, updatedMatch.AwayTeam, updatedMatch.MatchDate, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el partido"})
		return
	}

	updatedMatch.ID = id
	c.JSON(http.StatusOK, gin.H{"match": updatedMatch})
}

func deleteMatch(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	query := "DELETE FROM matches WHERE id = ?"
	if _, err := db.Exec(query, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el partido"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partido eliminado", "id": id})
}

// -------------------------------------------------------------------
// Handlers PATCH: incrementan en 1 las columnas
// -------------------------------------------------------------------

func updateGoals(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Sumar 1 a la columna 'goals'
	query := "UPDATE matches SET goals = goals + 1 WHERE id = ?"
	if _, err := db.Exec(query, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar goles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Gol registrado correctamente", "id": id})
}

func updateYellowCards(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Sumar 1 a la columna 'yellowcards'
	query := "UPDATE matches SET yellowcards = yellowcards + 1 WHERE id = ?"
	if _, err := db.Exec(query, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar tarjetas amarillas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tarjeta amarilla registrada", "id": id})
}

func updateRedCards(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Sumar 1 a la columna 'redcards'
	query := "UPDATE matches SET redcards = redcards + 1 WHERE id = ?"
	if _, err := db.Exec(query, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar tarjetas rojas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tarjeta roja registrada", "id": id})
}

func updateExtraTime(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Sumar 1 a la columna 'extratime'
	query := "UPDATE matches SET extratime = extratime + 1 WHERE id = ?"
	if _, err := db.Exec(query, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar tiempo extra"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tiempo extra incrementado", "id": id})
}
