package main

import (
	"backend/config"
	"backend/middleware"
	"backend/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// konek ke database dulu cuy
	if err := config.ConnectDatabase(); err != nil {
		log.Fatal("❌ gak bisa konek database bang:", err)
	}

	// bikin gin router nya
	r := gin.Default()

	// pasang cors biar flutter nya bisa manggil api ini
	r.Use(middleware.CORSMiddleware())

	// setup semua route endpoint disini
	routes.SetupRoutes(r)

	// gas jalan servernya di port 3000
	log.Println("🚀 server udah jalan nih di http://localhost:3000")
	log.Println("📡 cek health: http://localhost:3000/api/health")
	r.Run(":3000")
}
