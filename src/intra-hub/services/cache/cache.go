package cache

import "intra-hub/models"

var (
	Skills = make(map[string]*models.Skill)

	Themes = make(map[string]*models.Theme)

	Cities = make(map[string]*models.City)

	Promotions = make(map[string]*models.Promotion)
)

func SetSkills(skills []*models.Skill) {
	for _, s := range skills {
		Skills[s.Name] = s
	}
}

func SetThemes(themes []*models.Theme) {
	for _, s := range themes {
		Themes[s.Name] = s
	}
}

func SetCities(themes []*models.City) {
	for _, s := range themes {
		Cities[s.Name] = s
	}
}

func SetPromotions(themes []*models.Promotion) {
	for _, s := range themes {
		Promotions[s.Name] = s
	}
}
