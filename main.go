package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Piastrella struct {
	colore    string
	intensita int
}
type piano struct {
	piastrelle map[[2]int]*Piastrella
	regole     *[]*Regola
}

type Regola struct {
	condizione map[string]int
	risultato  string
	consumo    int
}

func (r Regola) String() string {
	s := r.risultato + ": "
	for a, k := range r.condizione {
		s += fmt.Sprint(k) + a + " "
	}
	return s
}

// (utility per questione di leggibilità)
// restituisce un puntatore alla piastrella in posizione x, y e false (no errore) se accesa,
// nil e true altrimenti (o 0x0, 0x0 è il puntatore a nil)
func (p piano) piastrella(x int, y int) (*Piastrella, bool) {
	ps, e := (p.piastrelle)[punto(x, y)]
	return ps, e
}

// (utility per inesistenza di un punto in go)
// restituisce un vettore di 2 posizioni (punto) contenente
// le coordinate x e y
func punto(x, y int) [2]int {
	return [2]int{x, y}
}

// colora una piastrella in posizione (x,y) sul piano p
// con colore alpha e intensità i
func colora(p piano, x int, y int, alpha string, i int) {
	(p.piastrelle)[punto(x, y)] = &Piastrella{colore: alpha, intensita: i}
}

// spegne la piastrella se accesa, non fa niente altrimenti
func spegni(p piano, x int, y int) {
	delete(p.piastrelle, punto(x, y))
}

// stampa in console lo stato di una piastrella nella forma "colore intensità"
// e restituisce i due valore sottoforma di stringa e intero
// non fa niente se la piastrella è spenta e restituisce una stringa vuota e zero
func stato(p piano, x int, y int) (string, int) {
	piastrella, err := p.piastrella(x, y)
	if err {
		fmt.Printf("%s %d\n", piastrella.colore, piastrella.intensita)
		return piastrella.colore, piastrella.intensita
	}
	return "", 0
}

// aggiunge una regola
func regola(p piano, regola string) {
	tokens := strings.Split(regola, " ")
	condizione := make(map[string]int)
	for i := 1; i < len(tokens)-1; i += 2 {
		// assumendo input corretto
		condizione[tokens[i+1]], _ = strconv.Atoi(tokens[i])
	}
	*p.regole = append(*(p.regole), &Regola{condizione: condizione, risultato: tokens[0]})
}

// stampa le regole di propagazione nell'ordine attuale
func stampa(p piano) {
	fmt.Println("(")
	for _, regola := range *p.regole {
		fmt.Println(*regola)
	}
	fmt.Println(")")
}

//TODO CONTINUTARE E TESTARE
//TODO creare struttura coda per chiarezza??
func blocco(p piano, x int, y int) int {
	start := punto(x, y)
	_, esiste := p.piastrelle[start]
	if !esiste {
		return 0
	}

	visitati := make(map[[2]int]bool)
	coda := [][2]int{start}
	visitati[start] = true
	sumIntensita := 0

	//TODO dovrebbero servire dopo estrarli e renderli utilizzabili globalmente
	//magari migliorando astrazione
	dirs := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // left, right, up, down
		{-1, -1}, {1, 1}, {-1, 1}, {1, -1}, // diagonals
	}

	// tot O(n^2) caso peggiore molto raro perchè deve sempre dover risolvere collissioni ad ogni accesso
	// quindi mediamente O(n)
	for len(coda) > 0 { //O(n) ripetizioni
		current := coda[0]
		coda = coda[1:]

		if ps, ok := p.piastrelle[current]; ok { // O(n) caso peggiore ricerca in hashtable
			sumIntensita += ps.intensita
		}

		for _, dir := range dirs { // O(1) sono sempre 8 (possibili) vicini da controllare
			vicino := punto(current[0]+dir[0], current[1]+dir[1])
			_, esiste = p.piastrelle[vicino] //O(n) caso peggiore
			if esiste && !visitati[vicino] { //O(n) caso peggiore
				visitati[vicino] = true
				coda = append(coda, vicino)
			}
		}
	}

	return sumIntensita
}

func esegui(p piano, s string) {
	tokens := strings.Split(s, " ")
	comando := tokens[0]
	switch comando {
	case "C":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		alpha := tokens[3]
		i, _ := strconv.Atoi(tokens[4])
		colora(p, x, y, alpha, i)
	case "S":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		spegni(p, x, y)
	case "?":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		stato(p, x, y)
	case "r":
		regola(p, s[2:])
	case "s":
		stampa(p)
	case "b":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		println(blocco(p, x, y))
	}
}

func creaPiano() piano {
	p := make(map[[2]int]*Piastrella)
	r := make([]*Regola, 0)
	return piano{p, &r}
}

func main() {
	p := creaPiano()
	file, err := os.Open("inputs/example1.txt")

	if err != nil {
		panic(err.Error())
	}

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		linea := sc.Text()
		esegui(p, linea)
	}
}
