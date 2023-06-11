package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Spot characteristics
type Spot struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Website     *string `json:"website"`
	Description *string `json:"description"`
	Rating      float64 `json:"rating"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:rootroot@/spotlas")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/spots", getSpots).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getSpots(w http.ResponseWriter, r *http.Request) {
	latitude, err := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	// error handling in case of invalid input
	if err != nil {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	longitude, err := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	if err != nil {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}
	radius, err := strconv.ParseFloat(r.URL.Query().Get("radius"), 64)
	if err != nil {
		http.Error(w, "Invalid radius", http.StatusBadRequest)
		return
	}

	spotType := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("type")))

	longitudeRadius := radius / (111412.84*math.Cos(latitude) - 93.5*math.Cos(3*latitude))
	latitudeRadius := radius / 111000

	var query string
	if spotType == "circle" {
		query = fmt.Sprintf(`
			SELECT id, name, website, ST_X(coordinates) as latitude, ST_Y(coordinates) as longitude, description, rating
			FROM MY_TABLE
			WHERE ST_Distance_Sphere(coordinates, ST_GeomFromText('POINT(%f %f)')) <= %f
			ORDER BY 
			CASE 
				WHEN ST_Distance_Sphere(coordinates, ST_GeomFromText('POINT(%f %f)')) <= 50 THEN rating
				ELSE ST_Distance_Sphere(coordinates, ST_GeomFromText('POINT(%f %f)'))
			END ASC`,
			latitude, longitude, radius, latitude, longitude, latitude, longitude)
	} else if spotType == "square" {
		query = fmt.Sprintf(`
			SELECT id, name, website, ST_X(coordinates) as latitude, ST_Y(coordinates) as longitude, description, rating
			FROM MY_TABLE
			WHERE MBRContains(
				ST_GeomFromText('Polygon((%f %f, %f %f, %f %f, %f %f, %f %f))'),
				coordinates)
			ORDER BY 
			CASE 
				WHEN ST_Distance_Sphere(coordinates, ST_GeomFromText('POINT(%f %f)')) <= 50 THEN rating
				ELSE ST_Distance_Sphere(coordinates, ST_GeomFromText('POINT(%f %f)'))
			END ASC`,
			latitude-latitudeRadius, longitude-longitudeRadius,
			latitude+latitudeRadius, longitude-longitudeRadius,
			latitude+latitudeRadius, longitude+longitudeRadius,
			latitude-latitudeRadius, longitude+longitudeRadius,
			latitude-latitudeRadius, longitude-longitudeRadius,
			latitude, longitude, latitude, longitude)
	} else {
		http.Error(w, "Invalid type parameter", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spots := make([]Spot, 0)
	for rows.Next() {
		var spot Spot
		err := rows.Scan(&spot.ID, &spot.Name, &spot.Website, &spot.Latitude, &spot.Longitude, &spot.Description, &spot.Rating)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		spots = append(spots, spot)
	}

	json.NewEncoder(w).Encode(spots)
}
