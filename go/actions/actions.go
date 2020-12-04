package actions

// ActionToken define a uma estrutura de dados de um ação com um Token de verificação
type ActionToken struct {
	Token string `json:"token"`
}

// Valid verifica se a estrutura da ação é válida.
// As regras para uma estrutura ser válida são:
// - Um token não vazio
func (at *ActionToken) Valid() bool {
	if at.Token == "" {
		return false
	}

	return true
}
