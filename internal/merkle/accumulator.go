package merkle

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"math/big"
)

// RSAAccumulator implements a simple RSA-based cryptographic accumulator
type RSAAccumulator struct {
	N           *big.Int
	G           *big.Int
	Acc         *big.Int
	elemProduct *big.Int
}

// NewRSAAccumulator generates an RSA modulus and initializes the accumulator
func NewRSAAccumulator(bits int) (*RSAAccumulator, error) {
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	N := key.N
	G := big.NewInt(2)
	return &RSAAccumulator{
		N:           N,
		G:           G,
		Acc:         new(big.Int).Set(G),
		elemProduct: big.NewInt(1),
	}, nil
}

// Add includes a prime element into the accumulator
func (r *RSAAccumulator) Add(element *big.Int) error {
	if !element.ProbablyPrime(20) {
		return errors.New("element must be prime")
	}
	r.elemProduct.Mul(r.elemProduct, element)
	r.elemProduct.Mod(r.elemProduct, new(big.Int).Sub(r.N, big.NewInt(1)))
	r.Acc.Exp(r.G, r.elemProduct, r.N)
	return nil
}

// Witness returns the witness for a given element
func (r *RSAAccumulator) Witness(element *big.Int) (*big.Int, error) {
	if !element.ProbablyPrime(20) {
		return nil, errors.New("element must be prime")
	}
	inv := new(big.Int).ModInverse(element, new(big.Int).Sub(r.N, big.NewInt(1)))
	if inv == nil {
		return nil, errors.New("no modular inverse")
	}
	exp := new(big.Int).Mul(r.elemProduct, inv)
	exp.Mod(exp, new(big.Int).Sub(r.N, big.NewInt(1)))
	w := new(big.Int).Exp(r.G, exp, r.N)
	return w, nil
}

// Verify checks that witness^element ≡ accumulator (mod N)
func (r *RSAAccumulator) Verify(element, witness *big.Int) bool {
	lhs := new(big.Int).Exp(witness, element, r.N)
	return lhs.Cmp(r.Acc) == 0
}
