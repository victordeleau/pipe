package pipe

type set map[string]struct{}

func (s *set) Slice() []string {
	list := make([]string, 0, len(*s))
	for key, _ := range *s {
		list = append(list, key)
	}
	return list
}
