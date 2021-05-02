package wallet

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

const (
	MetadataTypeInt    = "int"
	MetadataTypeString = "string"
	MetadataTypeBytes  = "bytes"
	MetadataTypeList   = "list"
	MetadataTypeMap    = "map"

	MetadataMapKey = "k"
	MetadataMapVal = "v"
)

// Metadata represents transaction metadata, as it is stored on the ledger.
// This type is manually added here, because oapi-codegen fails to generate it.
// The doc-comment below is copied from the Swagger definition.
//
// **⚠️ WARNING ⚠️**
//
// _Please note that metadata provided in a transaction will be
// stored on the blockchain forever. Make sure not to include any sensitive data,
// in particular personally identifiable information (PII)._
//
// Extra application data attached to the transaction.
//
// Cardano allows users and developers to embed their own
// authenticated metadata when submitting transactions. Metadata can
// be expressed as a JSON object with some restrictions:
//
// 1. All top-level keys must be integers between `0` and `2^64 - 1`.
//
// 2. Each metadata value is tagged with its type.
//
// 3. Strings must be at most 64 bytes when UTF-8 encoded.
//
// 4. Bytestrings are hex-encoded, with a maximum length of 64 bytes.
//
// Metadata aren't stored as JSON on the Cardano blockchain but are
// instead stored using a compact binary encoding (CBOR).
//
// The binary encoding of metadata values supports three simple types:
//
// * Integers in the range `-(2^64 - 1)` to `2^64 - 1`
// * Strings (UTF-8 encoded)
// * Bytestrings
//
// And two compound types:
//
// * Lists of metadata values
// * Mappings from metadata values to metadata values
//
// It is possible to transform any JSON object into this schema.
//
// However, if your application uses floating point values, they will
// need to be converted somehow, according to your
// requirements. Likewise for `null` or `bool` values. When reading
// metadata from chain, be aware that integers may exceed the
// javascript numeric range, and may need special "bigint" parsing.
type Metadata map[uint]interface{}

// TODO add Metadata.ParseInto by re-encoding the result of Parse() as clean json,
// and parsing it into an arbitrary interface{}

// TODO add EncodeMetadataFrom(interface{}) by dumping to json, re-parsing to a map[...]...
// and passing the result to EncodeMetadata()

// TODO how to handle bools, floats, and nil values?
// TODO Is it necessary to handle all primitive types in EncodeMetadata? Or better error out?

func (meta *Metadata) Parse() (map[uint]interface{}, error) {
	if meta == nil || len(*meta) == 0 {
		return nil, nil
	}
	result := make(map[uint]interface{})
	for key, val := range *meta {
		parsedVal, err := parseMetaValue(strconv.Itoa(int(key)), val)
		if err != nil {
			return nil, err
		}
		result[key] = parsedVal
	}
	return result, nil
}

func str(val interface{}) string {
	// Limit length of error messages
	s := fmt.Sprintf("%T: %v", val, val)
	return fmt.Sprintf("%.25s", s)
}

func parseMetaValue(path string, rawVal interface{}) (interface{}, error) {
	val, ok := rawVal.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%v: unexpected type for metadata object: %s", path, str(rawVal))
	}
	if len(val) != 1 {
		return nil, fmt.Errorf("%v: unexpected length of metadata object (%v): %s", path, len(val), str(val))
	}
	for valType, actualVal := range val {
		// This loop will be entered only once
		switch valType {
		case MetadataTypeInt:
			intVal, ok := actualVal.(int) // TODO this does not match int64 when encoding
			if !ok {
				return nil, fmt.Errorf("%v: expected type int, but got: %s", path, str(actualVal))
			}
			return intVal, nil
		case MetadataTypeString:
			strVal, ok := actualVal.(string)
			if !ok {
				return nil, fmt.Errorf("%v: expected type string, but got: %s", path, str(actualVal))
			}
			return strVal, nil
		case MetadataTypeBytes:
			strVal, ok := actualVal.(string)
			if !ok {
				return nil, fmt.Errorf("%v: expected type string, but got: %s", path, str(actualVal))
			}
			return parseByteString(path, strVal)
		case MetadataTypeList:
			listVal, ok := actualVal.([]interface{})
			if !ok {
				return nil, fmt.Errorf("%v: expected type []interface{} for list, but got: %s", path, str(actualVal))
			}
			return parseList(path, listVal)
		case MetadataTypeMap:
			listVal, ok := actualVal.([]interface{})
			if !ok {
				return nil, fmt.Errorf("%v: expected type []interface{} for map, but got: %s", path, str(actualVal))
			}
			return parseMap(path, listVal)
		default:
			return nil, fmt.Errorf("%v: unknown metadata type '%v' for %s", path, valType, str(actualVal))
		}
	}
	return val, nil
}

func parseByteString(path string, rawStr string) ([]byte, error) {
	result := make([]byte, hex.DecodedLen(len(rawStr)))
	length, err := hex.Decode(result, []byte(rawStr))
	if err != nil {
		return nil, fmt.Errorf("%v: failed to decode hex-bytestring (%v): %v", path, err, rawStr)
	}
	return result[:length], err
}

func parseList(path string, rawList []interface{}) ([]interface{}, error) {
	result := make([]interface{}, len(rawList))
	for i, rawVal := range rawList {
		parsedVal, err := parseMetaValue(path+"/"+strconv.Itoa(i), rawVal)
		if err != nil {
			return nil, err
		}
		result[i] = parsedVal
	}
	return result, nil
}

func parseMap(path string, rawList []interface{}) (map[interface{}]interface{}, error) {
	result := make(map[interface{}]interface{}, len(rawList))
	for i, rawPair := range rawList {
		itemPath := path + "/" + strconv.Itoa(i)
		rawPairObj, ok := rawPair.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%v: expected map[string]interface{}, but got %s", itemPath, str(rawPair))
		}
		if len(rawPairObj) != 2 {
			return nil, fmt.Errorf("%v: expected length 2, but got %v: %s", itemPath, len(rawPairObj), str(rawPairObj))
		}
		rawKey, keyOk := rawPairObj[MetadataMapKey]
		if !keyOk {
			return nil, fmt.Errorf("%v: missing key '%v' in map-pair: %s", itemPath, MetadataMapKey, str(rawPairObj))
		}
		rawVal, valOk := rawPairObj[MetadataMapVal]
		if !valOk {
			return nil, fmt.Errorf("%v: missing key '%v' in map-pair: %s", itemPath, MetadataMapVal, str(rawPairObj))
		}
		parsedKey, err := parseMetaValue(itemPath+"["+MetadataMapKey+"]", rawKey)
		if err != nil {
			return nil, err
		}
		parsedVal, err := parseMetaValue(itemPath+"["+MetadataMapVal+"]", rawVal)
		if err != nil {
			return nil, err
		}

		_, collision := result[parsedKey]
		if collision {
			// TODO any better error handling?
			return nil, fmt.Errorf("%v: duplicate map key %s", itemPath, str(parsedKey))
		}
		result[parsedKey] = parsedVal
	}
	return result, nil
}

func EncodeMetadata(input map[uint]interface{}) (Metadata, error) {
	result := make(Metadata, len(input))
	for key, val := range input {
		encodedVal, err := encodeMetaValue(strconv.Itoa(int(key)), val)
		if err != nil {
			return nil, err
		}
		result[key] = encodedVal
	}
	return result, nil
}

func encodeMetaValue(path string, val interface{}) (interface{}, error) {
	var encodedType string
	var encodedVal interface{}

	switch typedVal := val.(type) {
	case string:
		encodedType = MetadataTypeString
		encodedVal = typedVal
		break

	case []byte:
		encodedType = MetadataTypeBytes
		encodedVal = typedVal
		break

	case int:
	case int8: // covers byte
	case int16:
	case int32: // covers rune
	case int64:
	case uint:
	case uint8:
	case uint16:
	case uint32:
	case uint64:
		encodedType = MetadataTypeInt
		encodedVal = int64(typedVal)
		break

	case bool: // TODO bool is not covered when parsing
		encodedType = MetadataTypeInt
		if typedVal {
			encodedVal = 1
		} else {
			encodedVal = 0
		}

	// Avoid unexpected behavior due to rounding
	case float32:
	case float64:
	case complex64:
	case complex128:
		return nil, fmt.Errorf("%v: encoding floating point and complex values unsupported (value: %v)", path, val)

	case []interface{}:
		encodedType = MetadataTypeList
		var err error
		encodedVal, err = encodeMetaList(path, typedVal)
		if err != nil {
			return nil, err
		}
		break

	case map[interface{}]interface{}:
		encodedType = MetadataTypeList
		var err error
		encodedVal, err = encodeMetaMap(path, typedVal)
		if err != nil {
			return nil, err
		}
		break

	default:
		return nil, fmt.Errorf("%v: cannot encode value of unexpected type: %s", path, str(val))
	}

	return map[string]interface{}{
		encodedType: encodedVal,
	}, nil
}

func encodeMetaList(path string, listVal []interface{}) (interface{}, error) {
	encodedList := make([]interface{}, len(listVal))
	for i, val := range listVal {
		encodedVal, err := encodeMetaValue(path+"/"+strconv.Itoa(i), val)
		if err != nil {
			return nil, err
		}
		encodedList[i] = encodedVal
	}
	return map[string]interface{}{
		MetadataTypeList: encodedList,
	}, nil
}

func encodeMetaMap(path string, mapVal map[interface{}]interface{}) (interface{}, error) {
	encodedMap := make([]interface{}, 0, len(mapVal))
	for key, val := range mapVal {
		itemPath := path + "/" + fmt.Sprintf("%v", key)
		encodedKey, err := encodeMetaValue(itemPath+"[key]", key)
		if err != nil {
			return nil, err
		}
		encodedVal, err := encodeMetaValue(itemPath+"[value]", val)
		if err != nil {
			return nil, err
		}

		encodedMap = append(encodedMap, map[string]interface{}{
			MetadataMapKey: encodedKey,
			MetadataMapVal: encodedVal,
		})
	}
	return map[string]interface{}{
		MetadataTypeMap: encodedMap,
	}, nil
}
