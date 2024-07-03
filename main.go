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

// le possibili direzioni prendibili a partire da una piastrella (o per raggiungere i possibili vicini)
var dirs = [][2]int{
	{-1, 0}, {1, 0}, {0, -1}, {0, 1}, // destra, sinistra, sopra, giù
	{-1, -1}, {1, 1}, {-1, 1}, {1, -1}, // diagonali
}

func (r Regola) String() string {
	s := r.risultato + ": "
	for a, k := range r.condizione {
		s += fmt.Sprint(k) + a + " "
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

// (utility per questione di leggibilità)
// restituisce un puntatore alla piastrella in posizione x, y e false (no errore) se accesa,
// nil e true altrimenti (o 0x0, 0x0 è il puntatore a nil)
func (p piano) piastrella(x, y int) (*Piastrella, bool) {
	ps, e := p.piastrelle[punto(x, y)]
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
	p.piastrelle[punto(x, y)] = &Piastrella{colore: alpha, intensita: i}
}

// spegne la piastrella se accesa, non fa niente altrimenti
func spegni(p piano, x, y int) {
	delete(p.piastrelle, punto(x, y))
}

// stampa in console lo stato di una piastrella nella forma "colore intensità"
// e restituisce i due valore sottoforma di stringa e intero
// non fa niente se la piastrella è spenta e restituisce una stringa vuota e zero
func stato(p piano, x, y int) (string, int) {
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

//TODO creare struttura coda per chiarezza?? o se serve in altro successivamente
//funzione Helper per fare la BFS che veniva riutilizzata spesso diminuendo di molto il codice ripetuto,
//non cambia a livello di stime asintottiche rielaborare la visita
// tempo ammortizzato O(n) con n = #piastrelle del blocco
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

	// tot O(n^2) caso peggiore molto raro perchè deve sempre dover risolvere collissioni ad ogni accesso
	// quindi mediamente O(n)
	for len(coda) > 0 { //O(n) ripetizioni
		current := coda[0]
		coda = coda[1:]

		for _, dir := range dirs { // O(1) sono sempre 8 (possibili) vicini da controllare
			vicino := punto(current[0]+dir[0], current[1]+dir[1])
			ps, esiste := p.piastrelle[vicino]                                               //O(n) caso peggiore molto raro ammortizzato a O(1)
			if esiste && (!checkColore || ps.colore == pstart.colore) && !visitati[vicino] { //O(n) caso peggiore molto raro ammortizzato a O(1)
				visitati[vicino] = true
				coda = append(coda, vicino)
			}
		}
	}

	for v, _ := range visitati {
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
	return intensitaBlocco(p, bfsBlocco(p, x, y, false))
}

// restituisce l'intensità del blocco omogeneo a cui appartiene la pistrella
// di posizione x,y se accessa, 0 se spenta
func bloccoOmog(p piano, x, y int) int {
	return intensitaBlocco(p, bfsBlocco(p, x, y, true))
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
			// applica regola TODO funzione applica regola?
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

// TODO: valutare se è il modo più efficiente
// propaga le regole nel blocco di appartenenda di piatrella in posizione (x, y), tempo O(n*m)
func propagaBlocco(p piano, x, y int) {
	blocco := bfsBlocco(p, x, y, false) // O(n)
	if len(blocco) == 0 {
		return
	}

	supporto := creaPiano()
	supporto.regole = p.regole

	for _, pos := range blocco { // O(n)
		pias, _ := p.piastrelle[pos]                                        // non può non trovarle più se l'ha trovata sopra nella bfs
		supporto.piastrelle[pos] = &Piastrella{pias.colore, pias.intensita} // non devo passare il puntatore
	}

	for _, pos := range blocco { // O(n*m)
		//piano, x, y
		propaga(supporto, pos[0], pos[1]) // se m = len(r)
	}

	for pos, pias := range supporto.piastrelle { // O(n)
		p.piastrelle[pos] = pias
	}
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
	case "B":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		println(bloccoOmog(p, x, y))
	case "p":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		propaga(p, x, y)
	case "P":
		x, _ := strconv.Atoi(tokens[1])
		y, _ := strconv.Atoi(tokens[2])
		propagaBlocco(p, x, y)
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
	// TODO: rimuovi
	for _, r := range *p.regole {
		fmt.Println(*r, r.consumo)
	}
}
