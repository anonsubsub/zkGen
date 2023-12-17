package templates

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strconv"
	"time"
	c "transpiler/circuits"

	"github.com/rs/zerolog/log"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	gnarkHash "github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark-crypto/kzg"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/plonkfri"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test"
)

func TestCircuit(policyName, solidityFlag string) error {

	// sample age
	age := "23"
	threshold := "18"
	bytesAge, _ := hex.DecodeString(age)
	ageString := hex.EncodeToString(bytesAge)
	ageLen := len(bytesAge)
	witRandBytes := make([]byte, ageLen)
	rand.Read(witRandBytes)
	witRandString := hex.EncodeToString(witRandBytes) // used to be plain

	hashInput := make([]byte, ageLen)
	// xor witness with age bytes
	for j := 0; j < ageLen; j++ {
		hashInput[j] = witRandBytes[j] ^ bytesAge[j]
	}
	// hashInputString := hex.EncodeToString(hashInput)
	// hashInputBigInt, _ := new(big.Int).SetString(hashInputString, 16)
	// curve := ecc.BN254
	// modulus := curve.ScalarField()
	// hashInputBigInt.Mod(hashInputBigInt, modulus)
	// zkInputBigInt := new(big.Int).SetBytes(hashInput)

	hashFunc := gnarkHash.MIMC_BN254
	goMimc := hashFunc.New()
	var x fr.Element // fr element works for smaller inputs, check if fr.Hash must be used for larger inputs.
	x.SetBytes(hashInput)
	// elems, _ := fr.Hash(hashInput, []byte("string:"), 1)
	// bs := elems[0].Bytes()
	// goMimc.Write(bs[:])
	bs := x.Bytes()
	goMimc.Write(bs[:])
	expectedh := goMimc.Sum(nil)

	backend := "plonk"
	compile := false // true just compiles css and returns

	witAssign := StrToIntSlice(witRandString, true)
	ageAssign := StrToIntSlice(ageString, true)

	assignment := c.MimcWrapper{
		Age:       make([]frontend.Variable, ageLen),
		Witness:   make([]frontend.Variable, ageLen),
		Hash:      expectedh,
		Threshold: threshold,
	}
	for j := 0; j < ageLen; j++ {
		assignment.Age[j] = ageAssign[j]
	}
	for i := 0; i < ageLen; i++ {
		assignment.Witness[i] = witAssign[i]
	}

	circuit := c.MimcWrapper{
		Age:     make([]frontend.Variable, ageLen),
		Witness: make([]frontend.Variable, ageLen),
	}

	_, err := ProofWithBackend(backend, compile, &circuit, &assignment, ecc.BN254, solidityFlag, policyName)
	if err != nil {
		log.Error().Err(err).Msg("ProofWithBackend()")
		return err
	}

	return nil
}

// non-gnark zk system evalaution functions
func ProofWithBackend(backend string, compile bool, circuit frontend.Circuit, assignment frontend.Circuit, curveID ecc.ID, solidityFlag, policyName string) (map[string]time.Duration, error) {

	// time measures
	data := map[string]time.Duration{}

	// generate witness
	witness, err := frontend.NewWitness(assignment, curveID.ScalarField())
	if err != nil {
		log.Error().Msg("frontend.NewWitness")
		return nil, err
	}

	// init builders
	var builder frontend.NewBuilder
	var srs kzg.SRS
	switch backend {
	case "groth16":
		builder = r1cs.NewBuilder
	case "plonk":
		builder = scs.NewBuilder
	case "plonkFRI":
		builder = scs.NewBuilder
	}

	// generate CompiledConstraintSystem
	start := time.Now()
	ccs, err := frontend.Compile(curveID.ScalarField(), builder, circuit)
	if err != nil {
		log.Error().Msg("frontend.Compile")
		return nil, err
	}
	elapsed := time.Since(start)
	log.Debug().Str("elapsed", elapsed.String()).Msg("compile constraint system time.")

	data["compile"] = elapsed

	// measure byte size
	var buf bytes.Buffer
	bytesWritten, err := ccs.WriteTo(&buf)
	if err != nil {
		log.Error().Msg("ccs serialization error")
		return nil, err
	}
	log.Debug().Str("written", strconv.FormatInt(bytesWritten, 10)).Msg("compiled constraint system bytes")

	// kzg setup if using plonk
	if backend == "plonk" {
		srs, err = test.NewKZGSRS(ccs)
		if err != nil {
			log.Error().Msg("test.NewKZGSRS(ccs)")
			return nil, err
		}

		elapsed := time.Since(start)
		data["compile"] = elapsed

		// measure byte size
		var bufExtra bytes.Buffer
		bytesWritten, err := srs.WriteTo(&bufExtra)
		if err != nil {
			log.Error().Msg("srs serialization error")
			return nil, err
		}
		log.Debug().Str("written srs", strconv.FormatInt(bytesWritten, 10)).Msg("compiled constraint system bytes")

	}

	if compile {
		return nil, err
	}

	// proof system execution
	switch backend {
	case "groth16":

		// setup
		start = time.Now()
		pk, vk, err := groth16.Setup(ccs)
		if err != nil {
			log.Error().Msg("groth16.Setup")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("groth16.Setup time.")

		data["setup"] = elapsed

		// export solidity
		if solidityFlag == "solidity" {
			f, err := os.Create("circuits/" + policyName + ".sol")
			if err != nil {
				log.Error().Msg("os creation error")
				return nil, err
			}
			err = vk.ExportSolidity(f)
			if err != nil {
				log.Error().Msg("export solidity error")
				return nil, err
			}
		}

		// measure byte size
		var buf1 bytes.Buffer
		bytesWritten1, err := pk.WriteTo(&buf1)
		if err != nil {
			log.Error().Msg("ccs serialization error")
			return nil, err
		}
		log.Debug().Str("written", strconv.FormatInt(bytesWritten1, 10)).Msg("prover key bytes")
		var buf2 bytes.Buffer
		bytesWritten2, err := vk.WriteTo(&buf2)
		if err != nil {
			log.Error().Msg("ccs serialization error")
			return nil, err
		}
		log.Debug().Str("written", strconv.FormatInt(bytesWritten2, 10)).Msg("verifier key bytes")

		// prove
		start = time.Now()
		proof, err := groth16.Prove(ccs, pk, witness)
		if err != nil {
			log.Error().Msg("groth16.Prove")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("groth16.Prove time.")

		data["prove"] = elapsed

		// measure bytes
		var buf3 bytes.Buffer
		bytesWritten3, err := proof.WriteTo(&buf3)
		if err != nil {
			log.Error().Msg("proof serialization error")
			return nil, err
		}
		log.Debug().Str("written", strconv.FormatInt(bytesWritten3, 10)).Msg("proof bytes")

		// generate public witness
		publicWitness, _ := witness.Public()

		// measure bytes
		witnessBytes, err := publicWitness.MarshalBinary()
		if err != nil {
			log.Error().Msg("witness marshal binary error")
			return nil, err
		}
		log.Debug().Str("written", strconv.Itoa(len(witnessBytes)/8)).Msg("witness bytes")

		// verification
		start = time.Now()
		err = groth16.Verify(proof, vk, publicWitness)
		if err != nil {
			log.Error().Msg("groth16.Verify")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("groth16.Verify time.")

		data["verify"] = elapsed

	case "plonk":

		// setup
		start = time.Now()
		pk, vk, err := plonk.Setup(ccs, srs)
		if err != nil {
			log.Error().Msg("plonk.Setup")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("plonk.Setup time.")

		data["setup"] = elapsed

		// export solidity
		if solidityFlag == "solidity" {
			f, err := os.Create("circuits/" + policyName + ".sol")
			if err != nil {
				log.Error().Msg("os creation error")
				return nil, err
			}
			err = vk.ExportSolidity(f)
			if err != nil {
				log.Error().Msg("export solidity error")
				return nil, err
			}
		}

		// measure byte size
		var buf1 bytes.Buffer
		bytesWritten1, err := pk.WriteTo(&buf1)
		if err != nil {
			log.Error().Msg("ccs serialization error")
			return nil, err
		}
		log.Debug().Str("written", strconv.FormatInt(bytesWritten1, 10)).Msg("prover key bytes")
		var buf2 bytes.Buffer
		bytesWritten2, err := vk.WriteTo(&buf2)
		if err != nil {
			log.Error().Msg("ccs serialization error")
			return nil, err
		}
		log.Debug().Str("written", strconv.FormatInt(bytesWritten2, 10)).Msg("verifier key bytes")

		// prove
		start = time.Now()
		proof, err := plonk.Prove(ccs, pk, witness)
		if err != nil {
			log.Error().Msg("plonk.Prove")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("plonk.Prove time.")

		data["prove"] = elapsed

		// measure bytes
		var buf3 bytes.Buffer
		bytesWritten3, err := proof.WriteTo(&buf3)
		if err != nil {
			log.Error().Msg("proof serialization error")
			return nil, err
		}
		log.Debug().Str("written", strconv.FormatInt(bytesWritten3, 10)).Msg("proof bytes")

		// generate public witness
		publicWitness, _ := witness.Public()

		// measure bytes
		witnessBytes, err := publicWitness.MarshalBinary()
		if err != nil {
			log.Error().Msg("witness marshal binary error")
			return nil, err
		}
		// fmt.Println("witnessBytes:", witnessBytes)
		log.Debug().Str("written", strconv.Itoa(len(witnessBytes)/8)).Msg("witness bytes")

		// verify
		start = time.Now()
		err = plonk.Verify(proof, vk, publicWitness)
		if err != nil {
			log.Error().Msg("plonk.Verify")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("plonk.Verify time.")

		data["verify"] = elapsed

	case "plonkFRI":

		// setup
		start = time.Now()
		pk, vk, err := plonkfri.Setup(ccs)
		if err != nil {
			log.Error().Msg("plonkfri.Setup")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("plonkfri.Setup time.")

		data["setup"] = elapsed

		// prove
		start = time.Now()
		correctProof, err := plonkfri.Prove(ccs, pk, witness)
		if err != nil {
			log.Error().Msg("plonkfri.Prove")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("plonkfri.Prove time.")

		data["prove"] = elapsed

		// generate public witness
		publicWitness, _ := witness.Public()

		// verify
		start = time.Now()
		err = plonkfri.Verify(correctProof, vk, publicWitness)
		if err != nil {
			log.Error().Msg("plonkfri.Verify")
			return nil, err
		}
		elapsed = time.Since(start)
		log.Debug().Str("elapsed", elapsed.String()).Msg("plonkfri.Verify time.")

		data["verify"] = elapsed
	}

	return data, nil
}

func StrToIntSlice(inputData string, hexRepresentation bool) []int {

	// check if inputData in hex representation
	var byteSlice []byte
	if hexRepresentation {
		hexBytes, err := hex.DecodeString(inputData)
		if err != nil {
			log.Error().Msg("hex.DecodeString error.")
		}
		byteSlice = hexBytes
	} else {
		byteSlice = []byte(inputData)
	}

	// convert byte slice to int numbers which can be passed to gnark frontend.Variable
	var data []int
	for i := 0; i < len(byteSlice); i++ {
		data = append(data, int(byteSlice[i]))
	}

	return data
}
