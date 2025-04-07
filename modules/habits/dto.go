package habits

// TODO: Validation!
type createDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Remind      bool    `json:"remind"`
}
