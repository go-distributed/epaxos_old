package replica

type dependencies []InstanceIdType

// union unions the dep into the receiver
func (d dependencies) union(toUnion dependencies) bool {
	same := true
	for i := range d {
		if len(toUnion) < i {
			break
		}

		if d[i] != toUnion[i] {
			same = false
			if d[i] < toUnion[i] {
				d[i] = toUnion[i]
			}
		}
	}
	return same
}
