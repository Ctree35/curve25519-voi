// Copyright (c) 2016-2019 isis agora lovecruft. All rights reserved.
// Copyright (c) 2016-2019 Henry de Valence. All rights reserved.
// Copyright (c) 2020-2021 Oasis Labs Inc. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// 1. Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
// notice, this list of conditions and the following disclaimer in the
// documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS
// IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED
// TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
// PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
// TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
// PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
// LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
// NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package curve

import "github.com/oasisprotocol/curve25519-voi/internal/field"

type ProjectivePoint struct {
	X field.Element
	Y field.Element
	Z field.Element
}

type CompletedPoint struct {
	X field.Element
	Y field.Element
	Z field.Element
	T field.Element
}

type AffineNielsPoint struct {
	y_plus_x  field.Element
	y_minus_x field.Element
	xy2d      field.Element
}

type ProjectiveNielsPoint struct {
	Y_plus_X  field.Element
	Y_minus_X field.Element
	Z         field.Element
	T2d       field.Element
}

func (p *AffineNielsPoint) SetRaw(raw *[96]uint8) *AffineNielsPoint {
	_, _ = p.y_plus_x.SetBytes(raw[0:32])
	_, _ = p.y_minus_x.SetBytes(raw[32:64])
	_, _ = p.xy2d.SetBytes(raw[64:96])
	return p
}

// Note: dalek has the identity point as the defaut ctors for
// ProjectiveNielsPoint/AffineNielsPoint.

func (p *ProjectivePoint) Identity() *ProjectivePoint {
	p.X.Zero()
	p.Y.One()
	p.Z.One()
	return p
}

func (p *AffineNielsPoint) Identity() *AffineNielsPoint {
	p.y_plus_x.One()
	p.y_minus_x.One()
	p.xy2d.Zero()
	return p
}

func (p *ProjectiveNielsPoint) Identity() *ProjectiveNielsPoint {
	p.Y_plus_X.One()
	p.Y_minus_X.One()
	p.Z.One()
	p.T2d.Zero()
	return p
}

func (p *ProjectiveNielsPoint) ConditionalSelect(a, b *ProjectiveNielsPoint, choice int) {
	p.Y_plus_X.ConditionalSelect(&a.Y_plus_X, &b.Y_plus_X, choice)
	p.Y_minus_X.ConditionalSelect(&a.Y_minus_X, &b.Y_minus_X, choice)
	p.Z.ConditionalSelect(&a.Z, &b.Z, choice)
	p.T2d.ConditionalSelect(&a.T2d, &b.T2d, choice)
}

func (p *ProjectiveNielsPoint) ConditionalAssign(other *ProjectiveNielsPoint, choice int) {
	p.Y_plus_X.ConditionalAssign(&other.Y_plus_X, choice)
	p.Y_minus_X.ConditionalAssign(&other.Y_minus_X, choice)
	p.Z.ConditionalAssign(&other.Z, choice)
	p.T2d.ConditionalAssign(&other.T2d, choice)
}

func (p *AffineNielsPoint) ConditionalSelect(a, b *AffineNielsPoint, choice int) {
	p.y_plus_x.ConditionalSelect(&a.y_plus_x, &b.y_plus_x, choice)
	p.y_minus_x.ConditionalSelect(&a.y_minus_x, &b.y_minus_x, choice)
	p.xy2d.ConditionalSelect(&a.xy2d, &b.xy2d, choice)
}

func (p *AffineNielsPoint) ConditionalAssign(other *AffineNielsPoint, choice int) {
	p.y_plus_x.ConditionalAssign(&other.y_plus_x, choice)
	p.y_minus_x.ConditionalAssign(&other.y_minus_x, choice)
	p.xy2d.ConditionalAssign(&other.xy2d, choice)
}

func (p *EdwardsPoint) setProjective(pp *ProjectivePoint) *EdwardsPoint {
	p.Inner.X.Mul(&pp.X, &pp.Z)
	p.Inner.Y.Mul(&pp.Y, &pp.Z)
	p.Inner.Z.Square(&pp.Z)
	p.Inner.T.Mul(&pp.X, &pp.Y)
	return p
}

func (p *EdwardsPoint) setAffineNiels(ap *AffineNielsPoint) *EdwardsPoint {
	p.Identity()

	var sum CompletedPoint
	return p.setCompleted(sum.AddEdwardsAffineNiels(p, ap))
}

func (p *EdwardsPoint) setCompleted(cp *CompletedPoint) *EdwardsPoint {
	p.Inner.X.Mul(&cp.X, &cp.T)
	p.Inner.Y.Mul(&cp.Y, &cp.Z)
	p.Inner.Z.Mul(&cp.Z, &cp.T)
	p.Inner.T.Mul(&cp.X, &cp.Y)
	return p
}

func (p *ProjectivePoint) SetCompleted(cp *CompletedPoint) *ProjectivePoint {
	p.X.Mul(&cp.X, &cp.T)
	p.Y.Mul(&cp.Y, &cp.Z)
	p.Z.Mul(&cp.Z, &cp.T)
	return p
}

func (p *ProjectivePoint) SetEdwards(ep *EdwardsPoint) *ProjectivePoint {
	p.X.Set(&ep.Inner.X)
	p.Y.Set(&ep.Inner.Y)
	p.Z.Set(&ep.Inner.Z)
	return p
}

func (p *ProjectiveNielsPoint) SetEdwards(ep *EdwardsPoint) *ProjectiveNielsPoint {
	p.Y_plus_X.Add(&ep.Inner.Y, &ep.Inner.X)
	p.Y_minus_X.Sub(&ep.Inner.Y, &ep.Inner.X)
	p.Z.Set(&ep.Inner.Z)
	p.T2d.Mul(&ep.Inner.T, &constEDWARDS_D2)
	return p
}

func (p *AffineNielsPoint) SetEdwards(ep *EdwardsPoint) *AffineNielsPoint {
	var recip, x, y, xy field.Element
	recip.Invert(&ep.Inner.Z)
	x.Mul(&ep.Inner.X, &recip)
	y.Mul(&ep.Inner.Y, &recip)
	xy.Mul(&x, &y)
	p.y_plus_x.Add(&y, &x)
	p.y_minus_x.Sub(&y, &x)
	p.xy2d.Mul(&xy, &constEDWARDS_D2)
	return p
}

func (p *CompletedPoint) Double(pp *ProjectivePoint) *CompletedPoint {
	var XX, YY, ZZ2, X_plus_Y_sq field.Element
	XX.Square(&pp.X)
	YY.Square(&pp.Y)
	ZZ2.Square2(&pp.Z)
	X_plus_Y_sq.Add(&pp.X, &pp.Y)    // X+Y
	X_plus_Y_sq.Square(&X_plus_Y_sq) // (X+Y)^2

	p.Y.Add(&YY, &XX)
	p.X.Sub(&X_plus_Y_sq, &p.Y)
	p.Z.Sub(&YY, &XX)
	p.T.Sub(&ZZ2, &p.Z)

	return p
}

func (p *CompletedPoint) AddEdwardsProjectiveNiels(a *EdwardsPoint, b *ProjectiveNielsPoint) *CompletedPoint {
	var PP, MM, TT2d, ZZ, ZZ2 field.Element
	PP.Add(&a.Inner.Y, &a.Inner.X) // a.Y + a.X
	PP.Mul(&PP, &b.Y_plus_X)       // (a.Y + a.X) * b.Y_plus_X
	MM.Sub(&a.Inner.Y, &a.Inner.X) // a.Y - a.X
	MM.Mul(&MM, &b.Y_minus_X)      // (a.Y - a.X) * b.Y_minus_X
	TT2d.Mul(&a.Inner.T, &b.T2d)
	ZZ.Mul(&a.Inner.Z, &b.Z)
	ZZ2.Add(&ZZ, &ZZ)

	p.X.Sub(&PP, &MM)
	p.Y.Add(&PP, &MM)
	p.Z.Add(&ZZ2, &TT2d)
	p.T.Sub(&ZZ2, &TT2d)

	return p
}

func (p *CompletedPoint) SubEdwardsProjectiveNiels(a *EdwardsPoint, b *ProjectiveNielsPoint) *CompletedPoint {
	var PM, MP, TT2d, ZZ, ZZ2 field.Element
	PM.Add(&a.Inner.Y, &a.Inner.X) // a.Y + a.X
	PM.Mul(&PM, &b.Y_minus_X)      // (a.Y + a.X) * b.Y_minus_X
	MP.Sub(&a.Inner.Y, &a.Inner.X) // a.Y - a.X
	MP.Mul(&MP, &b.Y_plus_X)       // (a.Y - a.X) * b.Y_plus_X
	TT2d.Mul(&a.Inner.T, &b.T2d)
	ZZ.Mul(&a.Inner.Z, &b.Z)
	ZZ2.Add(&ZZ, &ZZ)

	p.X.Sub(&PM, &MP)
	p.Y.Add(&PM, &MP)
	p.Z.Sub(&ZZ2, &TT2d)
	p.T.Add(&ZZ2, &TT2d)
	return p
}

func (p *CompletedPoint) AddEdwardsAffineNiels(a *EdwardsPoint, b *AffineNielsPoint) *CompletedPoint {
	var PP, MM, Txy2d, Z2 field.Element
	PP.Add(&a.Inner.Y, &a.Inner.X) // a.Y + a.X
	PP.Mul(&PP, &b.y_plus_x)       // (a.Y + a.X) * b.y_plus_x
	MM.Sub(&a.Inner.Y, &a.Inner.X) // a.Y - a.X
	MM.Mul(&MM, &b.y_minus_x)      // (a.Y - a.X) * b.y_minus_x
	Txy2d.Mul(&a.Inner.T, &b.xy2d)
	Z2.Add(&a.Inner.Z, &a.Inner.Z)

	p.X.Sub(&PP, &MM)
	p.Y.Add(&PP, &MM)
	p.Z.Add(&Z2, &Txy2d)
	p.T.Sub(&Z2, &Txy2d)

	return p
}

func (p *CompletedPoint) AddCompletedAffineNiels(a *CompletedPoint, b *AffineNielsPoint) *CompletedPoint {
	var aTmp EdwardsPoint
	return p.AddEdwardsAffineNiels(aTmp.setCompleted(a), b)
}

func (p *CompletedPoint) SubEdwardsAffineNiels(a *EdwardsPoint, b *AffineNielsPoint) *CompletedPoint {
	var Y_plus_X, Y_minus_X, PM, MP, Txy2d, Z2 field.Element
	Y_plus_X.Add(&a.Inner.Y, &a.Inner.X)
	Y_minus_X.Sub(&a.Inner.Y, &a.Inner.X)
	PM.Mul(&Y_plus_X, &b.y_minus_x)
	MP.Mul(&Y_minus_X, &b.y_plus_x)
	Txy2d.Mul(&a.Inner.T, &b.xy2d)
	Z2.Add(&a.Inner.Z, &a.Inner.Z)

	p.X.Sub(&PM, &MP)
	p.Y.Add(&PM, &MP)
	p.Z.Sub(&Z2, &Txy2d)
	p.T.Add(&Z2, &Txy2d)

	return p
}

func (p *CompletedPoint) SubCompletedAffineNiels(a *CompletedPoint, b *AffineNielsPoint) *CompletedPoint {
	var aTmp EdwardsPoint
	return p.SubEdwardsAffineNiels(aTmp.setCompleted(a), b)
}

func (p *ProjectiveNielsPoint) ConditionalNegate(choice int) {
	p.Y_plus_X.ConditionalSwap(&p.Y_minus_X, choice)
	p.T2d.ConditionalNegate(choice)
}

func (p *AffineNielsPoint) ConditionalNegate(choice int) {
	p.y_plus_x.ConditionalSwap(&p.y_minus_x, choice)
	p.xy2d.ConditionalNegate(choice)
}
