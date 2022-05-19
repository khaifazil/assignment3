package limoBookingApp

import (
	"time"
)

var timeFormat = "02/01/2006"

//ConvertDate converts the user date input into an int to be used to find a booking. Returns int.
func ConvertDate(date string) int {
	// make array of months in accumulated days not including current month's days
	// add date given by user.
	var (
		daysInMonths = [12]int{0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}
	)

	userDate, _ := time.Parse(timeFormat, date)
	month := int(userDate.Month()) - 1
	return daysInMonths[month] + userDate.Day()
}

//ConvertTime converts the user time input into an int to be used to find a booking. Returns int.
func ConvertTime(time int) int {
	if time == 2400 {
		time = 0000
	}
	return time / 100
}

//func selectSort(arr []*BookingInfoNode, n int) {
//	for last := n - 1; last >= 1; last-- {
//		largest := indexOfLargest(arr, last+1)
//		swapPtr(&arr[largest], &arr[last])
//	}
//}

//func indexOfLargest(arr []*BookingInfoNode, n int) int {
//	largestIndex := 0
//	for i := 1; i < n; i++ {
//		checkAtIndex, _ := time.Parse(timeFormat, arr[i].Date)
//		currentLargest, _ := time.Parse(timeFormat, arr[largestIndex].Date)
//		if checkAtIndex.After(currentLargest) {
//			largestIndex = i
//		}
//	}
//	return largestIndex
//}

//func swapPtr(x **BookingInfoNode, y **BookingInfoNode) {
//	temp := *x
//	*x = *y
//	*y = temp
//}

func getTimeFromParse(d time.Time, _ error) time.Time {
	return d
}

//SortBookingsByDate sorts a given slice of *BookingInfoNodes by ascending date.
func SortBookingsByDate(arr []*BookingInfoNode, n int) []*BookingInfoNode {
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

//SortBookingsByTime sorts a given slice of *BookingInfoNodes by ascending time.
func SortBookingsByTime(arr []*BookingInfoNode, n int) []*BookingInfoNode {
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
	var temp []*BookingInfoNode
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

//func searchBookingByDate(arr []*BookingInfoNode, date string) ([]*BookingInfoNode, error) {
//	n := binarySearchDate(arr, date)
//	if n == -1 {
//		return nil, errors.New("no bookings found at that date")
//	}
//	return lookForDupsDate(arr, n, date), nil
//}

//Add returns sum of the given int arguments
func Add(x, y int) int {
	return x + y
}
