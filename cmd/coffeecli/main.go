package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	brewpb "github.com/jany/my-coffee/gen/proto/brew"
	menupb "github.com/jany/my-coffee/gen/proto/menu"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// func main() {
// 	// bufio is used to read user input from the command line
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		showMenu()

// 		fmt.Println("Choose an option:")

// 		input, _ := reader.ReadString('\n')
// 		input = strings.TrimSpace(input) // Remove whitespace and newline characters

// 		if input == "3" {
// 			fmt.Println("Goodbye!!!")
// 			break
// 		}

// 		handleChoice(input)
// 		fmt.Println()
// 	}
// }

func main() {
	// Connect to Menu Service (port 50052)
	menuConn, err := grpc.NewClient(
		"localhost:50052", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		fmt.Printf("Cannot connect to menu service: %v\n", err)
		return
	}

	// Close connection when the program ends
	defer menuConn.Close()

	// Connect to Brew Service (port 50051)
	brewConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		fmt.Printf("Cannot connect to brew service: %v\n", err)
		return
	}

	defer brewConn.Close()

	// Create clients
	menuClient := menupb.NewMenuServiceClient(menuConn)
	brewClient := brewpb.NewBrewServiceClient(brewConn)

	runMenuLoop(menuClient, brewClient)
}

func runMenuLoop(menuClient menupb.MenuServiceClient, brewClient brewpb.BrewServiceClient) {
	reader := bufio.NewReader(os.Stdin)

	for {
		showMenu()

		fmt.Println("Choose an option:")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "4" {
			fmt.Println("Goodbye!")
			break
		}

		handleChoice(input, menuClient, brewClient)
		fmt.Println()
	}
}

func handleChoice(choice string, menuClient menupb.MenuServiceClient, brewClient brewpb.BrewServiceClient) {
	switch choice {
	case "1":
		viewMenu(menuClient)
	case "2":
		checkListOrders(brewClient)
	case "3":
		createOrder(brewClient)
	default:
		fmt.Println("Invalid option. Please choose 1, 2, 3, or 4.")
	}
}

func viewMenu(client menupb.MenuServiceClient) {
	ctx := context.Background()

	// Call RPC
	resp, err := client.GetMenu(ctx, &menupb.GetMenuRequest{})
	if err != nil {
		fmt.Printf("Get menu error %v\n", err)
		return
	}
	
	// Display items
	fmt.Println("Menu Items:")
	fmt.Println("===============================")

	// fmt.Printf("resp %v", resp)

	for index, item := range resp.Items {
		fmt.Printf("%d_ %s _ $%.2f\n", index+1, item.Name, item.Price )
		fmt.Println()
		fmt.Println("===============================")
	}
}

func checkListOrders(client brewpb.BrewServiceClient) {
	ctx := context.Background()

	resp, err := client.ListOrders(ctx, &brewpb.ListOrdersRequest{})

	if err != nil {
		fmt.Printf("List orders error %v\n", err)
		return
	}

	fmt.Printf("Orders:\n")
	fmt.Println("===============================")

	for index, order := range resp.Orders {
		fmt.Printf("%d_ %s _ %s\n", index+1, order.MenuItemName, order.Status)
		fmt.Println()
		fmt.Println("===============================")
	}
}

func createOrder(client brewpb.BrewServiceClient) {
	fmt.Printf("Enter Menu Item Name: ")

	ctx := context.Background()

	reader := bufio.NewReader(os.Stdin)

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	resp, err := client.OrderDrink(ctx, &brewpb.OrderRequest{MenuItemName: input})

	if err != nil {
		fmt.Printf("Order Drink error: %v\n", err)
		return
	}

	fmt.Printf("Order drink %v \n", resp)
}

func showMenu() {
	fmt.Println("Welcome to the Coffee CLI!")
	fmt.Println()
	fmt.Println("Menu:")
	fmt.Println("1. View menu")
	fmt.Println("2. Check List orders status")
	fmt.Println("3. Create an order")
	fmt.Println("4. Quit")
	fmt.Println()
}
