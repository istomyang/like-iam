package json

import "github.com/vmihailenco/msgpack/v5"

var (
	// MPMarshal is exported by component-base/pkg/json package.
	MPMarshal = msgpack.Marshal
	// MPUnmarshal is exported by component-base/pkg/json package.
	MPUnmarshal = msgpack.Unmarshal
	// MPNewDecoder is exported by component-base/pkg/json package.
	MPNewDecoder = msgpack.NewDecoder
	// MPNewEncoder is exported by component-base/pkg/json package.
	MPNewEncoder = msgpack.NewEncoder
)
