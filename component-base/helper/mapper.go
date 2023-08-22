package helper

type MapStringAny map[string]any

func (m MapStringAny) Keys() []string {
	r := make([]string, len(m))
	for k, _ := range m {
		r = append(r, k)
	}
	return r
}
