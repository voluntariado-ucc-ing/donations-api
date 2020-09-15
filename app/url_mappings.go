package app

import "github.com/voluntariado-ucc-ing/donations-api/controllers"

func mapUrls() {
	router.GET("/ping", controllers.PingController.Ping)

	router.GET("/donations/get/:id", controllers.DonationController.GetDonation)
	router.GET("/donations/donators", controllers.DonationController.GetDonator)
	router.GET("/donations/all", controllers.DonationController.GetAllDonations)

	router.POST("/donations/create", controllers.DonationController.CreateDonation)

	router.PATCH("/donations/:id", controllers.DonationController.UpdateDonationStatus)
}
