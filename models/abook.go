package models

type Abook struct {
	ID       int
	Title    string
	Duration int
	Author   string
	Narrator string
	Price    float64
}

type Bestsellers struct {
	Copies int
	Abook
}
