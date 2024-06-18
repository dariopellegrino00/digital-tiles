## Operazioni

• **colora** *(x, y, α, i)* ✅
Colora Piastrella(x, y) di colore α e intensità i, qualunque sia lo stato di Piastrella(x, y) prima dell’operazione. 

• **spegni** *(x, y)* ✅
Spegne Piastrella(x, y). Se Piastrella(x, y) è già spenta, non fa nulla.

• **regola** *(k1, α1, k2, α2, . . . , kn, αn, β)* ✅
Definisce la regola di propagazione k1α1 + k2α2 + · · · + kn αn → β e la inserisce in fondo all’elenco delle regole.

• **stato** *(x, y)* ✅
Stampa e restituisce il colore e l’intensità di Piastrella(x, y). Se Piastrella(x, y) è spenta, non stampa nulla.

• **stampa** ✅
Stampa l’elenco delle regole di propagazione, nell’ordine attuale.

• **blocco** *(x, y)*
Calcola le stampa a somma delle intensità delle piastrelle contenute nel blocco di appartenenza di Piastrella(x, y). Se Piastrella(x, y) è spenta, restituisce 0.

• **bloccoOmog** *(x, y)*
Calcola e stampa la somma delle intensità delle piastrelle contenute nel blocco omogeneo di appartenenza di Piastrella(x, y). Se Piastrella(x, y) `e spenta, restituisce 0.

• **propaga** *(x, y)*
Applica a Piastrella(x, y) la prima regola di propagazione applicabile dell’elenco, ricolorando la piastrella. Se nessuna regola è applicabile, non viene eseguita alcuna operazione

• **propagaBlocco** *(x, y)*
Propaga il colore sul blocco di appartenenza di Piastrella(x, y).

• **ordina**
Ordina l’elenco delle regole di propagazione in base al consumo delle regole stesse: la regola con consumo maggiore diventa l’ultima dell’elenco. Se due regole hanno consumo uguale mantengono il loro ordine relativo.

• **pista** *(x, y, s)*
Stampa la pista che parte da Piastrella(x, y) e segue la sequenza di direzioni s, se tale pista è definita. Altrimenti non stampa nulla.

• **lung** *(x1, y1, x2, y2)*
Determina la lunghezza della pista più breve che parte da Piastrella(x1, y1) e arriva in Piastrella(x2, y2). Altrimenti non stampa nulla.

• **intensità** *(x1, y1, x2, y2)*
Determina l’intensità minima tra le intensità di tutte le piste che partono da Piastrella(x1, y1) e arrivano in Piastrella(x2, y2). Se non vi è alcuna pista tra queste due piastrelle, non stampa nulla.

• **perimetro** *(x, y)*
Calcola la lunghezza del perimetro del blocco di appartenenza della piastrella in (x, y).