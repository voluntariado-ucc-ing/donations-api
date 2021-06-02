package clients

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/voluntariado-ucc-ing/donations-api/config"
	"github.com/voluntariado-ucc-ing/donations-api/domain"
)

const (
	queryGetAllDonations    = "SELECT donation_id FROM donations"
	queryInsertDirection    = "INSERT INTO directions (street, number, details, city, postal_code) VALUES ($1,$2,$3,$4,$5) RETURNING direction_id"
	queryInsertDonation     = "INSERT INTO donations (quantity, unit, description, type_id, donator_id, direction_id, donation_date, status, element) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING donation_id"
	queryGetDonorByMail     = "SELECT donator_id, mail, first_name, last_name, phone_number FROM donators WHERE mail=$1"
	queryInsertDonor        = "INSERT INTO donators (mail, first_name, last_name, phone_number) VALUES ($1, $2, $3, $4) RETURNING donator_id"
	queryGetDonationById    = "SELECT d.donation_id, d.quantity, d.unit, d.description, d.type_id, d.donation_date, d.status, d.element, d.donator_id, i.direction_id, i.street, i.number, i.details, i.city, i.postal_code FROM directions i INNER JOIN donations d ON i.direction_id=d.direction_id WHERE d.donation_id=$1"
	queryGetDonorById       = "SELECT donator_id, mail, first_name, last_name, phone_number FROM donators WHERE donator_id=$1"
	queryUpdateStatusById   = "UPDATE donations SET status=$1 WHERE donation_id=$2"
	queryUpdateDonationById = "UPDATE donations SET quantity=$1, unit=$2, description=$3, type_id=$4, donator_id=$5, direction_id=$6, donation_date=$7, status=$8, element=$9 WHERE donation_id=$10"
	queryDeleteDonationById = "DELETE FROM donations WHERE donation_id=$1"
	queryEditDonorById = "UPDATE donators SET mail=$1, first_name=$2, last_name=$3, phone_number=$4 WHERE donator_id=$5"
)

var dbClient *sql.DB

func init() {
	var err error
	connStr := fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=disable",
		config.GetDatabaseHost(),
		config.GetDatabaseUser(),
		config.GetDatabasePort(),
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
	res := q.QueryRow(don.Quantity, don.Unit, don.Description, don.TypeId, don.DonorId, don.DirectionId, date, domain.StatusToBeConfirmed, don.Element)
	err = res.Scan(&id)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error scaning last insert id for create donation", err)
	}
	return id, nil
}

func GetDonatorByMail(mail string) (*domain.Donor, domain.ApiError) {
	var donor domain.Donor
	q, err := dbClient.Prepare(queryGetDonorByMail)
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

func GetDonation(id int64) (*domain.Donation, domain.ApiError) {
	var donation domain.Donation
	q, err := dbClient.Prepare(queryGetDonationById)
	if err != nil {
		return nil, domain.NewInternalServerApiError("Error preparing get donation statement", err)
	}
	res := q.QueryRow(id)
	donation.Direction = &domain.Direction{}
	err = res.Scan(&donation.DonationId, &donation.Quantity, &donation.Unit,
		&donation.Description, &donation.TypeId, &donation.Date, &donation.Status,
		&donation.Element, &donation.DonorId, &donation.Direction.DirectionId, &donation.Direction.Street,
		&donation.Direction.Number, &donation.Direction.Details, &donation.Direction.City, &donation.Direction.PostalCode)
	if err != nil {
		fmt.Print(err)
		return nil, domain.NewNotFoundApiError("Error donation not found")
	}

	return &donation, nil
}

func GetDonatorById(id int64) (*domain.Donor, domain.ApiError) {
	var donor domain.Donor
	q, err := dbClient.Prepare(queryGetDonorById)
	if err != nil {
		return nil, domain.NewInternalServerApiError("Error preparing get donator statement", err)
	}
	res := q.QueryRow(id)
	err = res.Scan(&donor.DonorId, &donor.Mail, &donor.FirstName, &donor.LastName, &donor.PhoneNumber)
	if err != nil {
		return nil, domain.NewApiError("Error donator not found", err.Error(), http.StatusNotFound, domain.CauseList{})
	}
	return &donor, nil
}

func GetAllDonationsIds() ([]int64, domain.ApiError) {
	ids := make([]int64, 0)
	q, err := dbClient.Prepare(queryGetAllDonations)
	if err != nil {
		fmt.Println(err)
		return nil, domain.NewInternalServerApiError("Error preparing get all volunteers statement", err)
	}

	res, err := q.Query()
	if err != nil {
		fmt.Println(err)
		return nil, domain.NewNotFoundApiError("no donations found")
	}

	defer res.Close()

	for res.Next() {
		var id int64
		err := res.Scan(&id)
		if err != nil {
			return nil, domain.NewNotFoundApiError("id not found in get all")
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func UpdateDonationStatus(donationId int64, status string) domain.ApiError {
	q, err := dbClient.Prepare(queryUpdateStatusById)
	if err != nil {
		fmt.Println(err)
		return domain.NewInternalServerApiError("error preparing query", err)
	}
	_, err = q.Exec(status, donationId)
	if err != nil {
		fmt.Println(err)
		return domain.NewInternalServerApiError("error updating status from db", err)
	}
	return nil
}

func UpdateDonation(donationUpdate domain.Donation) (int64, domain.ApiError) {
	panic("finish me")

	q, err := dbClient.Prepare(queryUpdateDonationById)
	if err != nil {
		fmt.Println(err)
		return 0, domain.NewInternalServerApiError("error preparing query", err)
	}
	id := donationUpdate.DonationId
	time := time.Now()
	date := fmt.Sprintf("%d-%d-%d", time.Year(), time.Month(), time.Day())
	res := q.QueryRow(donationUpdate.Quantity, donationUpdate.Unit, donationUpdate.Description, donationUpdate.TypeId, donationUpdate.DonorId, donationUpdate.DirectionId, date, domain.StatusToBeConfirmed, donationUpdate.Element)
	err = res.Scan(&id)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error scaning last insert id for create donation", err)
	}
	return id, nil
}

func DeleteDonation(donationId int64) domain.ApiError {
	panic("implement me")
	return nil
}

func EditDonorById(donorUpdate domain.Donor) (int64, domain.ApiError) {
	q, err := dbClient.Prepare(queryUpdateDonationById)
	if err != nil {
		fmt.Println(err)
		return 0, domain.NewInternalServerApiError("error preparing query", err)
	}
	id := donorUpdate.DonorId
	res := q.QueryRow(donorUpdate.Mail,donorUpdate.FirstName,donorUpdate.LastName,donorUpdate.PhoneNumber,donorUpdate.DonorId)
	err = res.Scan(&id)
	if err != nil {
		return 0, domain.NewInternalServerApiError("Error scaning last insert id for create donation", err)
	}
	return id, nil
}