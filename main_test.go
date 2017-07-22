package main

import (
	"reflect"
	"testing"
)

func TestRemoveTax(t *testing.T) {
	f := removeTax(110)
	if f != float64(100) {
		t.Errorf("%s Failed. Expected 100, got %f", t.Name(), f)
	}
}

var getTaxTests = []struct {
	name string
	cost float64
	amex bool
	expd float64
}{
	{"WithoutAmex", 100, false, 10},
	{"WithAmex", 100, true, 10.202},
	{"DecimalWithoutAmex", 0.001, false, 0.0001},
	{"DecimalWithAmex", 0.001, true, 0.00010202},
	{"BigWithoutAmex", 1000000000, false, 100000000},
	{"BigWithAmex", 100000000, true, 10202000},
}

func TestGetTax(t *testing.T) {
	for _, tst := range getTaxTests {
		tst.name = t.Name() + tst.name

		p := product{cost: tst.cost, incAmex: tst.amex}

		total, tax := p.getTax()
		if tax != tst.expd {
			t.Errorf("%s Failed. Expected %f, got %f", tst.name, tst.expd, tax)
		}
		// TODO: better testing of total.
		if (!tst.amex && total != tst.cost) || (tst.amex && total <= tst.cost) {
			t.Errorf("%s Failed. Expected %f, got %f", tst.name, tst.cost, total)
		}
	}
}

var getTotalTests = []struct {
	name string
	cost float64
	amex bool
	expd float64
}{
	{"WithoutAmex", 100, false, 110},
	{"WithAmex", 100, true, 112.222},
	{"DecimalWithoutAmex", 0.1, false, 0.11},
	{"DecimalWithAmex", 0.001, true, 0.00112222},
	{"BigWithoutAmex", 1000000000, false, 1100000000},
	{"BigWithAmex", 100000, true, 112222},
}

func TestGetTotal(t *testing.T) {
	for _, tst := range getTotalTests {
		tst.name = t.Name() + tst.name

		p := product{cost: tst.cost, incAmex: tst.amex}

		total := p.getTotal()
		if total != tst.expd {
			t.Errorf("%s Failed. Expected %f, got %f", tst.name, tst.expd, total)
		}
	}
}

func aProduct(c float64, a bool) product {
	return product{cost: c, incAmex: a}
}

var newProductAfterGSTTests = []struct {
	name string
	cost float64
	amex bool
	expd product
}{
	{"WithoutAmex", 550, false, aProduct(500, false)},
	{"WithAmex", 550, true, aProduct(500, true)},
	{"DecimalWithoutAmex", 1.1, false, aProduct(1, false)},
	{"DecimalWithAmex", 0.0011, true, aProduct(0.001, true)},
}

func TestNewProductAfterGST(t *testing.T) {
	for _, tst := range newProductAfterGSTTests {
		tst.name = t.Name() + tst.name

		p := newProductAfterGST(tst.cost, tst.amex)

		if !reflect.DeepEqual(p, tst.expd) {
			t.Errorf("%s Failed. Expected %v, got %v", tst.name, tst.expd, p)
		}
	}
}
