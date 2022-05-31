package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fakovacic/amadeus"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	client, err := amadeus.New(
		os.Getenv("AMADEUS_CLIENT_ID"),
		os.Getenv("AMADEUS_CLIENT_SECRET"),
		os.Getenv("AMADEUS_ENV"), // set to "TEST"
	)
	if err != nil {
		fmt.Println("not expected error while creating client", err)
	}

	// get offer flights request&response
	offerReq, offerResp, err := client.NewRequest(amadeus.ShoppingFlightOffers)

	// set offer flights params
	offerReq.(*amadeus.ShoppingFlightOffersRequest).SetCurrency("USD").SetSources("GDS").Return(
		"LAX",
		"NYC",
		time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		time.Now().AddDate(0, 2, 0).Format("2006-01-02"),
	).AddTravelers(1, 0, 0)

	// send request
	err = client.Do(offerReq, &offerResp, "GET")

	// fmt.Println(" ----- offerResp ------")
	// get response
	offerRespData := offerResp.(*amadeus.ShoppingFlightOffersResponse)
	// fmt.Println(offerRespData.Data)
	// fmt.Println(" ----- offerResp end ------")

	// get pricing request&response
	pricingReq, pricingResp, err := client.NewRequest(amadeus.ShoppingFlightPricing)

	// add offer from flight offers response
	// this should really be a parsed as json
	// and offered as a random id
	pricingReq.(*amadeus.ShoppingFlightPricingRequest).AddOffer(
		// use id 250 to avoid "No fare applicable" error
		offerRespData.GetOffer(249), // index count begins with zero
	)

	err = client.Do(pricingReq, &pricingResp, "POST")

	// get pricing request&response
	pricingRequest, pricingResponse, err := client.NewRequest(amadeus.ShoppingFlightPricing)

	// add offer from flight offers response
	// fmt.Println(" ----- pricingRequest start ------")
	pricingRequest.(*amadeus.ShoppingFlightPricingRequest).AddOffer(
		offerRespData.GetOffer(249),
	)
	// fmt.Println(pricingRequest)
	// fmt.Println(" ----- pricingRequest end ------")

	// send request
	err = client.Do(pricingRequest, &pricingResponse, "POST")

	// get response
	pricingRespData := pricingResponse.(*amadeus.ShoppingFlightPricingResponse)

	// get booking request
	bookingReq, bookingResp, err := client.NewRequest(amadeus.BookingFlightOrder)

	// add offer from flight pricing response
	bookingReq.(*amadeus.BookingFlightOrderRequest).AddOffers(
		pricingRespData.GetOffers(),
	).AddTicketingAgreement("DELAY_TO_CANCEL", "6D")

	// println(pricingRespData.GetOffers())

	// add payment
	bookingReq.(*amadeus.BookingFlightOrderRequest).AddPayment(
		bookingReq.(*amadeus.BookingFlightOrderRequest).
			NewCard("VI", "4111111111111111", "2023-01"),
	)

	// add traveler
	bookingReq.(*amadeus.BookingFlightOrderRequest).AddTraveler(
		bookingReq.(*amadeus.BookingFlightOrderRequest).
			NewTraveler(
				"Foo", "Bar", "MALE", "1990-02-15",
			).
			AddEmail("foo@bar.com").
			AddMobile("33", "480080072"),
	)

	// fmt.Println("--- send booking request start-----")
	// send request
	err = client.Do(bookingReq, &bookingResp, "POST")
	// fmt.Println("--- send booking request end -----")

	// get flight booking response
	bookingRespData := bookingResp.(*amadeus.BookingFlightOrderResponse)
	// fmt.Println("--- result -----")
	fmt.Println(bookingRespData.Data)
	// fmt.Println("--- result end -----")

}
