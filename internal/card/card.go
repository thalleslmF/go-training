package card


type ClientCardResponse struct {
	Id string `json:"id"`
	EmitidoEm string `json:"emitidoEm"`
	Titular string `json:"titular"`
	Limite int32 `json:"limite"`
	IdProposta string `json:"idProposta"`
}
