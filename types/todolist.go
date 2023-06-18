package types

type ToDoList []ListItem

type ListItem struct {
	Id           int64  `json:"id"`
	Taskname     string `json:"taskname"`
	Description  string `json:"description"`
	CreationDate string `json:"creationDate"`
	Done         bool   `json:"done"`
}
