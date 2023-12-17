package templates

func GetCircuitModel(algorithm string) string {

	if algorithm == "commitments:mimc" {
		// return parsed gadget
		return `
package circuits

import (
	"math"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// Circuit defines a pre-image knowledge proof
// mimc(secret preImage) = public hash
type MimcWrapper struct {
	{{.PrivateVariableDefinition}}
	Witness []frontend.Variable
	Hash frontend.Variable {{.Public}}
	{{.PublicVariableDefinition}}
}

// Define declares the circuit's constraints
func (circuit *MimcWrapper) Define(api frontend.API) error {
	
	mimc, _ := mimc.NewMiMC(api)

	// rearrange input to match mimc input requirements
	ddd := make([]frontend.Variable, len(circuit.{{.PrivateVariable}})*8) // used to be 256
	for j := 0; j < len(circuit.{{.PrivateVariable}}); j++ {
		witnessBits := api.ToBinary(circuit.Witness[j], 8)
		myBits := api.ToBinary(circuit.{{.PrivateVariable}}[j], 8)
		// get bits of ecb input, little endian!
		for k := 7; k >= 0; k-- {
			ddd[((len(circuit.{{.PrivateVariable}})-1-j)*8)+k] = api.Xor(witnessBits[k], myBits[k])
		}
	}
	// input data into mimc
	varSum := api.FromBinary(ddd...)
	mimc.Write(varSum)
	result := mimc.Sum()
	api.AssertIsEqual(circuit.Hash, result)

	// string to integer conversion
	valueString := circuit.{{.PrivateVariable}}[0:len(circuit.{{.PrivateVariable}})]

	// aggregation number
	sum := frontend.Variable(0)
	// loop from back to front
	for i := len(valueString); i > 0; i-- {
		idx := len(valueString) - i

		// expanded dezimal such that shift can be applied
		// 4 bits cover numbers 0-9, little endian result, IMPORTANT: 8 required, otherwise unsatisfied constraint error
		toInt := api.Sub(api.FromBinary(api.ToBinary(valueString[i-1], 8)...), 48)
		sum = api.MulAcc(sum, toInt, int(math.Pow(float64(10), float64(idx))))
	}

	{{.Comparison}}

	return nil
}	
		`
	} else {
		return ``
	}
}
