package main

import "fmt"

type Car struct {
	// Member Variables
	rpm       int
	noOfGears int
	noOfSeats int
}

//Member Functions

func (car Car) SetRPM(rpm int) {
	car.rpm = rpm
}

func (car Car) SetNoOfGears(noOfGears int) {
	car.noOfGears = noOfGears
}

func (car Car) SetNoOfSeats(noOfSeats int) {
	car.noOfSeats = noOfSeats
}

func (car Car) GetRPM() int {
	return car.rpm
}

func (car Car) GetNoOfGears() int {
	return car.noOfGears
}

func (car Car) GetNoOfSeats() int {
	return car.noOfSeats
}

func main() {
	var wagonR Car //Object of Car struct.

	wagonR.SetRPM(1100)    // accessing SetRPM function using wagonR object.
	wagonR.SetNoOfGears(6) // accessing SetNoOfGears function using wagonR object.
	wagonR.SetNoOfSeats(5) // accessing SetNoOfSeats function using wagonR object.

	fmt.Println("wagonR RPM:", wagonR.GetRPM())               // accessing GetRPM function using wagonR object.
	fmt.Println("wagonR No of Gears:", wagonR.GetNoOfGears()) // accessing GetNoOfGears function using wagonR object.
	fmt.Println("wagonR No of Seats:", wagonR.GetNoOfSeats()) // accessing GetNoOfSeats function using wagonR object.
}
