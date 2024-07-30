package main

import (
	"context"
	"log"
	"optimistic-booking-system/booking"
	"optimistic-booking-system/db"
	"time"
)

func main() {
	mongoCfg := db.New("mongodb://root:example@localhost:27017", "optimistic")
	mongo, err := mongoCfg.Connect()
	if err != nil {
		panic(err)
	}
	defer mongoCfg.Close(mongo)

	coll := mongo.Collection("bookings")

	// newBooking := booking.Booking{
	// 	UserID:     "user123",
	// 	ResourceID: "resource456",
	// 	StartTime:  time.Date(2024, 8, 1, 10, 0, 0, 0, time.Local),
	// 	EndTime:    time.Date(2024, 8, 1, 11, 0, 0, 0, time.Local),
	// }

	// _, err = booking.CreateBooking(context.Background(), coll, newBooking)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = booking.UpdateBookingWithOptimisticLock(context.Background(), coll, "66a8c60dc060aba34fb0b600", time.Date(2024, 8, 2, 11, 0, 0, 0, time.Local), time.Date(2024, 8, 1, 12, 0, 0, 0, time.Local), 1)
	if err != nil {
		log.Fatal(err)
	}
}
