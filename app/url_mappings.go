package app

import "github.com/voluntariado-ucc-ing/donations-api/controllers"

func mapUrls() {
	router.GET("/ping", controllers.PingController.Ping)

	router.GET("/donations/get/:id", controllers.DonationController.GetDonation)
	router.GET("/donations/donators", controllers.DonationController.GetDonator)
	router.POST("/donations/create", controllers.DonationController.CreateDonation)
	router.GET("/donations/all", controllers.DonationController.GetAllDonations)
	router.PATCH("/donations/:id", controllers.DonationController.UpdateDonationStatus)
}
