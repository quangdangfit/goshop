package notification

// Settings is the projection of pkg/config that the notifier factory needs. Defined here so
// the package doesn't depend on pkg/config (which would create a cycle for tests).
type Settings struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	EmailFrom    string
	// Prefs filters per-user delivery. Nil falls back to AlwaysOnPreferences.
	Prefs PreferenceChecker
	// DLQ receives notifications that exhausted retries. Nil disables retry+DLQ wrapping.
	DLQ DeadLetterSink
}

// BuildDefault returns a Notifier appropriate for the runtime config. SMTP is layered in
// only when both SMTPHost and EmailFrom are set; otherwise the notifier is logger-only so
// local development never tries to dial a non-existent mail server.
func BuildDefault(s Settings) Notifier {
	notifiers := []Notifier{NewLoggerNotifier()}
	if s.SMTPHost != "" && s.EmailFrom != "" {
		sender := NewSMTPSender(SMTPConfig{
			Host:     s.SMTPHost,
			Port:     s.SMTPPort,
			User:     s.SMTPUser,
			Password: s.SMTPPassword,
			From:     s.EmailFrom,
		})
		prefs := s.Prefs
		if prefs == nil {
			prefs = AlwaysOnPreferences{}
		}
		notifiers = append(notifiers, NewEmailNotifier(sender, prefs))
	}
	var n Notifier
	if len(notifiers) == 1 {
		n = notifiers[0]
	} else {
		n = NewMultiNotifier(notifiers...)
	}
	if s.DLQ != nil {
		n = NewRetryingNotifier(n, RetryConfig{}, s.DLQ)
	}
	return n
}
