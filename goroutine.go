package main

import (
	"flag"
	"log"
	"runtime"
	"time"
)

func dotp(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		log.Panic("Les vecteurs n'ont pas la même taille")
	}

	sum := 0.0
	for i := 0; i < len(v1); i++ {
		sum += v1[i] * v2[i]
	}

	return sum
}

func pdotp(v1, v2 []float64, c chan float64) {
	if len(v1) != len(v2) {
		log.Panic("Les vecteurs n'ont pas la même taille")
	}

	sum := 0.0
	for i := 0; i < len(v1); i++ {
		sum += v1[i] * v2[i]
	}

	c <- sum
}

func godotp(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		log.Panic("Les vecteurs n'ont pas la même taille")
	}

	sum := 0.0
	for i := 0; i < len(v1); i++ {
		sum += v1[i] * v2[i]
	}

	return sum
}

//func launch

func initv(n int) []float64 {
	v := make([]float64, n)

	for i, _ := range v {
		v[i] = float64(i) + 1.0
	}

	return v
}

func Help() {
	log.Println("Ce programme calcule le produit scalaire de 2 vecteurs de taille N, et le répartie sur p goroutines.")
	flag.PrintDefaults()
}

func main() {
	// Paramètre pour la parallélisation :
	var N, nbgo, nbproc int = 1000, 10, 4
	// Que faire tourner :
	var NoRout, NoSingle bool = false, false
	// Timer :
	var t0, t1 time.Time

	// Parsing de la ligne de commande :
	flag.Usage = Help
	flag.IntVar(&N, "N", 1000, "Taille des vecteurs.")
	flag.IntVar(&nbgo, "n", 10, "Nombre de GoRoutine à utiliser.")
	flag.IntVar(&nbproc, "p", 4, "Nombre de processeur à utiliser.")
	flag.BoolVar(&NoRout, "-no-goroutine", false, "Fait le calcul sans goroutine.")
	flag.BoolVar(&NoSingle, "-single", false, "Fait le calcul avec une goroutine simple (channel non-bufferisé).")
	flag.Parse()

	// On vérifie que l'utilisateur n'ait pas donné un nombre de CPU plus grand que celui disponible :
	nump := runtime.NumCPU()
	if nbproc > nump {
		log.Println("Set nbproc=", nbproc, " to nbproc=", nump, ".")
		nbproc = nump
	}

	// On dit à l'executable qu'il peut utiliser nbproc processeur :
	runtime.GOMAXPROCS(nbproc)

	// Création des variables nécessaire au calcul :
	var res float64 = 0.0         // Résultat de l'opération
	d := make(chan float64, nbgo) // Canal de communication

	// Création des 2 vecteurs :
	v1 := initv(N)
	v2 := initv(N)

	// Exécution de la méthode en linéaire :
	if NoRout {
		t0 = time.Now()
		log.Println("Résultat seul : ", dotp(v1, v2))
		t1 = time.Now()
		log.Println("Temps d'éxecution sur une routine simple : ", t1.Sub(t0))
	}

	// Exécution sur une seul Go routine :
	if NoSingle {
		c := make(chan float64)
		log.Println("Avec une goroutine :")
		t0 = time.Now()
		go pdotp(v1, v2, c)
		t1 = time.Now()

		log.Println("\tRésultat : ", <-c)
		log.Println("Temps d'éxecution sur une goroutine seule : ", t1.Sub(t0))
	}

	// Lancement des nbgo Go routine sur les nbproc processeurs :
	log.Println("Avec ", nbgo, " goroutines :")
	t0 = time.Now()
	//	-> lancement :
	for i := 0; i < nbgo; i++ {
		go pdotp(v1[i*N/nbgo:(i+1)*N/nbgo], v2[i*N/nbgo:(i+1)*N/nbgo], d)
	}
	//	-> récupération des résultats :
	for i := 0; i < nbgo; i++ {
		res += <-d
	}
	t1 = time.Now()
	log.Println("\tRésultat : ", res)
	log.Println("Temps d'éxecution avec ", nbgo, " goroutines : ", t1.Sub(t0))
}
