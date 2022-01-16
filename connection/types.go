package connection

type Params = map[string]interface{}

type Handler = func(request Request) (interface{}, error)

type Request struct {
	Id     string `json:"id"`
	Method string `json:"method"`
	Params Params `json:"params"`
}

type Response struct {
	Id     string      `json:"id"`
	Ok     bool        `json:"ok"`
	Result interface{} `json:"result"`
}
