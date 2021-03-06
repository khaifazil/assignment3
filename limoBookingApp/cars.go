package limoBookingApp

import (
	"errors"
)

var carsList = []string{"Car1", "Car2", "Car3", "Car4"}

var car1 [365][24]*BookingInfoNode
var car2 [365][24]*BookingInfoNode
var car3 [365][24]*BookingInfoNode
var car4 [365][24]*BookingInfoNode

func checkCarSelection(car string) error {
	for _, c := range carsList {
		if c == car {
			return nil
		}
	}
	return errors.New("car is not in selection")
}

//GetCarArr takes the user car selection input and returns the correct car array
func GetCarArr(car string) *[365][24]*BookingInfoNode {
	switch car {
	case "Car1":
		return &car1
	case "Car2":
		return &car2
	case "Car3":
		return &car3
	case "Car4":
		return &car4
	default:
		ErrorLogger.Println(errors.New("invalid car"))
		return nil
	}
}

//DeleteFromCarsArr deletes *BookingInfoNode from a car array
func DeleteFromCarsArr(ptr *BookingInfoNode) {
	t := ConvertTime(ptr.BookingTime)
	d := ConvertDate(ptr.Date)
	carArr := GetCarArr(ptr.Car)
	(*carArr)[d][t] = nil
}
