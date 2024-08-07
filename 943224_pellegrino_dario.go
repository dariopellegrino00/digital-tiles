package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	/*
		ordine serve per avere sempre lo stesso ordine di stampa
		quando si stampano le regole, per possibili test, la mappa non garantisce di ciclare sempre nello stesso ordine
		non serve un puntatore tanto non viene più modificato
	*/
	ordine    []string
	risultato string
	consumo   int
}

// le possibili direzioni prendibili a partire da una piastrella (o per raggiungere i possibili vicini)
var dirs = [][2]int{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // destra, sinistra, sopra, giù
	{-1, -1}, {1, 1}, {-1, 1}, {1, -1}, // diagonali
}

var spostamento = map[string][2]int{ // simil funzione da punto cardinale a simil versore
	"NN": {0, 1}, "NE": {1, 1}, "EE": {1, 0}, "SE": {1, -1},
	"SS": {0, -1}, "SO": {-1, -1}, "OO": {-1, 0}, "NO": {-1, 1},
}

func (p Piastrella) String() string {
	return fmt.Sprintf("%s %d\n", p.colore, p.intensita)
}

func (r Regola) String() string {
	s := r.risultato + ": "
	for _, c := range r.ordine {
		k, _ := r.condizione[c]
		s += fmt.Sprintf("%d %s ", k, c)
	}
	return s
}

// restituisce true se la regola è applicabile false altrimenti
func (r Regola) applicabile(colori map[string]int) bool {
	for colore, num := range r.condizione {
		if colori[colore] < num {
			return false
		}
	}
	return true
}

// utility
// restituisce un vettore di 2 posizioni (punto) contenente
// le coordinate x e y
func punto(x, y int) [2]int {
	return [2]int{x, y}
}

// colora una piastrella in posizione (x,y) sul piano p
// con colore alpha e intensità i
func colora(p piano, x int, y int, alpha string, i int) {
	if i == 0 {
		spegni(p, x, y)
	} else {
		p.piastrelle[punto(x, y)] = &Piastrella{colore: alpha, intensita: i}
	}
}

// spegne la piastrella se accesa, non fa niente altrimenti
func spegni(p piano, x int, y int) {
	delete(p.piastrelle, punto(x, y))
}

// stampa in console lo stato di una piastrella nella forma "colore intensità"
// e restituisce i due valore sottoforma di stringa e intero
// non fa niente se la piastrella è spenta e restituisce una stringa vuota e zero
func stato(p piano, x int, y int) (string, int) {
	piastrella, esiste := p.piastrelle[punto(x, y)]
	if esiste {
		fmt.Print(piastrella)
		return piastrella.colore, piastrella.intensita
	}
	return "", 0
}

// aggiunge una regola
func regola(p piano, r string) {
	tokens := strings.Split(r, " ")
	condizione := make(map[string]int)
	ordine := make([]string, 0)
	for i := 1; i < len(tokens)-1; i += 2 {
		// assumendo input corretto
		k, _ := strconv.Atoi(tokens[i])
		condizione[tokens[i+1]] = k
		ordine = append(ordine, tokens[i+1])
	}
	*p.regole = append(*(p.regole), &Regola{condizione: condizione, ordine: ordine, risultato: tokens[0]})
}

// stampa le regole di propagazione nell'ordine attuale
func stampa(p piano) {
	fmt.Println("(")
	for _, regola := range *p.regole {
		fmt.Println(*regola)
	}
	fmt.Println(")")
}

//funzione Helper per fare la BFS che veniva riutilizzata spesso diminuendo di molto il codice ripetuto,
//non cambia a livello di stime asintottiche rielaborare la visita
//tempo ammortizzato O(n) con n = #piastrelle del blocco
func bfsBlocco(p piano, x, y int, checkColore bool) [][2]int {
	start := punto(x, y)
	pstart, esiste := p.piastrelle[start]
	visita := make([][2]int, 0)
	if !esiste {
		return visita
	}

	visitati := make(map[[2]int]bool)
	coda := [][2]int{start}
	visitati[start] = true

	// n ripetizioni
	for len(coda) > 0 {
		current := coda[0]
		coda = coda[1:]

		for _, dir := range dirs { // O(1) sono sempre 8 (possibili) vicini da controllare e m = numero di elementi nella mappa
			vicino := punto(current[0]+dir[0], current[1]+dir[1])
			ps, esiste := p.piastrelle[vicino] //O(m) caso peggiore molto raro ammortizzato a O(1)
			if esiste && (!checkColore || ps.colore == pstart.colore) && !visitati[vicino] {
				visitati[vicino] = true
				coda = append(coda, vicino)
			}
		}
	}

	for v := range visitati { // n ripetizioni
		visita = append(visita, v)
	}

	return visita
}

// helper per calcoli dei blocchi riducendo ridondanza
func intensitaBlocco(p piano, blocco [][2]int) int {
	intensita := 0
	for _, pos := range blocco {
		intensita += p.piastrelle[pos].intensita
		//non controllo l'ok, se l'ha visitata prima la piastrelle deve per forza esserci
	}
	return intensita
}

// restituisce l'intensità il blocco a cui appartiene la pistrella di posizione x,y
// di posizione x,y se accessa, 0 se spenta
func blocco(p piano, x, y int) int {
	intensita := intensitaBlocco(p, bfsBlocco(p, x, y, false))
	fmt.Println(intensita)
	return intensita
}

// restituisce l'intensità del blocco omogeneo a cui appartiene la pistrella
// di posizione x,y se accessa, 0 se spenta
func bloccoOmog(p piano, x, y int) int {
	intensita := intensitaBlocco(p, bfsBlocco(p, x, y, true))
	fmt.Println(intensita)
	return intensita
}

// propaga la prima formula compatibile nella piastrella x y, tempo O(len(r))
func propaga(p piano, x, y int) {
	for _, r := range *p.regole { // O(len(r))
		coloriVicini := make(map[string]int)

		for _, dir := range dirs {
			vicina, ok := p.piastrelle[punto(x+dir[0], y+dir[1])]
			if ok {
				coloriVicini[vicina.colore] += 1
			}
		}

		if r.applicabile(coloriVicini) {
			ps, accesa := p.piastrelle[punto(x, y)]
			if !accesa {
				colora(p, x, y, r.risultato, 1)
			} else {
				ps.colore = r.risultato
			}
			r.consumo++
			return
		}
	}

}

// propaga le regole nel blocco di appartenenda di piatrella in posizione (x, y)
func propagaBlocco(p piano, x, y int) {
	blocco := bfsBlocco(p, x, y, false)
	if len(blocco) == 0 {
		return
	}

	supporto := make(map[[2]int]string)

	for _, pos := range blocco {
		pias, _ := p.piastrelle[pos]
		supporto[pos] = pias.colore
	}

	for _, pos := range blocco {
		for _, r := range *p.regole {
			coloriVicini := make(map[string]int)

			for _, dir := range dirs {
				vicina, ok := p.piastrelle[punto(pos[0]+dir[0], pos[1]+dir[1])]
				if ok {
					coloriVicini[vicina.colore] += 1
				}
			}

			if r.applicabile(coloriVicini) {
				supporto[pos] = r.risultato
			}
		}
	}

	for pos, pias := range supporto {
		p.piastrelle[pos].colore = pias
	}
}

/*
Funzione di lib:
utilizza un algoritmo mergesort modificato per funzionare in loco, nel complesso, il suo caso peggiore
fa O(n log n) confronti. il fatto che utlizza mergesort lo rende stabile ed era quello che ci serviva.
*/
func ordina(p piano) {
	sort.SliceStable(*p.regole, func(i, j int) bool { return (*p.regole)[i].consumo < (*p.regole)[j].consumo })
}

func pista(p piano, x, y int, s string) {
	seq := strings.Split(s, ",")
	_, ok := p.piastrelle[punto(x, y)]
	if !ok {
		return
	}

	pista := [][2]int{punto(x, y)}
	next := punto(x, y)
	for _, d := range seq {
		pos, _ := spostamento[d] // assumendo input corretto
		next[0] = next[0] + pos[0]
		next[1] = next[1] + pos[1]
		_, ok = p.piastrelle[next]
		if !ok {
			return
		}
		pista = append(pista, next)
	}
	fmt.Println("[")
	for _, v := range pista {
		fmt.Print(v[0], v[1], p.piastrelle[v])
	}
	fmt.Println("]")
}

// utilizza bfs perchè bfs con pesi tutti uguali sugli archi (tutti 1)
// a partire dal primo vertice s trova sempre il percorso minimo
// verso gli archi v che esplora successivamente
// restituisce la pista di lunghezza minima tra i due (numero archi)
// la lunghezza della pista è data da il numero di piastrelle che la compongono
// una pista da x y a x y è lunga 0 chiaramente
func lung(p piano, x1 int, y1 int, x2 int, y2 int) int {
	start := punto(x1, y1)
	goal := punto(x2, y2)

	_, esiste1 := p.piastrelle[start]
	_, esiste2 := p.piastrelle[goal]
	if !esiste1 || !esiste2 {
		return -1
	}

	visitati := make(map[[2]int]bool)
	coda := [][3]int{{x1, y1, 0}}
	visitati[start] = true

	for len(coda) > 0 {
		current := coda[0]
		coda = coda[1:]

		if punto(current[0], current[1]) == goal {
			fmt.Println(current[2])
			return current[2]
		}

		for _, dir := range dirs {
			vicino := punto(current[0]+dir[0], current[1]+dir[1])
			_, esiste := p.piastrelle[vicino]
			if esiste && !visitati[vicino] {
				visitati[vicino] = true
				coda = append(coda, [3]int{vicino[0], vicino[1], current[2] + 1})
			}
		}
	}

	return -1
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
		blocco(p, x, y)
	case "B":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		bloccoOmog(p, x, y)
	case "p":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		propaga(p, x, y)
	case "P":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		propagaBlocco(p, x, y)
	case "o":
		ordina(p)
	case "t":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		pista(p, x, y, tokens[3])
	case "L":
		x1, _ := strconv.Atoi(tokens[1])
		y1, _ := strconv.Atoi(tokens[2])
		x2, _ := strconv.Atoi(tokens[3])
		y2, _ := strconv.Atoi(tokens[4])
		lung(p, x1, y1, x2, y2)
	case "q":
		os.Exit(0)
	}
}

func creaPiano() piano {
	p := make(map[[2]int]*Piastrella)
	r := make([]*Regola, 0)
	return piano{p, &r}
}

func eseguiFile() {

	file, err := os.Open(os.Args[1:][0])
	if err != nil {
		panic(err.Error())
	}

	p := creaPiano()
	sc := bufio.NewScanner(file)

	for sc.Scan() {
		linea := sc.Text()
		esegui(p, linea)
	}

}

func main() {
	// TODO: spiegazione in relazione
	if len(os.Args) > 1 {
		eseguiFile()
	} else {
		p := creaPiano()
		sc := bufio.NewScanner(os.Stdin)

		for sc.Scan() {
			linea := sc.Text()
			esegui(p, linea)
		}
	}
}
