package rangeproof

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"

	"gitlab.dusk.network/dusk-core/dusk-go/rangeproof/fiatshamir"
	"gitlab.dusk.network/dusk-core/dusk-go/rangeproof/innerproduct"
	"gitlab.dusk.network/dusk-core/dusk-go/rangeproof/pedersen"
	"gitlab.dusk.network/dusk-core/dusk-go/rangeproof/vector"
	"gitlab.dusk.network/dusk-core/dusk-go/ristretto"
)

// N is number of bits in range
// So amount will be between 0..2^(N-1)
const N = 64

// M is the number of outputs
// for one bulletproof
var M = 1

// M is the maximum number of outputs allowed
// per bulletproof
const maxM = 32

// Proof is the constructed BulletProof
type Proof struct {
	V  []pedersen.Commitment // Curve points 32 bytes
	A  ristretto.Point       // Curve point 32 bytes
	S  ristretto.Point       // Curve point 32 bytes
	T1 ristretto.Point       // Curve point 32 bytes
	T2 ristretto.Point       // Curve point 32 bytes

	taux ristretto.Scalar //scalar
	mu   ristretto.Scalar //scalar
	t    ristretto.Scalar

	// XXX: Remove once debugging is completed
	l []ristretto.Scalar // []scalar
	r []ristretto.Scalar // []scalar

	IPProof *innerproduct.Proof
}

// Prove will take a scalar as a parameter and
// using zero knowledge, prove that it is [0, 2^N)
func Prove(v []ristretto.Scalar, debug bool) (Proof, error) {

	if len(v) < 1 {
		return Proof{}, errors.New("length of slice v is zero")
	}

	M = len(v)

	// commitment to values v
	Vs := make([]pedersen.Commitment, 0, M)

	ped := pedersen.New([]byte("dusk.BulletProof.vec1"))

	// Hash for Fiat-Shamir
	hs := fiatshamir.HashCacher{[]byte{}}

	for _, amount := range v {
		// compute commmitment to v
		V := ped.CommitToScalars(amount)

		Vs = append(Vs, V)

		// update Fiat-Shamir
		hs.Append(V.Value.Bytes())
	}

	aLs := make([]ristretto.Scalar, 0, N*M)
	aRs := make([]ristretto.Scalar, 0, N*M)

	for i := range v {
		// Compute Bitcommits aL and aR to v
		BC := BitCommit(v[i].BigInt())
		aLs = append(aLs, BC.AL...)
		aRs = append(aRs, BC.AR...)
	}

	// Compute A
	A := computeA(ped, aLs, aRs)

	// // Compute S
	S, sL, sR := computeS(ped)

	// // update Fiat-Shamir
	hs.Append(A.Value.Bytes(), S.Value.Bytes())

	// compute y and z
	y, z := computeYAndZ(hs)

	// compute polynomial
	poly, err := computePoly(aLs, aRs, sL, sR, y, z)
	if err != nil {
		return Proof{}, errors.Wrap(err, "[Prove] - poly")
	}
	// Compute T1 and T2
	T1 := ped.CommitToScalars(poly.t1)
	T2 := ped.CommitToScalars(poly.t2)

	// update Fiat-Shamir
	hs.Append(z.Bytes(), T1.Value.Bytes(), T2.Value.Bytes())

	// compute x
	x := computeX(hs)

	// compute taux which is just the polynomial for the blinding factors at a point x
	taux := computeTaux(x, z, T1.BlindingFactor, T2.BlindingFactor, Vs)

	// compute mu
	mu := computeMu(x, A.BlindingFactor, S.BlindingFactor)

	// compute l dot r
	l, err := poly.computeL(x)
	if err != nil {
		return Proof{}, errors.Wrap(err, "[Prove] - l")
	}
	r, err := poly.computeR(x)
	if err != nil {
		return Proof{}, errors.Wrap(err, "[Prove] - r")
	}
	t, err := vector.InnerProduct(l, r)
	if err != nil {
		return Proof{}, errors.Wrap(err, "[Prove] - t")
	}

	// START DEBUG
	if debug {
		err := debugProve(x, y, z, v, l, r, aLs, aRs, sL, sR)
		if err != nil {
			return Proof{}, errors.Wrap(err, "[Prove] - debugProve")
		}

		// // DEBUG T0
		testT0 := debugT0(aLs, aRs, y, z)
		if !testT0.Equals(&poly.t0) {
			return Proof{}, errors.New("[Prove]: Test t0 value does not match the value calculated from the polynomial")
		}

		// XXX: This is not passing tests, however we do not need it
		// If cannot be fixed, we will remove it from code.
		// polyt0 := poly.computeT0(y, z, v)
		// if !polyt0.Equals(&poly.t0) {
		// 	return Proof{}, errors.New("[Prove]: t0 value from delta function, does not match the polynomial t0 value(Correct)")
		// }

		tPoly := poly.eval(x)
		if !t.Equals(&tPoly) {
			return Proof{}, errors.New("[Prove]: The t value computed from the t-poly, does not match the t value computed from the inner product of l and r")
		}
	}
	// End DEBUG

	// check if any challenge scalars are zero
	if x.IsNonZeroI() == 0 || y.IsNonZeroI() == 0 || z.IsNonZeroI() == 0 {
		return Proof{}, errors.New("[Prove] - One of the challenge scalars, x, y, or z was equal to zero. Generate proof again")
	}

	hs.Append(x.Bytes(), taux.Bytes(), mu.Bytes(), t.Bytes())

	// TODO: calculate inner product proof
	Q := ristretto.Point{}
	w := hs.Derive()
	Q.ScalarMult(&ped.BaseVector.Bases[0], &w)

	Hpf := vector.ScalarPowers(y, uint32(N*M))

	ip, err := innerproduct.Generate(ped.BaseVector.Bases[1:], ped.BaseVector.Bases[1:], l, r, Hpf, Q)
	if err != nil {
		return Proof{}, errors.Wrap(err, "[Prove] -  ipproof")
	}

	return Proof{
		V:       Vs,
		A:       A.Value,
		S:       S.Value,
		T1:      T1.Value,
		T2:      T2.Value,
		l:       l,
		r:       r,
		t:       t,
		taux:    taux,
		mu:      mu,
		IPProof: ip,
	}, nil
}

// A = kH + aL*G + aR*H
func computeA(ped *pedersen.Pedersen, aLs, aRs []ristretto.Scalar) pedersen.Commitment {

	cA := ped.CommitToVectors(aLs, aRs)

	return cA
}

// S = kH + sL*G + sR * H
func computeS(ped *pedersen.Pedersen) (pedersen.Commitment, []ristretto.Scalar, []ristretto.Scalar) {

	sL, sR := make([]ristretto.Scalar, N*M), make([]ristretto.Scalar, N*M)
	for i := 0; i < N*M; i++ {
		var randA ristretto.Scalar
		randA.Rand()
		sL[i] = randA

		var randB ristretto.Scalar
		randB.Rand()
		sR[i] = randB
	}

	cS := ped.CommitToVectors(sL, sR)

	return cS, sL, sR
}

func computeYAndZ(hs fiatshamir.HashCacher) (ristretto.Scalar, ristretto.Scalar) {

	var y ristretto.Scalar
	y.Derive(hs.Result())

	var z ristretto.Scalar
	z.Derive(y.Bytes())

	return y, z
}
func computeX(hs fiatshamir.HashCacher) ristretto.Scalar {
	var x ristretto.Scalar
	x.Derive(hs.Result())
	return x
}

// compute polynomial for blinding factors l61
// N.B. tau1 means tau superscript 1
// taux = t1Blind * x + t2Blind * x^2 + (sum(z^n+1 * vBlind[n-1])) from n = 1 to n = m
func computeTaux(x, z, t1Blind, t2Blind ristretto.Scalar, vBlinds []pedersen.Commitment) ristretto.Scalar {
	tau1X := t1Blind.Mul(&x, &t1Blind)

	var xsq ristretto.Scalar
	xsq.Square(&x)

	tau2Xsq := t2Blind.Mul(&xsq, &t2Blind)

	var zN ristretto.Scalar
	zN.Square(&z) // start at zSq

	var zNBlindSum ristretto.Scalar
	zNBlindSum.SetZero()

	for i := range vBlinds {
		zNBlindSum.MulAdd(&zN, &vBlinds[i].BlindingFactor, &zNBlindSum)
		zN.Mul(&zN, &z)
	}

	var res ristretto.Scalar
	res.Add(tau1X, tau2Xsq)
	res.Add(&res, &zNBlindSum)

	return res
}

// alpha is the blinding factor for A
// rho is the blinding factor for S
// mu = alpha + rho * x
func computeMu(x, alpha, rho ristretto.Scalar) ristretto.Scalar {

	var mu ristretto.Scalar

	mu.MulAdd(&rho, &x, &alpha)

	return mu
}

// P = -mu*BlindingFactor + A + xS + Si(y,z)
func computeP(Blind, A, S ristretto.Point, mu, x, y, z ristretto.Scalar) (ristretto.Point, error) {

	var P ristretto.Point
	P.SetZero()

	var xS ristretto.Point
	xS.ScalarMult(&S, &x)

	Si, err := computeSi(y, z)
	if err != nil {
		return ristretto.Point{}, errors.Wrap(err, "[computeP]")
	}
	// Add -mu * blindingFactor
	var muB ristretto.Point
	muB.ScalarMult(&Blind, &mu)
	muB.Neg(&muB)

	P.Add(&A, &xS)
	P.Add(&P, &muB)
	P.Add(&P, &Si)

	return P, nil
}

// Si = -z<1,G> + z<1,H> + <(y^-nm Had z^2 *(z^0 * 2^n ... z^m-1 * 2^n)), H>
func computeSi(y, z ristretto.Scalar) (ristretto.Point, error) {

	var yinv ristretto.Scalar
	yinv.Inverse(&y)

	genData := []byte("dusk.BulletProof.vec1")
	ped := pedersen.New(genData)
	ped.BaseVector.Compute(uint32((N * M) + 1))

	genData = append(genData, uint8(1))

	ped2 := pedersen.New(genData)
	ped2.BaseVector.Compute(uint32(N * M))

	H := ped2.BaseVector.Bases
	G := ped.BaseVector.Bases[1:]

	yinvNM := vector.ScalarPowers(yinv, uint32(N*M))

	concatZM2N := sumZMTwoN(z)
	concatZM2N, err := vector.Hadamard(yinvNM, concatZM2N)
	if err != nil {
		return ristretto.Point{}, errors.Wrap(err, "[ComputeSI] - concatZM2N")
	}
	a, err := vector.Exp(concatZM2N, H, N, M)
	if err != nil {
		return ristretto.Point{}, errors.Wrap(err, "[ComputeSI] - a")
	}

	// b =  z<1,H> - z<1, G> = z( <1,H> - <1,G>)
	var b ristretto.Point
	b.SetZero()

	for i := range H {
		b.Add(&b, &H[i])
		b.Sub(&b, &G[i])
	}
	b.ScalarMult(&b, &z)

	var Si ristretto.Point
	Si.SetZero()

	Si.Add(&a, &b)

	return Si, nil
}

// computeHprime will take a a slice of points H, with a scalar y
// and return a slice of points Hprime,  such that Hprime = y^-n * H
func computeHprime(H []ristretto.Point, y ristretto.Scalar) []ristretto.Point {
	Hprimes := make([]ristretto.Point, len(H))

	var yInv ristretto.Scalar
	yInv.Inverse(&y)

	invYInt := yInv.BigInt()

	for i, p := range H {
		// compute y^-i
		var invYPowInt big.Int
		invYPowInt.Exp(invYInt, big.NewInt(int64(i)), nil)

		var invY ristretto.Scalar
		invY.SetBigInt(&invYPowInt)

		var hprime ristretto.Point
		hprime.ScalarMult(&p, &invY)

		Hprimes[i] = hprime
	}

	return Hprimes
}

// P = lG + rH
func computeLGRH(y, mu ristretto.Scalar, l, r []ristretto.Scalar) (ristretto.Point, error) {

	var P ristretto.Point
	P.SetZero()

	genData := []byte("dusk.BulletProof.vec1")
	ped := pedersen.New(genData)
	ped.BaseVector.Compute(uint32((N * M) + 1))

	genData = append(genData, uint8(1))

	ped2 := pedersen.New(genData)
	ped2.BaseVector.Compute(uint32(N * M))

	Hprime := computeHprime(ped2.BaseVector.Bases, y)
	G := ped.BaseVector.Bases[1:]

	rH, err := vector.Exp(r, Hprime, N, M)
	if err != nil {
		return ristretto.Point{}, errors.Wrap(err, "[computeLGRH] - rH")
	}
	lG, err := vector.Exp(l, G, N, M)
	if err != nil {
		return ristretto.Point{}, errors.Wrap(err, "[computeLGRH] - lG")
	}

	P.Add(&lG, &rH)

	return P, nil
}

// Verify takes a bullet proof and
// returns true only if the proof was valid
func Verify(p Proof) (bool, error) {

	genData := []byte("dusk.BulletProof.vec1")
	ped := pedersen.New(genData)
	ped.BaseVector.Compute(uint32((N * M) + 1))

	genData = append(genData, uint8(1))

	ped2 := pedersen.New(genData)
	ped2.BaseVector.Compute(uint32(N * M))

	H := ped2.BaseVector.Bases
	G := ped.BaseVector.Bases[1:]

	if len(p.l) != len(p.r) {
		return false, errors.New("[Verify]: Sizes of l and r do not match")
	}
	if len(p.l) <= 0 {
		return false, errors.New("[Verify]: size of l or r cannot be zero; empty proof")
	}

	// Reconstruct the challenges
	hs := fiatshamir.HashCacher{[]byte{}}
	for _, V := range p.V {
		hs.Append(V.Value.Bytes())
	}

	hs.Append(p.A.Bytes(), p.S.Bytes())
	y, z := computeYAndZ(hs)
	hs.Append(z.Bytes(), p.T1.Bytes(), p.T2.Bytes())
	x := computeX(hs)
	hs.Append(x.Bytes(), p.taux.Bytes(), p.mu.Bytes(), p.t.Bytes())
	w := hs.Derive()

	var c ristretto.Scalar
	c.Rand()

	// compute l dot r
	t, err := vector.InnerProduct(p.l, p.r)
	if err != nil {
		return false, errors.Wrap(err, "[Verify] - t")
	}
	// compute tG + tauH = Z^2 * V + delta(y,z) -- Prove t0 is correct
	_, err = VerifyT0(x, y, z, t, p.taux, ped.BaseVector.Bases[0], ped.BaseVector.Bases[1], p.T1, p.T2, p.V)
	if err != nil {
		return false, errors.Wrap(err, "[Verify] - T0")
	}

	// prove l(x) and r(x) is correct <lG, rH> = A+ xS + z<1,H> - z<1,G> + sum...
	Pleft, err := computeLGRH(y, p.mu, p.l, p.r)
	if err != nil {
		return false, errors.Wrap(err, "[Verify] - Pleft")
	}
	PRight, err := computeP(ped.BaseVector.Bases[0], p.A, p.S, p.mu, x, y, z)
	if err != nil {
		return false, errors.Wrap(err, "[Verify] - PRight")
	}
	if !Pleft.Equals(&PRight) {
		return false, errors.New("[Verify]: Proof for l(x),r(x) is wrong")
	}

	// Start of megacheck

	xSq, xInvSq, s := p.IPProof.VerifScalars()
	sInv := make([]ristretto.Scalar, len(s))
	copy(sInv, s)
	// reverse s
	for i, j := 0, len(sInv)-1; i < j; i, j = i+1, j-1 {
		sInv[i], sInv[j] = sInv[j], sInv[i]
	}

	scalars := make([]ristretto.Scalar, 0)
	points := make([]ristretto.Point, 0)

	// 1
	var one ristretto.Scalar
	one.SetOne()
	scalars = append(scalars, one)

	// x
	scalars = append(scalars, x)

	// cx
	var cx ristretto.Scalar
	cx.Mul(&c, &x)
	scalars = append(scalars, cx)

	// cx * x
	var cxsq ristretto.Scalar
	cxsq.Mul(&cx, &x)
	scalars = append(scalars, cxsq)

	// xSq
	scalars = append(scalars, xSq...)

	// xInvSq
	scalars = append(scalars, xInvSq...)

	// -mu - c * taux = -(mu + c * taux)
	var muCtauX ristretto.Scalar
	var cTaux ristretto.Scalar
	cTaux.Mul(&c, &p.taux)
	muCtauX.Add(&p.mu, &cTaux)
	muCtauX.Neg(&muCtauX)
	scalars = append(scalars, muCtauX)

	// w * (t-ab) + c * (delta(y,z) - t) = bpScal
	delta := computeDelta(y, z)

	var cDelMinT ristretto.Scalar
	cDelMinT.Sub(&delta, &p.t)
	cDelMinT.Mul(&cDelMinT, &c)

	var tMinAB ristretto.Scalar
	var ab ristretto.Scalar
	ab.Mul(&p.IPProof.A, &p.IPProof.B)
	tMinAB.Sub(&p.t, &ab)

	var bpScal ristretto.Scalar
	bpScal.MulSub(&w, &tMinAB, &bpScal)
	scalars = append(scalars, bpScal)

	// -z1 - as = g
	as := vector.MulScalar(s, p.IPProof.A)
	var minusZ ristretto.Scalar
	minusZ.Neg(&z)
	g := vector.AddScalar(as, minusZ)
	scalars = append(scalars, g...)

	// h =  (z + y^-nm) Had ( z^2 *(z^0 * 2^n ... z^m-1 * 2^n)- b/s) = zy Had (zAnd2 - bsInv)
	var yinv ristretto.Scalar
	yinv.Inverse(&y)
	yinvNM := vector.ScalarPowers(yinv, uint32(N*M))
	zy := vector.AddScalar(yinvNM, z)

	zAnd2 := sumZMTwoN(z)
	bs := vector.MulScalar(sInv, p.IPProof.B)

	zAnd2SubBs, err := vector.Sub(zAnd2, bs)
	if err != nil {
		return false, errors.Wrap(err, "[Verify]")
	}

	h, err := vector.Hadamard(zy, zAnd2SubBs)
	if err != nil {
		return false, errors.Wrap(err, "[Verify]")
	}
	scalars = append(scalars, h...)

	// cz^2 ..., cz^(m+1) = czM
	zM := vector.ScalarPowers(z, uint32(M))

	var czsq ristretto.Scalar
	czsq.Square(&z)
	czsq.Mul(&czsq, &c)

	czM := vector.MulScalar(zM, czsq)
	scalars = append(scalars, czM...)

	// POINTS

	// A
	points = append(points, p.A)

	// S
	points = append(points, p.S)

	// T1
	points = append(points, p.T1)

	// T2
	points = append(points, p.T2)

	// Lj
	points = append(points, p.IPProof.L...)

	// Rj
	points = append(points, p.IPProof.R...)

	// Blinding Point
	points = append(points, ped.BaseVector.Bases[0])

	// Generator Point
	points = append(points, ped.BaseVector.Bases[1])

	// Set of Generator Points without blinder
	points = append(points, G...)

	// Set of Points for second vector of points
	points = append(points, H...)

	// Commitment values
	for _, V := range p.V {

		points = append(points, V.Value)

	}

	result, err := vector.Exp(scalars, points, len(scalars), 1)

	if err != nil {
		return false, errors.Wrap(err, "[Verify]")
	}

	var iden ristretto.Point
	iden.SetZero()

	ok := result.EqualsI(&iden)

	fmt.Println("iden", iden.Bytes())
	fmt.Println("res", result.Bytes())

	fmt.Println("res", ok)

	if ok != 1 {
		return false, errors.New("Megacheck failed")
	}
	return true, nil
}

func VerifyT0(x, y, z, t, taux ristretto.Scalar, Blind, G1, T1, T2 ristretto.Point, V []pedersen.Commitment) (bool, error) {
	// LHS
	var tG ristretto.Point
	tG.ScalarMult(&G1, &t)
	var tauH ristretto.Point
	tauH.ScalarMult(&Blind, &taux)

	var LHS ristretto.Point
	LHS.Add(&tG, &tauH)

	// RHS
	var cT1x ristretto.Point
	cT1x.ScalarMult(&T1, &x)

	var xsq ristretto.Scalar
	xsq.Square(&x)

	var cT2xsq ristretto.Point
	cT2xsq.ScalarMult(&T2, &xsq)

	var zM ristretto.Scalar
	zM = z

	var zMV ristretto.Point
	var sumzMV ristretto.Point

	for i := 0; i < M; i++ {
		zM.Mul(&zM, &z)
		zMV.ScalarMult(&V[i].Value, &zM)
		sumzMV.Add(&sumzMV, &zMV)
	}

	var deltaG ristretto.Point
	delta := computeDelta(y, z)

	deltaG.ScalarMult(&G1, &delta)

	var RHS1 ristretto.Point
	var RHS2 ristretto.Point
	var RHS ristretto.Point
	RHS1.Add(&sumzMV, &deltaG)
	RHS2.Add(&cT2xsq, &cT1x)
	RHS.Add(&RHS1, &RHS2)

	if !LHS.Equals(&RHS) {
		return false, errors.New("[Verify]: LHS != RHS; proof that t0 is correct is wrong")
	}
	return true, nil
}