package controllers

import "net/http"

type Args struct {
	Who string
}

type Reply struct {
	Message string
}

type HelloService struct {
}

func (h *HelloService) Say(r *http.Request, a *Args, reply *Reply) error {

	reply.Message = "hello," + a.Who + " !"
	return nil
}


