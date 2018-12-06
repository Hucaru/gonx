package gonx

// Skill data from nx
type Skill struct {
	MaxLevel byte
	Mastery  []int16
}

// ExtractSkills from parsed nx
func ExtractSkills(nodes []Node, textLookup []string) map[int32]Skill {
	skills := make(map[int32]Skill)

	return skills
}
