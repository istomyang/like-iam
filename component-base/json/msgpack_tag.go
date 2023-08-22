//go:build msgpack

package json

import "github.com/vmihailenco/msgpack/v5"

var (
	// MGPMarshal is exported by component-base/pkg/json package.
	Marshal = msgpack.Marshal
	// MGPUnmarshal is exported by component-base/pkg/json package.
	Unmarshal = msgpack.Unmarshal
	// MGPNewDecoder is exported by component-base/pkg/json package.
	NewDecoder = msgpack.NewDecoder
	// MGPNewEncoder is exported by component-base/pkg/json package.
	NewEncoder = msgpack.NewEncoder
)
