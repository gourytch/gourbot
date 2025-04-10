package types

// Update represents an update received from Telegram.
type Update struct {
	Message            *Message
	CallbackQuery      *CallbackQuery
	InlineQuery        *InlineQuery
	ChosenInlineResult *ChosenInlineResult
}

// User represents a Telegram user.
type User struct {
	ID        int64
	Username  string
	FirstName string
	LastName  string
}

// Message represents a Telegram message.
type Message struct {
	ID   int64
	Chat *Chat
	From *User
}

// Chat represents a Telegram chat.
type Chat struct {
	ID int64
}

// CallbackQuery represents a Telegram callback query.
type CallbackQuery struct {
	From *User
}

// InlineQuery represents a Telegram inline query.
type InlineQuery struct {
	From *User
}

// ChosenInlineResult represents a Telegram chosen inline result.
type ChosenInlineResult struct {
	From *User
}
