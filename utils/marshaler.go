package utils

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoMarshaler struct{}

func (m *ProtoMarshaler) Marshal(message protoreflect.ProtoMessage) ([]byte, error) {
	b, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal message: %w", err)
	}
	return b, nil
}

func (m *ProtoMarshaler) Unmarshal(b []byte, message protoreflect.ProtoMessage) error {
	if err := proto.Unmarshal(b, message); err != nil {
		return fmt.Errorf("unable to unmarshal message: %w", err)
	}
	return nil
}
