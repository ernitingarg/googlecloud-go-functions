package entity

import "time"

type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

type FirestoreValue struct {
	CreateTime time.Time              `json:"createTime"`
	Fields     map[string]interface{} `json:"fields"`
	Name       string                 `json:"name"`
	UpdateTime time.Time              `json:"updateTime"`
}

type IntegerValue struct {
	IntegerValue string `json:"integerValue"`
}
type StringValue struct {
	StringValue string `json:"stringValue"`
}

type BooleanValue struct {
	BooleanValue bool `json:"booleanValue"`
}
type DoubleValue struct {
	DoubleValue float64 `json:"doubleValue"`
}
type TimestampValue struct {
	TimestampValue time.Time `json:"timestampValue"`
}
type BytesValue struct {
	BytesValue []byte `json:"bytesValue"`
}

type ArrayValue struct {
	Values []Value `json:"values"`
}

type Value struct {
	MapValue     MapValue     `json:"mapValue,omitempty"`
	StringValue  StringValue  `json:"stringValue,omitempty"`
	IntegerValue IntegerValue `json:"integerValue,omitempty"`
	ArrayValue   ArrayValue   `json:"arrayValue,omitempty"`
	BooleanValue BooleanValue `json:"booleanValue,omitempty"`
	DoubleValue  DoubleValue  `json:"doubleValue,omitempty"`
}

type MapValue struct {
	Fields interface{} `json:"fields"`
}
