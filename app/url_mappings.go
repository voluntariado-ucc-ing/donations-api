package app

import "github.com/voluntariado-ucc-ing/donations-api/controllers"

func mapUrls() {
	router.GET("/ping", controllers.PingController.Ping)

	router.GET("/donations/get/:donation_id", controllers.DonationController.GetDonation)
	router.GET("/donations/donators", controllers.DonationController.GetDonatorByMail)
	router.POST("/donations/create", controllers.DonationController.CreateDonation)
}
