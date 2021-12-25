package grpc_mapping

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/metadata"
)

func ExtractFromContext(ctx context.Context) (*Info, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("unable to extract context")
	}
	values := md.Get(InfoKey)
	if len(values) < 1 {
		return nil, fmt.Errorf("unable to extract context")
	}
	var info Info
	if err := json.Unmarshal([]byte(values[0]), &info); err != nil {
		return nil, fmt.Errorf("unable to unmarsharl: %w", err)
	}
	return &info, nil
}
