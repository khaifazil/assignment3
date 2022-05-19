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

	return newBookingInfoNode, nil
}

func AppendNodeToSlice(a []*BookingInfoNode, node *BookingInfoNode) []*BookingInfoNode {
	a = append(a, node)
	return a
}

func recursiveSeqSearchId(length int, start int, arr []*BookingInfoNode, target string) (*BookingInfoNode, int, error) {
	if start > length-1 {
		return nil, 0, errors.New("there are no bookings with that ID")
	} else {

		if target == arr[start].BookingId {
			return arr[start], start, nil
		} else {
			return recursiveSeqSearchId(length, start+1, arr, target)
		}
	}
}

func searchId(arr []*BookingInfoNode, target string) (*BookingInfoNode, error) {
	booking, _, err := recursiveSeqSearchId(len(arr), 0, arr, target)
	if err != nil {
		return nil, err
	}
	return booking, err
}

func (b *LinkedList) deleteBookingNode(ptr *BookingInfoNode) error {
	if b.Size == 0 {
		return errors.New("linked list is empty")
	}
	if b.Size == 1 {
		b.Head = nil
		b.Tail = nil
	}
	if b.Size == 2 {
		if b.Head == ptr {
			b.Head = b.Tail
			b.Tail.Prev = nil
			ptr.Next = nil
		}
		b.Tail = b.Head
		b.Head.Next = nil
		ptr.Prev = nil
	}
	if b.Size > 2 {
		if b.Head == ptr {
			b.Head = b.Head.Next
			b.Head.Prev = nil
			ptr.Next = nil
		} else if b.Tail == ptr {
			b.Tail = ptr.Prev
			b.Tail.Next = nil
			ptr.Prev = nil
		} else {
			ptr.Next.Prev = ptr.Prev
			ptr.Prev.Next = ptr.Next
			ptr.Next = nil
			ptr.Prev = nil
		}
	}
	b.Size--

	return nil
}

func (b *LinkedList) appendAllToSlice() ([]*BookingInfoNode, error) {
	if b.Head == nil {
		return nil, errors.New("there are no bookings")
	}

	currentNode := b.Head
	var temp = make([]*BookingInfoNode, 0, 5)
	for i := 1; i <= b.Size; i++ {
		temp = append(temp, currentNode)
		if currentNode.Next != nil {
			currentNode = currentNode.Next
		}
	}
	return temp, nil
}
