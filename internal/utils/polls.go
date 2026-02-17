package utils

import "realtime-poll/internal/models"

func CalculateTotalVotes(p *models.Poll) int64 {
	var total int64

	for _, opt := range p.Content.Options {
		total += int64(opt.Votes)
	}

	return total
}
