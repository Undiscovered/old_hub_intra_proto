package mail

const (
	templateHeader = `Hello {{ .User.Login }},<br/><br/>`
	templateFooter = `<br/><br/>Innovation Hub`

	subjectUserCreated  = "Bienvenue {{ .User.Login }}"
	templateUserCreated = templateHeader + `Activate your account <a href="{{ .Link }}">here</a>.` + templateFooter
)
