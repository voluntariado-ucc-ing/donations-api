package clients

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/voluntariado-ucc-ing/donations-api/config"
	"github.com/voluntariado-ucc-ing/donations-api/domain"
	"log"
	"net/http"
	"time"
)

const (
	queryInsertDirection = "INSERT INTO voluntariado_ing.directions (street, number, details, city, postal_code) VALUES ($1,$2,$3,$4,$5) RETURNING direction_id"
	queryInsertDonation  = "INSERT INTO voluntariado_ing.donations (quantity, unit, description, type_id, donator_id, direction_id, donation_date, status, element) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING donation_id"
	queryGetDonorById    = "SELECT donator_id, mail, first_name, last_name, phone_number FROM voluntariado_ing.donators WHERE mail=$1"
	queryInsertDonor     = "INSERT INTO voluntariado_ing.donators (mail, first_name, last_name, phone_number) VALUES ($1, $2, $3, $4) RETURNING donator_id"
)

var dbClient *sql.DB

func init() {
	var err error
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		config.GetDatabaseHost(),
		config.GetDatabaseUser(),
		config.GetDatabasePassword(),
		config.GetDatabaseName())

	dbClient, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := dbClient.Ping(); err != nil {
		log.Fatal(err)
	}
}

func InsertDirection(dir domain.Direction) (int64, domain.ApiError) {
	var id int64
	q, err := dbClient.Prepare(queryInsertDirection)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error preparing insert statement", err)
	}
	res := q.QueryRow(dir.Street, dir.Number, dir.Details, dir.City, dir.PostalCode)
	err = res.Scan(&id)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error scaning last insert id for create details", err)
	}
	return id, nil
}

func InsertDonation(don domain.Donation) (int64, domain.ApiError) {
	var id int64
	q, err := dbClient.Prepare(queryInsertDonation)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error preparing insert statement", err)
	}
	time := time.Now()
	date := fmt.Sprintf("%d-%d-%d", time.Year(), time.Month(), time.Day())
	res := q.QueryRow(don.Quantity, don.Unit, don.Description, don.TypeId, don.DonorId, don.DirectionId, date, "active", don.Element)
	err = res.Scan(&id)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error scaning last insert id for create donation", err)
	}
	return id, nil
}

func GetDonator(mail string) (*domain.Donor, domain.ApiError) {
	var donor domain.Donor
	q, err := dbClient.Prepare(queryGetDonorById)
	if err != nil {
		return nil, domain.NewInternalServerApiError("Error preparing get donator statement", err)
	}
	res := q.QueryRow(mail)
	err = res.Scan(&donor.DonorId, &donor.Mail, &donor.FirstName, &donor.LastName, &donor.PhoneNumber)
	if err != nil {
		return nil, domain.NewApiError("Error donator not found", err.Error(), http.StatusNotFound, domain.CauseList{})
	}
	return &donor, nil
}

func InsertDonor(don domain.Donor) (int64, domain.ApiError) {
	var id int64
	q, err := dbClient.Prepare(queryInsertDonor)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error preparing insert statement", err)
	}
	res := q.QueryRow(don.Mail, don.FirstName, don.LastName, don.PhoneNumber)
	err = res.Scan(&id)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error scaning last insert id for create donator", err)
	}
	return id, nil
}
