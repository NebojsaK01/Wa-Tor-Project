GitHub Link: https://github.com/NebojsaK01/Wa-Tor-Project

# Wa-Tor Simulation in Go

This project implements a **Wa-Tor predator-prey simulation** in Go.  
The simulation consists of **fish** and **sharks** on a toroidal grid. Each time step (chronon) updates the world 
according to the rules:

## For the fish
- At each chronon, a fish moves randomly to one of the adjacent unoccupied squares. If there are no free squares, 
  no movement takes place.

- Once a fish has survived a certain number of chronons it may reproduce. This is done as it moves to a neighbouring
  square, leaving behind a new fish in its old position. Its reproduction time is also reset to zero.

## For the sharks

- At each chronon, a shark moves randomly to an adjacent square occupied by a fish. If there is none, the shark moves 
  to a random adjacent unoccupied square. If there are no free squares, no movement takes place.
- At each chronon, each shark is deprived of a unit of energy.
- Upon reaching zero energy, a shark dies.
- If a shark moves to a square occupied by a fish, it eats the fish and earns a certain amount of energy.
- Once a shark has survived a certain number of chronons it may reproduce in exactly the same way as the fish.

## How to Run

1. Make sure you have **Go Version: go version go1.22.4 linux/arm64** installed.
2. Run the simulation:

## Run
go run main.go
