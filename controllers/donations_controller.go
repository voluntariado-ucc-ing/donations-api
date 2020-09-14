package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/voluntariado-ucc-ing/donations-api/domain"
	"github.com/voluntariado-ucc-ing/donations-api/services"
)

var (
	DonationController donationsControllerInterface = &donationController{}
)

type donationsControllerInterface interface {
	CreateDonation(c *gin.Context)
	GetDonation(c *gin.Context)
	GetDonator(c *gin.Context)
	GetAllDonations(c *gin.Context)
	EditDonation(c *gin.Context)
	DeleteDonation(c *gin.Context)
	UpdateDonationStatus(c *gin.Context)
}

type donationController struct{}

func (d donationController) UpdateDonationStatus(c *gin.Context) {
	donationId, parseErr := strconv.Atoi(c.Param("id"))
	if parseErr != nil {
		err := domain.NewBadRequestApiError("invalid donation id")
		c.JSON(err.Status(), err)
		return
	}
	var newStatus domain.StatusRequest
	if err := c.ShouldBindJSON(&newStatus); err != nil {
		err := domain.NewBadRequestApiError("invalid json body")
		c.JSON(err.Status(), err)
		return
	}
	res, err := services.DonationService.UpdateStatus(int64(donationId), newStatus)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (d donationController) GetDonator(c *gin.Context) {
	mail := c.Query("mail")
	id := c.Query("id")
	if mail == "" && id == "" {
		err := domain.NewBadRequestApiError("must pass mail or id param")
		c.JSON(err.Status(), err)
		return
	}

	var data *domain.Donor
	var err domain.ApiError
	if mail != "" {
		data, err = services.DonationService.GetDonatorByMail(mail)
		if err != nil {
			c.JSON(err.Status(), err)
			return
		}
	} else if id != "" {
		donorId, parseErr := strconv.ParseInt(id, 10, 64)
		if parseErr != nil {
			fmt.Println(parseErr)
			badRequest := domain.NewBadRequestApiError("donator id must be a number " + parseErr.Error())
			c.JSON(badRequest.Status(), badRequest)
			return
		}

		data, err = services.DonationService.GetDonatorById(donorId)
		if err != nil {
			c.JSON(err.Status(), err)
			return
		}
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
	donationId := c.Param("id")
	if donationId == "" {
		err := domain.NewBadRequestApiError("donation id must not be empty")
		fmt.Println(err)
		c.JSON(err.Status(), err)
		return
	}

	id, err := strconv.ParseInt(donationId, 10, 64)
	if err != nil {
		fmt.Println(err)
		badRequest := domain.NewBadRequestApiError("donation id must be a number " + err.Error())
		c.JSON(badRequest.Status(), badRequest)
		return
	}

	donation, dErr := services.DonationService.GetDonation(id)
	if dErr != nil {
		fmt.Println(dErr)
		c.JSON(dErr.Status(), dErr)
		return
	}

	c.JSON(http.StatusOK, donation)
}

func (d donationController) GetAllDonations(c *gin.Context) {
	userFilter, parseErr := strconv.Atoi(c.Query("user"))
	statusFilter := c.Query("status")
	typeFilter, parseErr := strconv.Atoi(c.Query("type"))
	if parseErr != nil && (c.Query("user") != "" || c.Query("type") != "") {
		br := domain.NewBadRequestApiError("query params user or type must be int64")
		c.JSON(br.Status(), br)
		return
	}
	res, err := services.DonationService.GetAllDonations(int64(userFilter), statusFilter, int64(typeFilter))
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, res)
	return
}

func (d donationController) EditDonation(c *gin.Context) {
	panic("implement me")
}

func (d donationController) DeleteDonation(c *gin.Context) {
	panic("implement me")
}
