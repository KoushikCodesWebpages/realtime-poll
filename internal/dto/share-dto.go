package dto

type SharedViewer struct {
	AlreadyVoted   bool   `json:"already_voted"`
	CanVote        bool   `json:"can_vote"`
	CanViewResults bool   `json:"can_view_results"`
	Changed        bool   `json:"changed"`
	SelectedOption string `json:"selected_option"`
	Started        bool   `json:"started"`
	Ended          bool   `json:"ended"`
}
