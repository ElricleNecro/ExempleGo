package main

import "fmt"

func main() {
	fmt.Println("Hello world !")

	sum   := 0
	nbite := 100
	for i:=0 ; i < nbite ; i++ {
		sum += i
	}
	fmt.Println("J'ai fait ", nbite, " itérations. Ce qui nous donne la somme : ", sum)


	fmt.Println("Voyons la correspondance caractére -- entier en unicode (Utf-8) :")
	for i:=0 ; i < 85 ; i++ {
		if i < 26 {
			fmt.Printf("%c : %d\t", i + 'a', i + 'a')
		} else if i >= 26 && i - 26 < 26 {
			fmt.Printf("%c : %d\t", i + 'A' - 26, i + 'A' - 26)
		} else {
			fmt.Printf("       \t")
		}
		fmt.Printf("\t%c : %d\t%c : %d\n", i + 'あ', i + 'あ', i + 'ア', i + 'ア')
	}
}

