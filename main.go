package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // Driver para MySQL
	"github.com/joho/godotenv"
)

// Match representa la estructura de un partido.
type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	MatchDate string `json:"matchDate"` // O usa time.Time si prefieres
	ScoreA    int    `json:"score_a"`
	ScoreB    int    `json:"score_b"`
}

var db *sql.DB

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, se usarán variables de entorno existentes")
	}
	// Obtener variables de entorno usando os.Getenv
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construir DSN: user:password@tcp(host:port)/dbname?parseTime=true
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	log.Println("Conectado a la base de datos exitosamente")

	// Inicializar Gin
	router := gin.Default()

	// Servir archivos estáticos desde la carpeta "public"
	router.Static("/static", "./public")
	// Servir directamente el index (frontend)
	router.StaticFile("/", "./public/LaLigaTracker.html")

	// Configuración de rutas API
	api := router.Group("/api")
	{
		matches := api.Group("/matches")
		{
			matches.GET("", getAllMatches)
			matches.GET("/:id", getMatchByID)
			matches.POST("", createMatch)
			matches.PUT("/:id", updateMatch)
			matches.DELETE("/:id", deleteMatch)
		}
	}

	router.Run(":8080")
}

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
		err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.MatchDate, &m.ScoreA, &m.ScoreB)
		if err != nil {
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
	// Aquí consultarías la BD y devolverías el partido correspondiente.
	c.JSON(http.StatusOK, gin.H{"match": Match{ID: id}})
}

func createMatch(c *gin.Context) {
	var newMatch Match
	if err := c.ShouldBindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Aquí insertarías newMatch en la base de datos, por ejemplo:
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
	updatedMatch.ID = id
	// Aquí actualizarías el registro en la BD.
	c.JSON(http.StatusOK, gin.H{"match": updatedMatch})
}

func deleteMatch(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}
	// Aquí eliminarías el partido de la BD.
	c.JSON(http.StatusOK, gin.H{"message": "Partido eliminado", "id": id})
}
