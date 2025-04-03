package main

type Game struct {
	Secret, Tries, MaxTries int
}

func NewGame(secret, maxTries int) *Game {
	return &Game{Secret: secret, MaxTries: maxTries}
}

func (g *Game) Guess(n int) string {
	return "You have won!"
}
