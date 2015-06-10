package mail

const (
	templateHeader      = `Hello .User.Login,

	`
    templateFooter = `

    Innovation Hub`

	templateUserCreated = templateHeader + `Activate your account <a href="{{ .Link }}">here</a>.` + templateFooter
)
