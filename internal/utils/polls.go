package utils

import "realtime-poll/internal/models"

func CalculateTotalVotes(p *models.Poll) int64 {
	var total int64

	for _, opt := range p.Content.Options {
		total += int64(opt.Votes)
	}

	return total
}

func SanitizePollForViewer(poll *models.Poll) *models.Poll {

	clone := *poll

	clone.Content.Options = make([]models.Option, len(poll.Content.Options))

	for i, opt := range poll.Content.Options {
		clone.Content.Options[i] = opt
		clone.Content.Options[i].Votes = 0
	}

	clone.Meta.TotalVotes = 0

	return &clone
}

func Bool(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

func Int(v *int, def int) int {
	if v == nil {
		return def
	}
	return *v
}

func String(v *string, def string) string {
	if v == nil || *v == "" {
		return def
	}
	return *v
}
