package dto

type APIResponse struct {
	Data *CEP
	Api  string
}

type CEP struct {
	Cep         string `json:"cep,omitempty"`
	Logradouro  string `json:"logradouro,omitempty"`
	Complemento string `json:"complemento,omitempty"`
	Unidade     string `json:"unidade,omitempty"`
	Bairro      string `json:"bairro,omitempty"`
	Rua         string `json:"rua,omitempty"`
	Localidade  string `json:"localidade,omitempty"`
	Uf          string `json:"uf,omitempty"`
	Cidade      string `json:"cidade,omitempty"`
	Estado      string `json:"estado,omitempty"`
	Regiao      string `json:"regiao,omitempty"`
	Ibge        string `json:"ibge,omitempty"`
	Gia         string `json:"gia,omitempty"`
	Ddd         string `json:"ddd,omitempty"`
	Siafi       string `json:"siafi,omitempty"`
	Servico     string `json:"servico,omitempty"`
}
