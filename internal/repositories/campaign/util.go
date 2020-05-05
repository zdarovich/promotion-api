package campaign

func GetIds(c []Campaign) []int {
	var res []int
	for _, el := range c {
		res = append(res, el.ID)
	}
	return res
}
