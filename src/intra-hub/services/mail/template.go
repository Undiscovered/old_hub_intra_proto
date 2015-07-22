package mail

const (
	templateHeader = `Hello {{ .User.Login }},<br/><br/>`
	templateFooter = `<br/><br/>Innovation Hub`

	subjectUserCreated  = "Innovation Hub - Bienvenue {{ .User.Login }}"
	templateUserCreated = templateHeader + `Activate your account <a href="{{ .Link }}">here</a>.` + templateFooter

	subjectForgotPassword  = "Innovation Hub - Password reset"
	templateForgotPassword = templateHeader + `Your password has been reset. Set it again <a href="{{ .Link }}">here</a>.` + templateFooter
)
