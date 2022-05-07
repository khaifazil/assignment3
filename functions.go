package main

import (
	"errors"
	"time"
)

var timeFormat = "02/01/2006"

func convertDate(date string) int {
	// make array of months in accumalated days not including current month's days
	// add date given by user.
	var (
		daysInMonths = [12]int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}
	)

	userDate, _ := time.Parse(timeFormat, date)
	month := int(userDate.Month()) - 1
	return daysInMonths[month] + userDate.Day()
}

func convertTime(time int) int {
	if time == 2400 {
		time = 0000
	}
	return time / 100
}

func selectSort(arr []*BookingInfoNode, n int) {
	for last := n - 1; last >= 1; last-- {
		largest := indexOfLargest(arr, last+1)
		swapPtr(&arr[largest], &arr[last])
	}
}

func indexOfLargest(arr []*BookingInfoNode, n int) int {
	largestIndex := 0
	for i := 1; i < n; i++ {
		checkAtIndex, _ := time.Parse(timeFormat, arr[i].Date)
		currentLargest, _ := time.Parse(timeFormat, arr[largestIndex].Date)
		if checkAtIndex.After(currentLargest) {
			largestIndex = i
		}
	}
	return largestIndex
}

func swapPtr(x **BookingInfoNode, y **BookingInfoNode) {
	temp := *x
	*x = *y
	*y = temp
}

func getTimeFromParse(d time.Time, _ error) time.Time {
	return d
}

func sortBookingsByDate(arr []*BookingInfoNode, n int) []*BookingInfoNode {
	for i := 1; i < n; i++ {
		data := arr[i]
		last := i
		dataDate, _ := time.Parse(timeFormat, data.Date)
		for (last > 0) && (getTimeFromParse(time.Parse(timeFormat, arr[last-1].Date)).After(dataDate)) {
			arr[last] = arr[last-1]
			last--
		}

		arr[last] = data
	}
	return arr
}

func sortBookingsByTime(arr []*BookingInfoNode, n int) []*BookingInfoNode {
	for i := 1; i < n; i++ {
		data := arr[i]
		last := i

		for (last > 0) && (arr[last-1].BookingTime > data.BookingTime) {
			arr[last] = arr[last-1]
			last--
		}

		arr[last] = data
	}
	return arr
}

// func updateCarArr(ptr *[365][24]*BookingInfoNode, index1 int, index2 int, address *BookingInfoNode) {
// 	ptr[index1][index2] = address
// }

func checkDate(date string) error {
	parsedDate, err := time.Parse(timeFormat, date)
	if err != nil {
		return err
	} else if parsedDate.Before(time.Now()) {
		return errors.New("date given has passed")
	}
	return nil
}

func binarySearchDate(arr []*BookingInfoNode, target string) int {
	first := 0
	last := len(arr) - 1

	for first <= last {
		mid := (first + last) / 2
		if arr[mid].Date == target {
			return mid
		} else {
			if target < arr[mid].Date {
				last = mid - 1
			} else {
				first = mid + 1
			}
		}
	}
	return -1
}

func lookForDupsDate(arr []*BookingInfoNode, n int, target string) []*BookingInfoNode {
	temp := []*BookingInfoNode{}
	for i := n; i >= 0; i-- {
		if target != arr[i].Date {
			break
		} else {
			temp = append(temp, arr[i])
		}
	}
	for i := n + 1; i < len(arr); i++ {
		if target != arr[i].Date {
			break
		} else {
			temp = append(temp, arr[i])
		}
	}
	return temp
}

func searchBookingByDate(arr []*BookingInfoNode, date string) ([]*BookingInfoNode, error) {
	n := binarySearchDate(arr, date)
	if n == -1 {
		return nil, errors.New("no bookings found at that date")
	}
	return lookForDupsDate(arr, n, date), nil
}

func add(x, y int) int {
	return x + y
}
