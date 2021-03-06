package twistededwards

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// CurveParams curve parameters: ax^2 + y^2 = 1 + d*x^2*y^2
type CurveParams struct {
	A, D     fr.Element // in Montgomery form
	Cofactor fr.Element // not in Montgomery form
	Order    big.Int
	Base     PointAffine
}

var edwards CurveParams

// GetEdwardsCurve returns the twisted Edwards curve on BN254's Fr
func GetEdwardsCurve() CurveParams {
	// copy to keep Order private
	var res CurveParams

	res.A.Set(&edwards.A)
	res.D.Set(&edwards.D)
	res.Cofactor.Set(&edwards.Cofactor)
	res.Order.Set(&edwards.Order)
	res.Base.Set(&edwards.Base)

	return res
}

func init() {

	edwards.A.SetUint64(168700)
	edwards.D.SetUint64(168696)
	edwards.Cofactor.SetUint64(8).FromMont()
	edwards.Order.SetString("2736030358979909402780800718157159386076813972158567259200215660948447373041", 10)

	edwards.Base.X.SetString("5299619240641551281634865583518297030282874472190772894086521144482721001553")
	edwards.Base.Y.SetString("16950150798460657717958625567821834550301663161624707787222815936182638968203")
}
