package execution

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/cube2222/octosql"
)

type Datatype string

const (
	DatatypeBoolean Datatype = "boolean"
	DatatypeInt     Datatype = "int"
	DatatypeFloat64 Datatype = "float64"
	DatatypeString  Datatype = "string"
	DatatypeTuple   Datatype = "octosql.Tuple"
)

func GetType(i octosql.Value) Datatype {
	if _, ok := i.(octosql.Bool); ok {
		return DatatypeBoolean
	}
	if _, ok := i.(octosql.Int); ok {
		return DatatypeInt
	}
	if _, ok := i.(octosql.Float); ok {
		return DatatypeFloat64
	}
	if _, ok := i.(octosql.String); ok {
		return DatatypeString
	}
	if _, ok := i.(octosql.Tuple); ok {
		return DatatypeTuple
	}
	return DatatypeString // TODO: Unknown
}

// ParseType tries to parse the given string into any type it succeeds to. Returns back the string on failure.
func ParseType(str string) octosql.Value {
	integer, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return octosql.MakeInt(int(integer))
	}

	float, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return octosql.MakeFloat(float)
	}

	boolean, err := strconv.ParseBool(str)
	if err == nil {
		return octosql.MakeBool(boolean)
	}

	var jsonObject map[string]interface{}
	err = json.Unmarshal([]byte(str), &jsonObject)
	if err == nil {
		return octosql.NormalizeType(jsonObject)
	}

	t, err := time.Parse(time.RFC3339Nano, str)
	if err == nil {
		return octosql.MakeTime(t)
	}

	return octosql.MakeString(str)
}
