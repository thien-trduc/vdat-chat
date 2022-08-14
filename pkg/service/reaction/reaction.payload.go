package reaction

type Payload struct {
	IdArticle int64 `json:"idArticle"`
	Type      int8  `json:"type"`
}

func (p Payload) convertToModel() *Reaction {
	reaction := &Reaction{
		IdArticle: p.IdArticle,
		Type:      p.Type,
	}
	return reaction
}
