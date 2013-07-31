package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"
	//"unsafe"
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
	fmt.Println("Ce programme calcule le produit scalaire de 2 vecteurs de taille N, et le répartie sur p goroutines.")
	flag.PrintDefaults()
}

func main() {
	var N, nbgo, nbproc int = 1000, 10, 4
	var NoRout, NoSingle bool = false, false
	var t0, t1 time.Time

	flag.Usage = Help
	flag.IntVar(&N, "N", 1000, "Taille des vecteurs.")
	flag.IntVar(&nbgo, "n", 10, "Nombre de GoRoutine à utiliser.")
	flag.IntVar(&nbproc, "p", 4, "Nombre de processeur à utiliser.")
	flag.BoolVar(&NoRout, "-no-goroutine", false, "Fait le calcul sans goroutine.")
	flag.BoolVar(&NoSingle, "-single", false, "Fait le calcul avec une goroutine simple (channel non-bufferisé).")
	flag.Parse()

	nump := runtime.NumCPU()
	if nbproc > nump {
		fmt.Println("Set nbproc=", nbproc, " to nbproc=", nump, ".")
		nbproc = nump
	}

	var res float64 = 0.0
	d := make(chan float64, nbgo)

	v1 := initv(N)
	v2 := initv(N)

	if NoRout {
		t0 = time.Now()
		fmt.Println("Résultat seul : ", dotp(v1, v2))
		t1 = time.Now()
		fmt.Println("Temps d'éxecution sur une routine simple : ", t1.Sub(t0))
	}

	if NoSingle {
		c := make(chan float64)
		fmt.Println("Avec une goroutine :")
		t0 = time.Now()
		go pdotp(v1, v2, c)
		t1 = time.Now()

		fmt.Println("\tRésultat : ", <-c)
		fmt.Println("Temps d'éxecution sur une goroutine seule : ", t1.Sub(t0))
	}

	fmt.Println("Avec ", nbgo, " goroutines :")
	t0 = time.Now()
	for i := 0; i < nbgo; i++ {
		//fmt.Println("Lancement de la goroutine n°=", i+1, " [", i*N/nbgo,":", (i+1)*N/nbgo, "]")

		//fmt.Println(i, N, nbgo, N/nbgo, i+1, (i+1)*N/nbgo, len(v1), len(v2), int(^uint(0) >> 1) )
		//fmt.Printf("%T %d %T (%d)\n", i, i, (i+1)*N/nbgo, unsafe.Sizeof(i))
		//fmt.Printf("%T %T %T (%d)\n", i, N, nbgo, unsafe.Sizeof(i))

		go pdotp(v1[i*N/nbgo:(i+1)*N/nbgo], v2[i*N/nbgo:(i+1)*N/nbgo], d)
	}

	for i := 0; i < nbgo; i++ {
		res += <-d
	}
	t1 = time.Now()
	fmt.Println("\tRésultat : ", res)
	fmt.Println("Temps d'éxecution avec ", nbgo, " goroutines : ", t1.Sub(t0))
}
