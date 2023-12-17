package templates

import (
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"

	p "transpiler/zkpolicy"
)

type Circuit struct {
	FileName      string
	Template      *template.Template
	Constraints   []string
	Relations     []p.Relation
	CircuitString string
}

func NewCircuit(policy p.ZkPolicy) *Circuit {

	algorithm := policy.Relations[0].PrivateArgument.Protection.Algorithm
	circuitModel := GetCircuitModel(algorithm)

	t := template.Must(template.New("zkcircuit").Parse(circuitModel))
	c := &Circuit{Template: t, FileName: policy.Name, Constraints: policy.Constraints, Relations: policy.Relations}

	return c
}

func (t *Circuit) Transpile() error {

	operator := strings.Split(t.Constraints[0], "-")[1]

	// template data
	var data map[string]string
	if operator == "gt" {
		data = dataAttributesGT(t.Relations[0].PrivateArgument.Name, t.Relations[0].PublicArgument.Name)
	} else if operator == "lt" {
		data = dataAttributesLT(t.Relations[0].PrivateArgument.Name, t.Relations[0].PublicArgument.Name)
	} else if operator == "eq" {
		data = getComparatorEQ()
	} else {
		err := errors.New("constraint not supported")
		log.Error().Err(err).Msg("transpile error: constraint not found")
		return err
	}

	// transpile constraints into generator file
	builder := &strings.Builder{}
	if err := t.Template.Execute(builder, data); err != nil {
		log.Error().Err(err).Msg("t.Template.Execute()")
		return err
	}
	t.CircuitString = builder.String()

	return nil
}

func (t *Circuit) StoreCircuit() error {

	// store generator file in jsnark-demo
	generatorFile, err := os.OpenFile("circuits/"+t.FileName+".go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Error().Err(err).Msg("os.OpenFile()")
		return err
	}
	_, err = generatorFile.WriteString(t.CircuitString)
	defer generatorFile.Close()
	if err != nil {
		log.Error().Err(err).Msg("WriteString()")
		return err
	}

	return nil
}

func dataAttributesLT(privateVariable, publicVariable string) map[string]string {

	data := map[string]string{
		"PrivateVariable":           privateVariable,
		"PrivateVariableDefinition": privateVariable + " []frontend.Variable",
		"Public":                    "`gnark:\",public\"`",
		"PublicVariableDefinition":  publicVariable + " frontend.Variable `gnark:\",public\"`",
		"Comparison":                "api.AssertIsLessOrEqual(sum , circuit." + publicVariable + ")",
	}

	return data
}

func dataAttributesGT(privateVariable, publicVariable string) map[string]string {

	data := map[string]string{
		"PrivateVariable":           privateVariable,
		"PrivateVariableDefinition": privateVariable + " []frontend.Variable",
		"Public":                    "`gnark:\",public\"`",
		"PublicVariableDefinition":  publicVariable + " frontend.Variable `gnark:\",public\"`",
		"Comparison":                "api.AssertIsLessOrEqual(circuit." + publicVariable + ", sum)",
	}

	return data
}

func getComparatorEQ() map[string]string {

	data := map[string]string{
		"GeneratorThreshold": "args[24]",
	}

	// sample_data := map[string]interface{}{
	// "Name":     "Bob",
	// "UserName": "bob92",
	// "Roles":    []string{"dbteam", "uiteam", "tester"},
	// }

	// `gnark:",public"`

	return data
}
