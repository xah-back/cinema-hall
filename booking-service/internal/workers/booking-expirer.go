package workers

import (
	"booking-service/internal/config"
	"booking-service/internal/services"
	"time"
)

// StartExpiredBookingsWorker запускает фоновый воркер для отмены просроченных броней (15 минут)
func StartExpiredBookingsWorker(bookingService services.BookingService) {
	ticker := time.NewTicker(1 * time.Minute) // Проверка каждую минуту
	defer ticker.Stop()

	logger := config.GetLogger()
	logger.Info("Expired bookings worker started", "interval", "1 minute")

	// Запускаем сразу при старте
	if err := bookingService.ExpireOldBookings(); err != nil {
		logger.Error("Failed to expire old bookings on startup", "error", err)
	}

	// Затем каждую минуту
	for range ticker.C {
		if err := bookingService.ExpireOldBookings(); err != nil {
			logger.Error("Failed to expire old bookings", "error", err)
		}
	}
}

// StartEndedSessionsWorker запускает фоновый воркер для освобождения мест после окончания сеансов
func StartEndedSessionsWorker(bookingService services.BookingService) {
	ticker := time.NewTicker(1 * time.Minute) // Проверка каждую минуту
	defer ticker.Stop()

	logger := config.GetLogger()
	logger.Info("Ended sessions worker started", "interval", "1 minute")

	// Запускаем сразу при старте
	if err := bookingService.FreeSeatsForEndedSessions(); err != nil {
		logger.Error("Failed to free seats for ended sessions on startup", "error", err)
	}

	// Затем каждую минуту
	for range ticker.C {
		if err := bookingService.FreeSeatsForEndedSessions(); err != nil {
			logger.Error("Failed to free seats for ended sessions", "error", err)
		}
	}
}
