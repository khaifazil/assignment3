package main

import (
	"errors"
	"math/rand"
	"time"
)

type BookingInfoNode struct {
	Car         string
	Date        string
	BookingTime int
	UserName    string
	PickUp      string
	DropOff     string
	ContactInfo int
	Remarks     string
	BookingId   string
	Prev        *BookingInfoNode
	Next        *BookingInfoNode
}

type LinkedList struct {
	Head *BookingInfoNode
	Tail *BookingInfoNode
	Size int
}

var bookings = &LinkedList{nil, nil, 0}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func makeRandomBookingId(length int) string {

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bookingId := make([]byte, length)
	for i := range bookingId {
		bookingId[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bookingId)
}

func (b *LinkedList) makeNewBooking(car string, date string, bookingTime int, userName string, pickUp string, dropOff string, contactInfo int, remarks string) (*BookingInfoNode, error) {
	t := convertTime(bookingTime)
	d := convertDate(date)
	carArr := getCarArr(car)

	//fmt.Println(carArr[d][t])
	if carArr[d][t] != nil {
		return nil, errors.New("there is already a booking at that time and date")
	}

	bookingId := makeRandomBookingId(6)

	newBookingInfoNode := &BookingInfoNode{
		Car:         car,
		Date:        date,
		BookingTime: bookingTime,
		UserName:    userName,
		PickUp:      pickUp,
		DropOff:     dropOff,
		ContactInfo: contactInfo,
		Remarks:     remarks,
		BookingId:   bookingId,
		Next:        nil,
		Prev:        nil,
	}
	if b.Head == nil {
		b.Head = newBookingInfoNode
		b.Tail = newBookingInfoNode
	} else {
		b.Tail.Next = newBookingInfoNode
		newBookingInfoNode.Prev = b.Tail
		b.Tail = newBookingInfoNode
	}
	b.Size++

	(*carArr)[d][t] = newBookingInfoNode

	myUser := mapUsers[userName]
	myUser.UserBookings = append(myUser.UserBookings, newBookingInfoNode)
	myUser.UserBookings = sortBookingsByTime(myUser.UserBookings, len(myUser.UserBookings))
	myUser.UserBookings = sortBookingsByDate(myUser.UserBookings, len(myUser.UserBookings))

	return newBookingInfoNode, nil
}
