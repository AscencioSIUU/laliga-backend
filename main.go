package main

import (
    "net/http"
    "github.com/gin-gonic/gin" // <-- Agrega este import
)

func main() {
    router := gin.Default()

    router.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "Bienvenido a La Liga Tracker Backend",
        })
    })

    // Aquí defines los demás endpoints:
    // GET /api/matches
    // GET /api/matches/:id
    // POST /api/matches
    // PUT /api/matches/:id
    // DELETE /api/matches/:id

    router.Run(":8080")
}
