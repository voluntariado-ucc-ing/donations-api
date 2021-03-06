package services

import (
	"fmt"
	"net/http"

	"github.com/voluntariado-ucc-ing/donations-api/clients"
	"github.com/voluntariado-ucc-ing/donations-api/domain"
)

var (
	DonationService donationServiceInterface = &donationService{}
)

type donationServiceInterface interface {
	CreateDonation(request domain.DonationRequest) (*domain.DonationRequest, domain.ApiError)
	CreateDonator(request domain.DonatorRequest) (*domain.DonatorRequest, domain.ApiError)
	GetDonatorByMail(mail string) (*domain.Donor, domain.ApiError)
	GetDonation(id int64) (*domain.Donation, domain.ApiError)
	GetDonatorById(id int64) (*domain.Donor, domain.ApiError)
	GetAllDonations(userFilter int64, statusFilter string, typeFilter int64) ([]domain.Donation, domain.ApiError)
	UpdateStatus(donationId int64, request domain.StatusRequest) (*domain.Donation, domain.ApiError)
}

type donationService struct{}

func (d donationService) UpdateStatus(donationId int64, request domain.StatusRequest) (*domain.Donation, domain.ApiError) {
	if !request.IsValidStatus() {
		return nil, domain.NewBadRequestApiError("Invalid status for donation")
	}

	err := clients.UpdateDonationStatus(donationId, request.Status)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return clients.GetDonation(donationId)
}

func (d donationService) GetDonation(id int64) (*domain.Donation, domain.ApiError) {
	return clients.GetDonation(id)
}

func (d donationService) CreateDonation(request domain.DonationRequest) (*domain.DonationRequest, domain.ApiError) {
	var donorId int64 = 0
	donor, err := clients.GetDonatorByMail(request.Donor.Mail)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			fmt.Println(err)
			return nil, err
		}
		donorId, err = clients.InsertDonor(request.Donor)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else {
		donorId = donor.DonorId
	}

	for index := range request.Donations {
		directionId, err := clients.InsertDirection(*request.Donations[index].Direction)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		request.Donations[index].DirectionId = directionId
		request.Donations[index].DonorId = donorId

		donationId, err := clients.InsertDonation(request.Donations[index])
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		request.Donations[index].DonationId = donationId
		request.Donations[index].Direction = nil
	}

	return &request, nil
}


func (d donationService) CreateDonator(request domain.DonatorRequest) (*domain.DonatorRequest, domain.ApiError) {

	_, err := clients.GetDonatorByMail(request.Donor.Mail)
	if err != nil {
		if err.Status() != http.StatusNotFound {
			fmt.Println(err)
			return nil, err
		}
		_, err = clients.InsertDonor(request.Donor)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return &request, nil
}


func (d donationService) GetDonatorByMail(mail string) (*domain.Donor, domain.ApiError) {
	donor, err := clients.GetDonatorByMail(mail)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return donor, nil
}

func (d donationService) GetDonatorById(id int64) (*domain.Donor, domain.ApiError) {
	donor, err := clients.GetDonatorById(id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return donor, nil
}

func (d donationService) GetAllDonations(userFilter int64, statusFilter string, typeFilter int64) ([]domain.Donation, domain.ApiError) {
	donationsList := make([]domain.Donation, 0)
	ids, err := clients.GetAllDonationsIds()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	input := make(chan domain.DonationConcurrent)
	defer close(input)
	for _, id := range ids {
		go d.getConcurrentDonation(id, input)
	}

	for i := 0; i < len(ids); i++ {
		result := <-input
		if result.Error != nil {
			fmt.Println(result.Error	)
			return nil, result.Error
		}
		donationsList = append(donationsList, *result.Donation)
	}

	filteredResult := make([]domain.Donation, 0)
	if statusFilter != "" || typeFilter != 0 || userFilter != 0 {
		for i := range donationsList {
			if statusFilter != "" {
				if typeFilter != 0 {
					if userFilter != 0 {
						if donationsList[i].Status == statusFilter &&
							donationsList[i].TypeId == typeFilter &&
							donationsList[i].DonorId == userFilter {
							filteredResult = append(filteredResult, donationsList[i])
						}
					} else {
						if donationsList[i].Status == statusFilter &&
							donationsList[i].TypeId == typeFilter {
							filteredResult = append(filteredResult, donationsList[i])
						}
					}
				} else {
					if donationsList[i].Status == statusFilter {
						filteredResult = append(filteredResult, donationsList[i])
					}
				}
			} else if typeFilter != 0 {
				if userFilter != 0 {
					if donationsList[i].TypeId == typeFilter &&
						donationsList[i].DonorId == userFilter {
						filteredResult = append(filteredResult, donationsList[i])
					}
				} else {
					if donationsList[i].TypeId == typeFilter {
						filteredResult = append(filteredResult, donationsList[i])
					}
				}
			} else if userFilter != 0 {
				if donationsList[i].DonorId == userFilter {
					filteredResult = append(filteredResult, donationsList[i])
				}
			}
		}
		return filteredResult, nil
	}

	return donationsList, nil
}

func (d donationService) getConcurrentDonation(id int64, output chan domain.DonationConcurrent) {
	vol, err := d.GetDonation(id)
	output <- domain.DonationConcurrent{Donation: vol, Error: err}
	return
}
