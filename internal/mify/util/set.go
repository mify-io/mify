package util

func StringSetAppend(set []string, vals ...string) []string {
	mp := make(map[string]struct{}, len(set))
	for _, val := range set {
		mp[val] = struct{}{}
	}
	for _, v := range vals {
		if _, ok := mp[v]; !ok {
			mp[v] = struct{}{}
			set = append(set, v)
		}
	}
	return set
}
