package main

import (
    "fmt"
    "math/rand"
    "time"
)

// automatically called on start
func init() {
    // new random seed
    rand.Seed(time.Now().UnixNano())
}

func main() {
    // guess variable 
    var guess int

    // number of guesses variable
    var count int

    // pick a no between 1-100, [0,99] + 1
    // rand.Intn(100) returns a number between 0-99
    num := rand.Intn(100) + 1

    fmt.Println("I'm thinking of a number between 1-100 ")

    // loop until user guesses the number
    for {
        fmt.Print("Guess: ")
        _, err := fmt.Scanf("%d", &guess)
        if err == nil {
            count += 1 // increment guess counter
            if guess > num {
                fmt.Println(" Too high ")
            } else if guess < num {
                fmt.Println(" Too low ")
            } else {
                fmt.Printf("Correct! It took you %d guesses!\n", count)
                break
            }
        } else { // an error with input
            fmt.Println("Please input a number")
        }
    }
}

/*
Guess the Number Game

Problem:
The program randomly selects a number between 1 and 100.
The user must guess the number in a limited number of tries (variable), receiving feedback:
- "Too low!" if the guess is below the number
- "Too high!" if the guess is above the number
- "Correct!" if the guess is right

If the user has lost, return "You have lost!" every time a guess is tried
If the user won, return "You have won!" every time a guess is tried

Input: integers from user input (guesses)
Output: feedback strings

Implement the Guess function that would pass the tests, as well as a
command line user interface that would generate a game (with user given
max retries), generate a secret random number, then let the user play the game.
*/