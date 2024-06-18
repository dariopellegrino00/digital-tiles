package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Piastrella struct {
	colore    string
	intensita int
}
type piano struct {
	piastrelle *map[[2]int]*Piastrella
}

/*
//TODO valurare se serve
// restituisce un puntatore alla piastrella in posizione x, y se accesa,
// nil altrimenti (o 0x0, 0x0 è il puntatore a nil)
func (p piano) piastrella(x int, y int) *Piastrella {
	return (*p.piastrelle)[pair(x, y)]
}*/

//(utility per inesistenza di un pair in go)
//restituisce un vettore di 2 posizioni (pair) contenente
//le coordinate x e y
func pair(x, y int) [2]int {
	return [2]int{x, y}
}

// colora una piastrella in posizione (x,y) sul piano p
// con colore alpha e intensità i
func colora(p piano, x int, y int, alpha string, i int) {
	(*p.piastrelle)[pair(x, y)] = &Piastrella{colore: alpha, intensita: i}
}

// spegne la piastrella se accesa, non fa niente altrimenti
func spegni(p piano, x int, y int) {
	delete(*p.piastrelle, pair(x, y))
}

// stampa in console lo stato di una piastrella nella forma "colore intensità"
// e restituisce i due valore sottoforma di stringa e intero
// non fa niente se la piastrella è spenta e restituisce una stringa vuota e zero
func stato(p piano, x int, y int) (string, int) {
	piastrella, err := (*p.piastrelle)[[2]int{x, y}]
	if err {
		fmt.Printf("%s %d\n", piastrella.colore, piastrella.intensita)
		return piastrella.colore, piastrella.intensita
	}
	return "", 0
}

func esegui(p piano, s string) {
	tokens := strings.Split(s, " ")
	comando := tokens[0]
	switch comando {
	case "C":
		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])
		alpha := tokens[3]
		i, _ := strconv.Atoi(tokens[4])
		colora(p, x, y, alpha, i)
	case "S":
		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])
		spegni(p, x, y)
	case "?":
		x, _ := strconv.Atoi(tokens[0])
		y, _ := strconv.Atoi(tokens[1])
		stato(p, x, y)
	}
}

func creaPiano() piano {
	p := make(map[[2]int]*Piastrella)
	return piano{&p}
}

func main() {
	p := creaPiano()
	esegui(p, "C 1 2 r 4")
	esegui(p, "C 0 2 g 4")
	esegui(p, "C 1 1 b 5")
	esegui(p, "? 0 2")
	esegui(p, "? 1 1")
}
