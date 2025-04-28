package crack

func HasRepeatingBlock(bytes []byte, bs int) bool {
	block_map := make(map[string]int)
	for i := 0; i < len(bytes); i += bs {
		block := string(bytes[i : i+bs])
		block_map[block]++
	}
	for _, v := range block_map {
		if v > 1 {
			return true
		}
	}
	return false
}
