package limoBookingApp

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
var booking *BookingInfoNode

func init() {
	rand.Seed(time.Now().UnixNano())
}

//MakeRandomBookingId makes a random string of variable length and returns the string
func MakeRandomBookingId(length int) string {

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bookingId := make([]byte, length)
	for i := range bookingId {
		bookingId[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bookingId)
}

//MakeNewBooking collects the user input and proceeds to make a new booking node in the linked list. Returns a pointer to the new booking node and any errors
func (b *LinkedList) MakeNewBooking(car string, date string, bookingTime int, userName string, pickUp string, dropOff string, contactInfo int, remarks string) (*BookingInfoNode, error) {
	t := ConvertTime(bookingTime)
	d := ConvertDate(date)
	carArr := GetCarArr(car)

	//fmt.Println(carArr[d][t])
	if carArr[d][t] != nil {
		return nil, errors.New("there is already a booking at that time and date")
	}

	bookingId := MakeRandomBookingId(6)

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

//AppendNodeToSlice appends a pointer to a bookingInfoNode to the slice parameter. Returns slice.
func AppendNodeToSlice(a []*BookingInfoNode, node *BookingInfoNode) []*BookingInfoNode {
	a = append(a, node)
	return a
}

//RecursiveSeqSearchId is a recursive function that searches in the slice parameter for the target string parameter. Returns pointer to BookingInfoNode, index of target and error if target is not found in slice.
func RecursiveSeqSearchId(length int, start int, arr []*BookingInfoNode, target string) (*BookingInfoNode, int, error) {
	if start > length-1 {
		return nil, 0, errors.New("there are no bookings with that ID")
	} else {

		if target == arr[start].BookingId {
			return arr[start], start, nil
		} else {
			return RecursiveSeqSearchId(length, start+1, arr, target)
		}
	}
}

//SearchId is the wrapper function for RecursiveSeqSearchId
func SearchId(arr []*BookingInfoNode, target string) (*BookingInfoNode, error) {
	booking, _, err := RecursiveSeqSearchId(len(arr), 0, arr, target)
	if err != nil {
		return nil, err
	}
	return booking, err
}

//DeleteBookingNode deletes the *BookingInfoNode in a linked list. Returns any errors.
func (b *LinkedList) DeleteBookingNode(ptr *BookingInfoNode) error {
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

//AppendAllToSlice appends all the *BookingInfoNodes in a linked list into a slice. Returns the slice and any errors.
func (b *LinkedList) AppendAllToSlice() ([]*BookingInfoNode, error) {
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
