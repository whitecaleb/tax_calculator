package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	amexPercent = 2.02 // percent you charge customer for amex.
	taxPercent  = 10   // percent of government tax.
	amexFee     = 1.8  // percent of transaction amex will take.
)

type input struct {
	string
}

func (i *input) set(reader *bufio.Reader) {
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	i.string = text[:len(text)-1]
}

func main() {
	var i = &input{}
	// create input reader
	reader := bufio.NewReader(os.Stdin)
	for i.string != "q" {
		printTitle()
		// Get cost
		fmt.Print("\nProduct cost (inc. GST): $")
		i.set(reader)
		if i.string == "q" {
			break
		}
		cost, err := strconv.ParseFloat(i.string, 64)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Find out if amex is incl.
		amex := false
		fmt.Print("Amex? (y/n): ")
		i.set(reader)
		if i.string == "q" {
			break
		}
		if i.string == "y" {
			amex = true
		}
		p := newProductAfterGST(cost, amex)
		fmt.Println("\n******************************")
		fmt.Printf("Product: $%f\n", p.cost)
		fmt.Printf("Amex charge: $%f\n", p.getAmexCharge())
		fmt.Println("******************************")
		sub, tax := p.getTax()
		fmt.Printf("Subtotal: $%f\n", sub)
		fmt.Println("******************************")
		fmt.Printf("GST: $%f\n", tax)
		fmt.Println("******************************")
		fmt.Printf("Total: $%f\n", p.getTotal())

		fmt.Printf("\nIf Amex take 1.8%% of total transaction and government take gst\n")
		fmt.Printf("Pocket: $%f\n\n\n", p.getPocket())
	}
	fmt.Println("Good bye")
}

func printTitle() {
	fmt.Println("\n******************************")
	fmt.Println("*** Josh's Cost Calculator ***")
	fmt.Println("******************************")
	fmt.Println("\nEnter q to exit")
}

type product struct {
	// cost without tax
	cost    float64
	incAmex bool
}

func newProductAfterGST(c float64, incAmex bool) product {
	return product{
		cost:    removeTax(c),
		incAmex: incAmex,
	}
}

func (p *product) getPocket() float64 {
	_, tax := p.getTax()
	return p.getTotal() - p.getAmexCharge() - tax
}

// getAmexCharge returns what we will be charged by amex upon transaction.
//
// c is the total cost with tax and charges already added.
func (p *product) getAmexCharge() float64 {
	if !p.incAmex {
		return 0
	}
	return p.getTotal() * (amexFee / 100)
}

func (p *product) getTotal() float64 {
	sub, tax := p.getTax()
	return sub + tax
}

// Workout tax charge
// return total with tax and just tax
func (p *product) getTax() (total float64, tax float64) {
	// If paying by amex, we need to add the charge to the pre-tax total, otherwise subtotal remains at cost.
	subtotal := p.cost
	if p.incAmex {
		subtotal = p.cost + percent(p.cost, amexPercent)
	}

	return subtotal, percent(subtotal, taxPercent)
}

// c is a non-amex product cost with tax already added.
func removeTax(c float64) float64 {
	var total float64 = 100 + taxPercent
	return (c / total) * 100
}

// Work out the given percent of the given value.
func percent(value, percent float64) float64 {
	return (value / 100) * percent
}
