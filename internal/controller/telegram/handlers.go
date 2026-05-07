package telegram

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"task-tracker-clean/internal/entity"
	"task-tracker-clean/internal/usecase"
)

type TaskHandler struct {
	uc usecase.TaskUsecase
}

func NewTaskHandler(uc usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{uc: uc}
}

func (h *TaskHandler) RegisterBotHandler(bh *th.BotHandler) {
	bh.Handle(h.handleStart, th.CommandEqual("start"))
	bh.Handle(h.handleHelp, th.CommandEqual("help"))
	bh.Handle(h.handleCreate, th.CommandEqual("create"))
	bh.Handle(h.handleList, th.CommandEqual("list"))
	bh.Handle(h.handleGet, th.CommandEqual("get"))
	bh.Handle(h.handleDone, th.CommandEqual("done"))
	bh.Handle(h.handleTodo, th.CommandEqual("todo"))
	bh.Handle(h.handleInProgress, th.CommandEqual("in_progress"))
	bh.Handle(h.handleTrash, th.CommandEqual("trash"))
	bh.Handle(h.handleUnknown, th.AnyCommand())
}

func (h *TaskHandler) handleStart(ctx *th.Context, update telego.Update) error {
	msg := "👋 Welcome to Task Tracker!\n\n" +
		"Available commands:\n" +
		"/create <title> - Create a new task\n" +
		"/list - List all tasks\n" +
		"/list <status> - List tasks by status\n" +
		"/get <id> - Get task details\n" +
		"/done <id> - Mark task as done\n" +
		"/todo <id> - Mark task as todo\n" +
		"/in_progress <id> - Mark task as in progress\n" +
		"/trash <id> - Move task to trash\n" +
		"/help - Show this help"
	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID), msg,
	).WithParseMode(telego.ModeHTML))
	return nil
}

func (h *TaskHandler) handleHelp(ctx *th.Context, update telego.Update) error {
	msg := "📖 <b>Available Commands</b>\n\n" +
		"/create <title> - Create a new task\n" +
		"/list - List all tasks\n" +
		"/list <status> - Filter by: todo, in_progress, done\n" +
		"/get <id> - Get task details\n" +
		"/done <id> - Mark done\n" +
		"/todo <id> - Mark todo\n" +
		"/in_progress <id> - Mark in progress\n" +
		"/trash <id> - Move to trash"
	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID), msg,
	).WithParseMode(telego.ModeHTML))
	return nil
}

func (h *TaskHandler) handleCreate(ctx *th.Context, update telego.Update) error {
	args := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/create"))
	if args == "" {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID), "Usage: /create <task title>",
		))
		return nil
	}

	task, err := h.uc.CreateTask(ctx.Context(), args)
	if err != nil {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID), fmt.Sprintf("Error: %v", err),
		))
		return nil
	}

	_, _ = ctx.Bot().SendMessage(ctx, tu.Messagef(
		tu.ID(update.Message.Chat.ID),
		"✅ Task created: %s\nID: %s", task.Title, task.ID,
	))
	return nil
}

func (h *TaskHandler) handleList(ctx *th.Context, update telego.Update) error {
	filter := usecase.TaskFilter{}

	args := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/list"))
	if args != "" {
		status := entity.TaskStatus(args)
		if status.Valid() {
			filter.Status = &status
		} else {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
				tu.ID(update.Message.Chat.ID),
				"Invalid status. Use: todo, in_progress, done",
			))
			return nil
		}
	}

	tasks, err := h.uc.ListTasks(ctx.Context(), filter)
	if err != nil {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Messagef(
			tu.ID(update.Message.Chat.ID), "Error: %v", err,
		))
		return nil
	}

	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID), FormatTaskList(tasks),
	).WithParseMode(telego.ModeHTML))
	return nil
}

func (h *TaskHandler) handleGet(ctx *th.Context, update telego.Update) error {
	idStr := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/get"))
	if idStr == "" {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID), "Usage: /get <task_id>",
		))
		return nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID), "Invalid task ID",
		))
		return nil
	}

	task, err := h.uc.GetTask(ctx.Context(), id)
	if err != nil {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID), "Task not found",
		))
		return nil
	}

	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID), FormatTask(*task),
	).WithParseMode(telego.ModeHTML))
	return nil
}

func (h *TaskHandler) updateStatus(ctx *th.Context, update telego.Update, command string, status entity.TaskStatus) error {
	idStr := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/"+command))
	if idStr == "" {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID), fmt.Sprintf("Usage: /%s <task_id>", command),
		))
		return nil
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID), "Invalid task ID",
		))
		return nil
	}

	_, err = h.uc.UpdateTask(ctx.Context(), id, nil, &status)
	if err != nil {
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID), fmt.Sprintf("Error: %v", err),
		))
		return nil
	}

	_, _ = ctx.Bot().SendMessage(ctx, tu.Messagef(
		tu.ID(update.Message.Chat.ID),
		"✅ Task updated to %s", status,
	))
	return nil
}

func (h *TaskHandler) handleDone(ctx *th.Context, update telego.Update) error {
	return h.updateStatus(ctx, update, "done", entity.TaskStatusDone)
}

func (h *TaskHandler) handleTodo(ctx *th.Context, update telego.Update) error {
	return h.updateStatus(ctx, update, "todo", entity.TaskStatusToDo)
}

func (h *TaskHandler) handleInProgress(ctx *th.Context, update telego.Update) error {
	return h.updateStatus(ctx, update, "in_progress", entity.TaskStatusInProgress)
}

func (h *TaskHandler) handleTrash(ctx *th.Context, update telego.Update) error {
	return h.updateStatus(ctx, update, "trash", entity.TaskStatusTrashed)
}

func (h *TaskHandler) handleUnknown(ctx *th.Context, update telego.Update) error {
	_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
		tu.ID(update.Message.Chat.ID), "Unknown command. Use /help for list of commands",
	))
	return nil
}
