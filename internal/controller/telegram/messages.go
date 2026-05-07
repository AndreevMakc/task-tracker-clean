package telegram

import (
	"fmt"
	"time"

	"task-tracker-clean/internal/entity"
)

func FormatTask(t entity.Task) string {
	return fmt.Sprintf("📋 <b>%s</b>\n"+
		"ID: %s\n"+
		"Status: %s\n"+
		"Created: %s\n"+
		"Updated: %s",
		t.Title,
		t.ID,
		formatStatus(t.Status),
		t.CreatedAt.Format(time.RFC822),
		t.UpdatedAt.Format(time.RFC822),
	)
}

func FormatTaskList(tasks []entity.Task) string {
	if len(tasks) == 0 {
		return "No tasks found."
	}

	msg := "📋 <b>Your Tasks</b>\n\n"
	for i, t := range tasks {
		emoji := statusEmoji(t.Status)
		msg += fmt.Sprintf("%d. %s %s\n   ID: %s\n",
			i+1, emoji, t.Title, t.ID)
	}
	return msg
}

func formatStatus(s entity.TaskStatus) string {
	switch s {
	case entity.TaskStatusToDo:
		return "📝 To Do"
	case entity.TaskStatusInProgress:
		return "🔄 In Progress"
	case entity.TaskStatusDone:
		return "✅ Done"
	case entity.TaskStatusTrashed:
		return "🗑️ Trashed"
	default:
		return string(s)
	}
}

func statusEmoji(s entity.TaskStatus) string {
	switch s {
	case entity.TaskStatusToDo:
		return "📝"
	case entity.TaskStatusInProgress:
		return "🔄"
	case entity.TaskStatusDone:
		return "✅"
	case entity.TaskStatusTrashed:
		return "🗑️"
	default:
		return "❓"
	}
}