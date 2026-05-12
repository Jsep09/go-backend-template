// Server + Graceful Shutdown
package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
)

func startServer(app *fiber.App, port string) {
	// Channel รับ OS signal (Ctrl+C, kill command)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server ใน goroutine แยก (ไม่ block main)
	go func() {
		addr := ":" + port
		slog.Info("server listening",
			"port", port,
			"url", "http://localhost:"+port,
			"swagger", "http://localhost:"+port+"/swagger/index.html",
		)

		if err := app.Listen(addr); err != nil {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// Block จนกว่าจะได้รับ signal
	<-quit
	slog.Info("shutting down server...")

	// Graceful shutdown: รอ request ที่ค้างอยู่ให้เสร็จก่อน (max 5 วิ)
	if err := app.ShutdownWithTimeout(5 * time.Second); err != nil {
		slog.Error("forced shutdown", "error", err)
	}

	slog.Info("server stopped")
}
