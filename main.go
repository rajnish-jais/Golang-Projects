//package main
//
//import (
//	"fmt"
//	"os"
//	"os/exec"
//)
//
//// Player represents a player in the game.
//type Player struct {
//	Name   string
//	Symbol string
//}
//
//// Board represents the game board.
//type Board struct {
//	Grid         [3][3]string
//	ActivePlayer *Player
//	Player1      *Player
//	Player2      *Player
//}
//
//// Game represents the Tic Tac Toe game.
//type Game struct {
//	Player1 *Player
//	Player2 *Player
//	Board   *Board
//}
//
//// NewGame creates a new instance of the Tic Tac Toe game.
//func NewGame(player1Name string, player2Name string) *Game {
//	player1 := &Player{Name: player1Name, Symbol: "X"}
//	player2 := &Player{Name: player2Name, Symbol: "O"}
//	board := &Board{Player1: player1, Player2: player2}
//	return &Game{Player1: player1, Player2: player2, Board: board}
//}
//
//// Start starts the Tic Tac Toe game.
//func (g *Game) Start() {
//	clearScreen()
//	fmt.Println("Let's play Tic Tac Toe!")
//
//	g.Board.ActivePlayer = g.Player1
//	g.Board.Print()
//
//	for !g.isGameOver() {
//		g.makeMove()
//		g.Board.SwitchActivePlayer()
//		g.Board.Print()
//	}
//
//	g.declareWinner()
//}
//
//// makeMove allows the active player to make a move.
//func (g *Game) makeMove() {
//	player := g.Board.ActivePlayer
//
//	var row, col int
//	for {
//		fmt.Printf("%s's turn. Enter the row (0-2): ", player.Name)
//		fmt.Scanln(&row)
//		fmt.Printf("%s's turn. Enter the column (0-2): ", player.Name)
//		fmt.Scanln(&col)
//
//		if g.Board.IsValidMove(row, col) {
//			break
//		}
//
//		fmt.Println("Invalid move. Please try again.")
//	}
//
//	g.Board.MakeMove(row, col, player.Symbol)
//}
//
//// isGameOver checks if the game is over.
//func (g *Game) isGameOver() bool {
//	if g.Board.HasWin(g.Player1.Symbol) {
//		return true
//	}
//
//	if g.Board.HasWin(g.Player2.Symbol) {
//		return true
//	}
//
//	if g.Board.IsFull() {
//		return true
//	}
//
//	return false
//}
//
//// declareWinner declares the winner of the game.
//func (g *Game) declareWinner() {
//	if g.Board.HasWin(g.Player1.Symbol) {
//		fmt.Printf("Congratulations! %s wins!\n", g.Player1.Name)
//	} else if g.Board.HasWin(g.Player2.Symbol) {
//		fmt.Printf("Congratulations! %s wins!\n", g.Player2.Name)
//	} else {
//		fmt.Println("It's a tie!")
//	}
//}
//
//// clearScreen clears the terminal screen.
//func clearScreen() {
//	cmd := exec.Command("clear")
//	cmd.Stdout = os.Stdout
//	cmd.Run()
//}
//
//// Print prints the current state of the game board.
//func (b *Board) Print() {
//	fmt.Println("---------")
//	for _, row := range b.Grid {
//		fmt.Printf("| %s | %s | %s |\n", row[0], row[1], row[2])
//		fmt.Println("---------")
//	}
//}
//
//// IsValidMove checks if the given row and column are valid for a move.
//func (b *Board) IsValidMove(row, col int) bool {
//	if row < 0 || row >= len(b.Grid) || col < 0 || col >= len(b.Grid[0]) {
//		return false
//	}
//
//	if b.Grid[row][col] != "" {
//		return false
//	}
//
//	return true
//}
//
//// MakeMove makes a move by placing the player's symbol at the given row and column.
//func (b *Board) MakeMove(row, col int, symbol string) {
//	b.Grid[row][col] = symbol
//}
//
//// SwitchActivePlayer switches the active player.
//func (b *Board) SwitchActivePlayer() {
//	if b.ActivePlayer == nil || b.ActivePlayer == b.Player1 {
//		b.ActivePlayer = b.Player2
//	} else {
//		b.ActivePlayer = b.Player1
//	}
//}
//
//// HasWin checks if the given symbol has won the game.
//func (b *Board) HasWin(symbol string) bool {
//	// Check rows
//	for _, row := range b.Grid {
//		if row[0] == symbol && row[1] == symbol && row[2] == symbol {
//			return true
//		}
//	}
//
//	// Check columns
//	for col := 0; col < len(b.Grid[0]); col++ {
//		if b.Grid[0][col] == symbol && b.Grid[1][col] == symbol && b.Grid[2][col] == symbol {
//			return true
//		}
//	}
//
//	// Check diagonals
//	if b.Grid[0][0] == symbol && b.Grid[1][1] == symbol && b.Grid[2][2] == symbol {
//		return true
//	}
//
//	if b.Grid[0][2] == symbol && b.Grid[1][1] == symbol && b.Grid[2][0] == symbol {
//		return true
//	}
//
//	return false
//}
//
//// IsFull checks if the game board is full.
//func (b *Board) IsFull() bool {
//	for _, row := range b.Grid {
//		for _, cell := range row {
//			if cell == "" {
//				return false
//			}
//		}
//	}
//
//	return true
//}
//
//// Entry point of the program
//func main() {
//	game := NewGame("Player 1", "Player 2")
//	game.Start()
//}

package main

import (
	"fmt"
	"sync"
)

// Event represents a generic event type
type Event struct {
	Type string
	Data interface{}
}

// EventBroker represents the event broker component
type EventBroker struct {
	subscribers map[string][]chan Event
	mu          sync.Mutex
}

// NewEventBroker creates a new EventBroker instance
func NewEventBroker() *EventBroker {
	return &EventBroker{
		subscribers: make(map[string][]chan Event),
	}
}

// Subscribe adds a new subscriber to the event broker
func (eb *EventBroker) Subscribe(eventType string) chan Event {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	ch := make(chan Event, 10) // Buffer channel to avoid blocking publishers
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)

	return ch
}

// Publish sends an event to the event broker
func (eb *EventBroker) Publish(event Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subscribers := eb.subscribers[event.Type]
	for _, ch := range subscribers {
		go func(ch chan Event) {
			ch <- event
		}(ch)
	}
}

// EventDispatcher represents the event dispatcher component
type EventDispatcher struct {
	eventBroker *EventBroker
}

// NewEventDispatcher creates a new EventDispatcher instance
func NewEventDispatcher(eventBroker *EventBroker) *EventDispatcher {
	return &EventDispatcher{
		eventBroker: eventBroker,
	}
}

// Start starts the event dispatcher to handle incoming events
func (ed *EventDispatcher) Start() {
	for eventType, subscribers := range ed.eventBroker.subscribers {
		go func(eventType string, subscribers []chan Event) {
			for {
				select {
				case event := <-subscribers[0]:
					// Process the event based on the subscriber's requirements or business logic
					fmt.Printf("Received event of type '%s': %+v\n", eventType, event)
				}
			}
		}(eventType, subscribers)
	}
}

func main() {
	eventBroker := NewEventBroker()
	eventDispatcher := NewEventDispatcher(eventBroker)

	// Start the event dispatcher in a separate goroutine
	go eventDispatcher.Start()

	// Create subscribers for different event types
	subscriber1 := eventBroker.Subscribe("type1")
	subscriber2 := eventBroker.Subscribe("type2")

	// Publish events
	eventBroker.Publish(Event{Type: "type1", Data: "Event 1"})
	eventBroker.Publish(Event{Type: "type2", Data: "Event 2"})
	eventBroker.Publish(Event{Type: "type1", Data: "Event 3"})

	// Wait for events to be processed
	fmt.Println(<-subscriber1)
	fmt.Println(<-subscriber2)
	fmt.Println(<-subscriber1)
}
