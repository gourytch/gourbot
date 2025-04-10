package tgbot

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gourbot/internal/config"
	"gourbot/internal/storage"
	"gourbot/internal/types"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

// TgBot represents the Telegram bot instance.
type TgBot struct {
	config    *config.Config
	logger    *logrus.Logger
	context   context.Context
	cancel    context.CancelFunc
	bot       *bot.Bot
	wgWorkers sync.WaitGroup
	chanQuit  chan struct{}
	commands  map[string]string // Store handler IDs as strings
	storage   *storage.Storage  // Add a new field for storage
}

// NewTgBot initializes a new TgBot instance.
func NewTgBot(cfg *config.Config, logger *logrus.Logger) (*TgBot, error) {
	tgBot := &TgBot{
		config:   cfg,
		logger:   logger,
		chanQuit: make(chan struct{}, 1),
		commands: make(map[string]string),
		storage:  storage.NewStorage(cfg), // Initialize the storage field
	}
	if err := tgBot.storage.Open(); err != nil {
		tgBot.logger.Fatalf("Failed to open storage: %v", err)
		return nil, err
	}
	master, err := tgBot.storage.GetTgUser(cfg.MasterUID)
	if err != nil {
		master = types.NewTgUser(cfg.MasterUID, "master", nil)
	}
	master.AddPermission(types.CanEverything)
	tgBot.storage.UpdateTgUser(master)

	opts := []bot.Option{
		bot.WithDefaultHandler(func(_ context.Context, _ *bot.Bot, update *models.Update) {
			tgBot.wgWorkers.Add(1)
			defer tgBot.wgWorkers.Done()
			tgBot.storage.AddTgRecord(false, update)
			if tgBot.Guard(update) {
				tgBot.DefaultHandler(update)
			} else {
				tgBot.logger.Infof("ignore user %d", update.Message.From.ID)
			}
		}),
	}
	// Initialize the Telegram bot
	b, err := bot.New(cfg.TGBotToken, opts...)
	if err != nil {
		return nil, err
	}
	tgBot.bot = b
	return tgBot, nil
}

// RegisterCommand registers a command with the bot.
func (tgBot *TgBot) RegisterCommand(command string, handler func(update *models.Update)) {
	if _, exists := tgBot.commands[command]; exists {
		return // Command already registered
	}

	handlerID := tgBot.bot.RegisterHandler(bot.HandlerTypeMessageText, command, bot.MatchTypeExact, bot.HandlerFunc(
		func(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
			tgBot.wgWorkers.Add(1)
			defer tgBot.wgWorkers.Done()
			tgBot.storage.AddTgRecord(false, update)
			if tgBot.Guard(update) {
				handler(update)
			} else {
				tgBot.logger.Infof("ignore user %d", update.Message.From.ID)
			}
		},
	))
	tgBot.commands[command] = handlerID // Store handler ID as string
}

// Start begins the bot's operation.
func (tgBot *TgBot) Start(ctx context.Context) error {
	// Register commands
	tgBot.RegisterCommand("ping", tgBot.CmdPing)
	tgBot.RegisterCommand("/list", tgBot.CmdList)
	tgBot.RegisterCommand("/stop", tgBot.CmdStop)

	tgBot.context, tgBot.cancel = context.WithCancel(context.Background())
	go func() {
		defer tgBot.cancel()
		tgBot.logger.Info("start proxy canceller ...")
		<-ctx.Done()
		tgBot.Notify("got signal from outer space")
		time.Sleep(100 * time.Millisecond)
		tgBot.Stop()
	}()

	// Start a goroutine to handle shutdown logic
	go func() {
		tgBot.logger.Info("start watching chanQuit ...")
		<-tgBot.chanQuit
		tgBot.logger.Info("got chanQuit")
		tgBot.Notify("bot got chanQuit signal")
		tgBot.wgWorkers.Wait() // Wait for all workers to finish
		tgBot.logger.Info("all workers finished - pull the trigger")
		tgBot.cancel() // Cancel the context
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		tgBot.Notify("bot started")
	}()
	// Start the bot
	tgBot.logger.Info("TgBot instance starting...")
	tgBot.bot.Start(tgBot.context)
	tgBot.logger.Info("TgBot instance finished...")
	return nil
}

// Stop gracefully stops the bot.
func (tgBot *TgBot) Stop() {
	tgBot.logger.Info("emit chanQuit signal...")
	select {
	case tgBot.chanQuit <- struct{}{}:
		tgBot.logger.Info("Shutdown signal sent to chanQuit")
	default:
		tgBot.logger.Warn("chanQuit already closed or full")
	}
}

// IsAllowed checks if the given ID is allowed to perform certain actions.
func (tgBot *TgBot) IsAllowed(id int64) bool {
	return id == tgBot.config.MasterUID
}

// Guard processes an incoming update, registers the user, and determines if further interaction is allowed.
func (tgBot *TgBot) Guard(update *models.Update) bool {
	tgBot.storage.AddTgRecord(false, update)

	// Extract user information from the update
	user := GetUserFromUpdate(update)
	if user == nil {
		tgBot.logger.Warn("Update does not contain user information")
		return false
	}
	info, _ := json.Marshal(user)

	username := user.Username
	if username == "" {
		username = user.FirstName + " " + user.LastName
	}

	// Check if the user exists in the system
	exists, err := tgBot.storage.TgUserExists(user.ID)
	if err != nil {
		tgBot.logger.Errorf("Failed to check user existence: %v", err)
		return false
	}

	if !exists {
		// Use the NewTgUser constructor to create a new TgUser instance
		tgUser := types.NewTgUser(user.ID, username, info)
		err = tgBot.storage.AddTgUser(tgUser)
		if err != nil {
			tgBot.logger.Errorf("Failed to add new user: %v", err)
			return false
		}

		// Notify the master about the new user
		message := "New user detected: " + username + ", ID: " + fmt.Sprint(user.ID) + ". To approve, use /approve_" + fmt.Sprint(user.ID)
		tgBot.Notify(message)
		return false
	}

	// Update existing user information
	tgUser, err := tgBot.storage.GetTgUser(user.ID)
	if err != nil {
		tgBot.logger.Errorf("Failed to retrieve user: %v", err)
		return false
	}

	tgUser.SeenAt = time.Now()
	tgUser.Name = username
	tgUser.Info = info
	err = tgBot.storage.UpdateTgUser(tgUser)
	if err != nil {
		tgBot.logger.Errorf("Failed to update user: %v", err)
		return false
	}

	// Check if the user has no permission to chat
	if !tgUser.HasPermission(types.CanChat) {
		return false
	}

	return true
}

func (tgBot *TgBot) SendMessage(smp *bot.SendMessageParams) (*models.Message, error) {
	tgBot.wgWorkers.Add(1)
	defer tgBot.wgWorkers.Done()
	msg, err := tgBot.bot.SendMessage(tgBot.context, smp)
	if err != nil {
		tgBot.logger.Errorf("SendMessage failed: %v", err)
	} else {
		tgBot.storage.AddTgRecord(true, msg)
	}
	return msg, err
}

// Notify sends a message to the master user.
func (tgBot *TgBot) Notify(message string) {
	tgBot.logger.Info("Notify: " + message)
	_, err := tgBot.SendMessage(&bot.SendMessageParams{
		ChatID: tgBot.config.MasterUID,
		Text:   message,
	})
	if err != nil {
		tgBot.logger.Errorf("Notify failed: %v", err)
	} else {
		tgBot.logger.Info("Notify sent")
	}
}

func (tgBot *TgBot) Reply(update *models.Update, text string) (*models.Message, error) {
	return tgBot.SendMessage(&bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
		Text: text,
	})
}

// CmdStop handles the "/stop" command.
func (tgBot *TgBot) DefaultHandler(update *models.Update) {
	blob, err := json.Marshal(update)
	if err == nil {
		tgBot.logger.Infof("GOT::: %s", string(blob))
	}
	tgBot.Reply(update, "IDK what to do with your stuff")
}

// CmdPing handles the "ping" command.
func (tgBot *TgBot) CmdPing(update *models.Update) {
	tgBot.Reply(update, "pong")
}

// CmdList handles the "/list" command.
func (tgBot *TgBot) CmdList(update *models.Update) {
	// Generate the list of commands dynamically
	commandList := "Available commands:\n"
	for command := range tgBot.commands {
		commandList += "- " + command + "\n"
	}

	tgBot.Reply(update, commandList)
}

// CmdStop handles the "/stop" command.
func (tgBot *TgBot) CmdStop(update *models.Update) {
	if !tgBot.IsAllowed(update.Message.From.ID) {
		tgBot.Reply(update, "You are not authorized to stop the bot.")
		return
	}
	tgBot.Notify("Bot is stopping...")
	tgBot.Stop()
}
