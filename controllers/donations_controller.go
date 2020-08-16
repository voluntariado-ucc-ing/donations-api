package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/voluntariado-ucc-ing/donations-api/domain"
	"github.com/voluntariado-ucc-ing/donations-api/services"
	"net/http"
)

var(
	DonationController donationsControllerInterface = &donationController{}
)

type donationsControllerInterface interface {
	CreateDonation(c *gin.Context)
	GetDonation(c *gin.Context)
	GetDonatorByMail(c *gin.Context)
	GetAllDonations(c *gin.Context)
	EditDonation(c *gin.Context)
	DeleteDonation(c *gin.Context)
}

type donationController struct {}

func (d donationController) GetDonatorByMail(c *gin.Context) {
	mail := c.Query("mail")
	if mail == "" {
		err := domain.NewBadRequestApiError("must pass mail param")
		c.JSON(err.Status(), err)
		return
	}

	data, err := services.DonationService.GetDonatorByMail(mail)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, data)
}


func (d donationController) CreateDonation(c *gin.Context) {
	var donationRequest domain.DonationRequest
	if err := c.ShouldBindJSON(&donationRequest); err != nil {
		fmt.Println(err)
		err := domain.NewBadRequestApiError("Invalid donation body")
		c.JSON(err.Status(), err)
		return
	}

	r, err := services.DonationService.CreateDonation(donationRequest)
	if err != nil {
		fmt.Println(err)
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, r)
}

func (d donationController) GetDonation(c *gin.Context) {
	panic("implement me")
}

func (d donationController) GetAllDonations(c *gin.Context) {
	panic("implement me")
}

func (d donationController) EditDonation(c *gin.Context) {
	panic("implement me")
}

func (d donationController) DeleteDonation(c *gin.Context) {
	panic("implement me")
}
