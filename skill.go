package gonx

// PlayerSkill data from nx
type PlayerSkill struct {
	MaxLevel byte
	Mastery  []int16
}

// MobSkill data from nx
type MobSkill struct {
}

// ExtractSkills from parsed nx
func ExtractSkills(nodes []Node, textLookup []string) (map[int32][]PlayerSkill, map[int32][]MobSkill) {
	playerSkills := make(map[int32][]PlayerSkill)
	mobSkills := make(map[int32][]MobSkill)

	return playerSkills, mobSkills
}
