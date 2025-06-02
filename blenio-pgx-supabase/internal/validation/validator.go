// validation/validator.go

package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Raw *validator.Validate

var regex = regexp.MustCompile(`[^.\[\]]+|\[\d+\]`)

func Validate(i interface{}) (map[string]string, error) {
	fmt.Println("✅ Validate function was called")
	err := Raw.Struct(i)
	if err == nil {
		return nil, nil
	}

	errs := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		nameSpaces := regex.FindAllString(e.StructNamespace(), -1)

		var jsonParts []string

		currentVal := reflect.ValueOf(i)
		if currentVal.Kind() == reflect.Ptr {
			currentVal = currentVal.Elem()
		}

		fmt.Printf("\n🔍 Validation Error: %s\n", e)
		fmt.Printf("Struct Namespace: %s\n", e.StructNamespace())
		fmt.Printf("Namespace Parts after split: %+v\n", nameSpaces)

		for idx := 1; idx < len(nameSpaces); idx++ {
			ns := nameSpaces[idx]
			fmt.Printf("\n📦 Processing Part [%d]: %q\n", idx, ns)

			// Handle map keys
			if currentVal.Kind() == reflect.Map {
				key := reflect.ValueOf(ns)
				value := currentVal.MapIndex(key)
				if value.IsValid() {
					fmt.Printf("→ Map key %q found. Moving into value\n", ns)
					currentVal = value
					continue
				} else {
					fmt.Printf("⚠️ Map key %q not found\n", ns)
					break
				}
			}

			// Handle array/slice index
			if strings.HasPrefix(ns, "[") {
				fmt.Printf("→ Handling array/slice index: %s\n", ns)
				jsonParts = append(jsonParts, ns)

				indexStr := strings.Trim(ns, "[]")
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					fmt.Printf("⚠️ Invalid index format: %s\n", indexStr)
					continue
				}

				if currentVal.Kind() == reflect.Slice || currentVal.Kind() == reflect.Array {
					if index < 0 || index >= currentVal.Len() {
						fmt.Printf("⚠️ Index out of range: %d for length %d\n", index, currentVal.Len())
						continue
					}

					element := currentVal.Index(index)
					if element.Kind() == reflect.Ptr && !element.IsNil() {
						currentVal = element.Elem()
					} else {
						currentVal = element
					}

					fmt.Printf("→ Moved into slice/array element at index %d\n", index)
				} else {
					fmt.Printf("⚠️ Cannot index into non-slice/array type: %s\n", currentVal.Kind())
				}

				continue
			}

			// Now process the actual struct field
			fieldVal := currentVal.FieldByName(ns)
			fmt.Printf("→ Struct Field Name: %q\n", ns)
			fmt.Printf("→ Is field valid? %t\n", fieldVal.IsValid())

			if !fieldVal.IsValid() {
				fmt.Printf("⚠️ Field %q not found. Using lowercase fallback.\n", ns)
				jsonParts = append(jsonParts, strings.ToLower(ns))
				continue
			}

			fieldType, ok := currentVal.Type().FieldByName(ns)
			fmt.Printf("→ Field Type Found? %t\n", ok)

			if !ok {
				fmt.Printf("⚠️ Field type not found for %q. Using lowercase fallback.\n", ns)
				jsonParts = append(jsonParts, strings.ToLower(ns))
				continue
			}

			jsonTag := fieldType.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				jsonTag = strings.ToLower(ns)
				fmt.Printf("→ No json tag found. Using lowercase: %q\n", jsonTag)
			} else {
				jsonTag = strings.Split(jsonTag, ",")[0]
				fmt.Printf("→ Found json tag: %q\n", jsonTag)
			}

			jsonParts = append(jsonParts, jsonTag)

			if strings.HasPrefix(ns, "[") {
				continue
			}
			// Update currentVal for next level if needed
			if fieldVal.Kind() == reflect.Ptr && fieldVal.Elem().IsValid() {
				fmt.Printf("→ Following pointer to nested struct\n")
				currentVal = fieldVal.Elem()
			} else if fieldVal.IsValid() {
				fmt.Printf("→ Moving into nested struct or value\n")
				currentVal = fieldVal
			}
		}

		jsonName := strings.Join(jsonParts, ".")
		fmt.Printf("✅ Final JSON Path: %q\n", jsonName)
		errs[jsonName] = e.Tag()
	}

	return errs, err
}

func init() {
	Raw = validator.New()

	// Register custom validators
	_RegisterCustomValidators()

	// Register custom translation functions if needed
	// _RegisterTranslations()
}
