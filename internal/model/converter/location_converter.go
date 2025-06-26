package converter

import (
	"fmt"

	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func parseLocation(wkt string) *model.LocationResponse {
	var lon, lat float64
	_, err := fmt.Sscanf(wkt, "POINT(%f %f)", &lon, &lat)
	if err != nil {
		return nil
	}
	return &model.LocationResponse{
		Latitude:  lat,
		Longitude: lon,
	}
}
