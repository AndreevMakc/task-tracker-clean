package request

type CreateTask struct {
	Title string `json:"title"`
}

type UpdateTask struct {
	Title  *string `json:"title,omitempty"`
	Status *string `json:"status,omitempty"`
}
