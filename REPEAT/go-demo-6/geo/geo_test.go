package geo_test

import (
	"demo/weather/geo"
	"testing"
)

// Пример позитивного теста
func TestGetMyLocation(t *testing.T) {
	city := "Moscow"
	expectedCity := "Moscow" // Отдельная переменная для сравнения

	got, err := geo.GetMyLocation(city)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err) // Используем Fatalf для немедленного выхода
	}

	if got.City != expectedCity {
		t.Errorf("Expected city %q, got %q", expectedCity, got.City)
	}
}

// Пример негативного теста
func TestGetMyLocationNoCity(t *testing.T) {
	city := "Londonsss"
	_, err := geo.GetMyLocation(city)
	if err != geo.ErrorNoCity {
		t.Errorf("Ожидалось %v, получение %v", geo.ErrorNoCity, err)
	}
}
