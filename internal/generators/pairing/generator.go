package pairing

import (
	"path/filepath"

	"github.com/consensys/bavard"
)

// Data data used to generate the templates
type Data struct {
	Fpackage string
	// FpModulus string
	// FrModulus string
	// Fp2NonResidue string
	Fp6NonResidue   string
	EmbeddingDegree int
	T               string
	TNeg            bool

	// data needed in the template, always set to constants
	Fp2Name  string // TODO this name cannot change; remove it
	Fp6Name  string // TODO this name cannot change; remove it
	Fp12Name string // TODO this name cannot change; remove it

	// these members are computed as needed
	Frobenius [][]fp2Template // constants used Frobenius
}

// Generate generates pairing
func Generate(d Data, outputDir string) error {

	rootPath := filepath.Join(outputDir, d.Fpackage)

	d.InitFrobenius()

	// pairing.go
	{
		src := []string{
			Pairing,
			ExtraWork,
			MulAssign,
			Expt,
		}
		if err := bavard.Generate(filepath.Join(rootPath, "pairing.go"), src, d,
			bavard.Package(d.Fpackage),
			bavard.Apache2("ConsenSys AG", 2020),
			bavard.GeneratedBy("gurvy/internal/generators"),
		); err != nil {
			return err
		}
	}

	// frobenius.go
	{
		src := []string{
			Frobenius,
		}
		if err := bavard.Generate(filepath.Join(rootPath, "frobenius.go"), src, d,
			bavard.Package(d.Fpackage),
			bavard.Apache2("ConsenSys AG", 2020),
			bavard.GeneratedBy("gurvy/internal/generators"),
		); err != nil {
			return err
		}
	}

	return nil
}

const ConstantsTemplate = `
import (
	"math/big"
	"math/bits"
	"strings"

	"github.com/consensys/gurvy/{{$.Fpackage}}"
	"github.com/consensys/gurvy/{{$.Fpackage}}/fp"
)

type fp2 = {{$.Fpackage}}.{{$.Fp2Name}}

type fp2Template struct {
	fp2
	A0String, A1String string // base 10
}

// InitFrobenius set z.Frobenius to constants gamma[i][j]
// from https://eprint.iacr.org/2010/354.pdf (Section 3.2)
// gamma[i][j] = Fp6NonResidue^(j*(p^i-1)/d)
// where:
//   d = EmbeddingDegree / 2
//   i = range [1,2,3]
//   j = range [1,...,d-1]
func (z *Data) InitFrobenius() *Data {

	// compute d
	switch z.EmbeddingDegree {
	case 6, 12:
	default:
		panic("unsupported embedding degree")
	}
	d := (uint64)(z.EmbeddingDegree / 2)

	// allocate memory
	z.Frobenius = make([][]fp2Template, 3) // constants for Frobenius up to exponent 3
	for i := range z.Frobenius {
		z.Frobenius[i] = make([]fp2Template, d-1)
	}

	// parse Fp6NonResidue
	var nonResidue fp2Template
	{
		nonResidueStrings := strings.Split(z.Fp6NonResidue, ",")
		if len(nonResidueStrings) != 2 {
			panic("malformed Fp6NonResidue string")
		}
		nonResidue.A0.SetString(nonResidueStrings[0])
		nonResidue.A1.SetString(nonResidueStrings[1])
	}

	// compute exponent (p-1)/d
	var exponent big.Int
	exponent.Set(fp.ElementModulus())
	{
		var i big.Int
		i.SetUint64(1)
		exponent.Sub(&exponent, &i)
		i.SetUint64(d)
		exponent.Div(&exponent, &i)
	}

	// compute gamma = nonResidue^exponent
	// all other constants are derived from gamma
	var gamma fp2Template
	gamma.exp(&nonResidue, toUint64Slice(&exponent)...)

	// compute gamma[i][j] as in https://eprint.iacr.org/2010/354.pdf (Section 3.2)
	z.Frobenius[0][0].Set(&gamma.fp2)
	for j := 1; j < len(z.Frobenius[0]); j++ {
		z.Frobenius[0][j].Set(&z.Frobenius[0][j-1].fp2).
			MulAssign(&gamma.fp2)
	}
	for j := range z.Frobenius[1] {
		z.Frobenius[1][j].Conjugate(&z.Frobenius[0][j].fp2).
			MulAssign(&z.Frobenius[0][j].fp2)
	}
	for j := range z.Frobenius[2] {
		z.Frobenius[2][j].Set(&z.Frobenius[0][j].fp2).
			MulAssign(&z.Frobenius[1][j].fp2)
	}

	// compute helper data on gamma[i][j] for the template generator
	for i := range z.Frobenius {
		for j := range z.Frobenius[i] {
			f := &z.Frobenius[i][j]
			f.A0String = f.A0.String()
			f.A1String = f.A1.String()
		}
	}

	return z
}

// code copied from goff
func toUint64Slice(b *big.Int) (s []uint64) {
	s = make([]uint64, len(b.Bits()))
	for i, v := range b.Bits() {
		s[i] = (uint64)(v)
	}
	return
}

// exp sets z = x^exponent in fp2Template, return z
// exponent (non-montgomery form) is ordered from least significant word to most significant word
// code copied from goff
func (z *fp2Template) exp(x *fp2Template, exponent ...uint64) *fp2Template {
	r := 0
	msb := 0
	for i := len(exponent) - 1; i >= 0; i-- {
		if exponent[i] == 0 {
			r++
		} else {
			msb = (i * 64) + bits.Len64(exponent[i])
			break
		}
	}
	exponent = exponent[:len(exponent)-r]
	if len(exponent) == 0 {
		z.A0.SetOne()
		z.A1.SetZero()
		return z
	}

	z.Set(&x.fp2)

	l := msb - 2
	for i := l; i >= 0; i-- {
		z.Square(&z.fp2)
		if exponent[i/64]&(1<<uint(i%64)) != 0 {
			z.MulAssign(&x.fp2)
		}
	}
	return z
}
`