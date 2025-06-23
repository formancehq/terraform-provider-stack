package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DynamicObjectValue struct {
	Attributes map[string]attr.Value
}

func NewDynamicObjectValue(attrs map[string]attr.Value) DynamicObjectValue {
	return DynamicObjectValue{
		Attributes: attrs,
	}
}

func (d DynamicObjectValue) Schema() map[string]attr.Type {
	return GetMapTypeFromAttrTypes(d.Attributes)
}

func (d DynamicObjectValue) Value() basetypes.ObjectValue {
	if d.Attributes == nil {
		return types.ObjectNull(d.Schema())
	}

	return types.ObjectValueMust(d.Schema(), d.Attributes)
}

func GetMapTypeFromAttrTypes(attrTypes map[string]attr.Value) map[string]attr.Type {
	result := make(map[string]attr.Type)

	for k, v := range attrTypes {
		result[k] = v.Type(context.Background())
	}

	return result
}

type DynamicTupleValue struct {
	Elements []attr.Value
}

func (d DynamicTupleValue) Schema() []attr.Type {
	return GetListTypeFromAttrTypes(d.Elements)
}

func (d DynamicTupleValue) Value() basetypes.TupleValue {
	if d.Elements == nil {
		return types.TupleNull(d.Schema())
	}

	return types.TupleValueMust(d.Schema(), d.Elements)
}

func NewDynamicTupleValue(elems []attr.Value) DynamicTupleValue {
	return DynamicTupleValue{
		Elements: elems,
	}
}

func GetListTypeFromAttrTypes(attrTypes []attr.Value) []attr.Type {
	result := make([]attr.Type, len(attrTypes))

	for i, v := range attrTypes {
		result[i] = v.Type(context.Background())
	}

	return result
}

// ConvertToAttrValues converts a map[string]any to a map[string]attr.Value.
// This is useful for converting dynamic data structures to Terraform attribute values.
func ConvertToAttrValues(input map[string]any) map[string]attr.Value {
	result := make(map[string]attr.Value)

	for k, v := range input {
		attrVal := convertAnyToAttrValue(v)
		result[k] = attrVal
	}

	return result
}

// convertAnyToAttrValue converts a value of any type to an attr.Value.
func convertAnyToAttrValue(v any) attr.Value {
	switch val := v.(type) {
	case string:
		return types.StringValue(val)
	case int64:
		return types.Int64Value(val)
	case float64:
		if val == float64(int64(val)) { // Check if it can be represented as int64
			return types.Int64Value(int64(val))
		}
		return types.Float64Value(val)
	case bool:
		return types.BoolValue(val)
	case nil:
		return types.DynamicNull()
	case map[string]any:
		innerMap := ConvertToAttrValues(val)
		return NewDynamicObjectValue(innerMap).Value()
	case []any:
		elems := make([]attr.Value, len(val))
		for i, item := range val {
			elems[i] = convertAnyToAttrValue(item)
		}
		return NewDynamicTupleValue(elems).Value()
	default:
		panic(fmt.Sprintf("unsupported type %T for value %v", v, v))
	}
}
