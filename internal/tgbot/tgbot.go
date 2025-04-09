package tgbot

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"gourbot/internal/config"

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
}

// NewTgBot initializes a new TgBot instance.
func NewTgBot(cfg *config.Config, logger *logrus.Logger) (*TgBot, error) {

	tgBot := &TgBot{
		config:   cfg,
		logger:   logger,
		chanQuit: make(chan struct{}, 1),
		commands: make(map[string]string),
	}
	// Initialize the Telegram bot
	b, err := bot.New(cfg.TGBotToken, bot.WithDefaultHandler(
		func(ctx context.Context, _ *bot.Bot, update *models.Update) {
			tgBot.DefaultHandler(ctx, update)
		}))
	if err != nil {
		return nil, err
	}
	tgBot.bot = b
	return tgBot, nil
}

// RegisterCommand registers a command with the bot.
func (tgBot *TgBot) RegisterCommand(command string, handler func(ctx context.Context, update *models.Update)) {
	if _, exists := tgBot.commands[command]; exists {
		return // Command already registered
	}

	handlerID := tgBot.bot.RegisterHandler(bot.HandlerTypeMessageText, command, bot.MatchTypeExact, bot.HandlerFunc(
		func(ctx context.Context, botInstance *bot.Bot, update *models.Update) {
			tgBot.wgWorkers.Add(1)
			defer tgBot.wgWorkers.Done()
			handler(ctx, update)
		},
	))
	tgBot.commands[command] = handlerID // Store handler ID as string
}

// Start begins the bot's operation.
func (tgBot *TgBot) Start(ctx context.Context) error {
	tgBot.context, tgBot.cancel = context.WithCancel(context.Background())
	// Register commands
	tgBot.RegisterCommand("ping", tgBot.CmdPing)
	tgBot.RegisterCommand("/list", tgBot.CmdList)
	tgBot.RegisterCommand("/stop", tgBot.CmdStop)

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

// Notify sends a message to the master user.
func (tgBot *TgBot) Notify(message string) {
	tgBot.wgWorkers.Add(1)
	defer tgBot.wgWorkers.Done()
	tgBot.logger.Info("Notify: " + message)
	_, err := tgBot.bot.SendMessage(tgBot.context, &bot.SendMessageParams{
		ChatID: tgBot.config.MasterUID,
		Text:   message,
	})
	if err != nil {
		tgBot.logger.Errorf("Notify failed: %v", err)
	} else {
		tgBot.logger.Info("Notify sent")
	}
}

func (tgBot *TgBot) Reply(ctx context.Context, update *models.Update, text string) (*models.Message, error) {
	return tgBot.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		ReplyParameters: &models.ReplyParameters{
			MessageID: update.Message.ID,
		},
		Text: text,
	})
}

// CmdStop handles the "/stop" command.
func (tgBot *TgBot) DefaultHandler(ctx context.Context, update *models.Update) {
	blob, err := json.Marshal(update)
	if err == nil {
		tgBot.logger.Infof("GOT::: %s", string(blob))
	}
	tgBot.Reply(ctx, update, "IDK what to do with your stuff")
}

// CmdPing handles the "ping" command.
func (tgBot *TgBot) CmdPing(ctx context.Context, update *models.Update) {
	tgBot.Reply(ctx, update, "pong")
}

// CmdList handles the "/list" command.
func (tgBot *TgBot) CmdList(ctx context.Context, update *models.Update) {
	// Generate the list of commands dynamically
	commandList := "Available commands:\n"
	for command := range tgBot.commands {
		commandList += "- " + command + "\n"
	}

	tgBot.Reply(ctx, update, commandList)
}

// CmdStop handles the "/stop" command.
func (tgBot *TgBot) CmdStop(ctx context.Context, update *models.Update) {
	if !tgBot.IsAllowed(update.Message.From.ID) {
		tgBot.Reply(ctx, update, "You are not authorized to stop the bot.")
		return
	}
	tgBot.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Bot is stopping...",
	})
	tgBot.Stop()
}
