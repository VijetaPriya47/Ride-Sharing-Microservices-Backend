package main

import (
	"fmt"
	"strings"
)	

func main() {

	conferenceName := "Go Conference"
	const conferenceTickets int = 50
	var remainingTickets uint = 50

	fmt.Printf("conferenceTickets is %T, remainingTickets is %T, conferenceName is %T\n",conferenceTickets,remainingTickets,conferenceName)

	fmt.Printf("Welcome to %v booking application.\n", conferenceName)
	fmt.Printf("We have total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend")

	var bookings []string
	
	for {
		var firstname string
		var lastname string
		var userTickets uint
		var email string
		// ask user for their name

		// fmt.Print(userName) // Prints the value if the variable
		// fmt.Print(&userName) // Prints the address of the memory (Printf   x)  (Print, Println     v/)

		fmt.Println("Enter your first name: ")
		fmt.Scan(&firstname)

		fmt.Println("Enter your last name: ")
		fmt.Scan(&lastname)

		fmt.Println("Enter your Tickets: ")
		fmt.Scan(&userTickets)

		fmt.Println("Enter your Email: ")
		fmt.Scan(&email)


		isValidName := len(firstName) >= 2 && len(lastName) >= 2
		isValidEmail := strings.Contains(email, "@")
		isValidTicketNumber := userTickets > 0 && userTickets <= remainingTickets

		// isValidCity := city = "Singapore" || city = "London"
		// !isValidCity
		

		// bookings[0]  = firstname + " " + lastname   //  ->array
		bookings = append(bookings, firstname + " " + lastname)



		if isValidEmail && isValidName && isValidTicketNumber {


			remainingTickets = remainingTickets - userTickets

			fmt.Printf("The whiole array: %v\n", bookings)
			fmt.Printf("The first value: %v\n", bookings[0])
			fmt.Printf("Array type: %T\n", bookings)
			fmt.Printf("Array length: %v\n", len(bookings))

			fmt.Printf("User %v booked %v tickets and the remaining tickets are %v.\n", bookings[0], userTickets, remainingTickets) 


			firstNames := []string{}
			for _, booking := range bookings {
				var names = strings.Fields(booking)
				firstNames = append(firstNames, names[0])
			}
			fmt.Printf("The first names of the bookings are: %v\n", firstNames)

			// noTicketsRemaining := remainingTickets == 0 //bool
			if remainingTickets == 0 {
				// end program
				fmt.Println("Our conference is booked out. Come back next year.")
				break
			}
		} else {
			fmt.Printf("We have only %v tickets remaining, so you can't book %v tickets\n", remainingTickets, userTickets)
			fmt.Print("Input is invalid, Please try again.")

		}
	}
}