package mail

const (
	templateHeader = `Bonjour {{ .User.Login }},<br/><br/>`
	templateFooter = `<br/><br/>Innovation Hub`

	subjectUserCreated  = "Innovation Hub - Bienvenue {{ .User.Login }}"
	templateUserCreated = templateHeader + `Set ton mot de passe <a href="{{ .Link }}">ici</a>.` + templateFooter

	subjectForgotPassword  = "Innovation Hub - Mot de passe reset"
	templateForgotPassword = templateHeader + `Votre mot de passe a été reset. Réassigner le <a href="{{ .Link }}">ici</a>.` + templateFooter
)
