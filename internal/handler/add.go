package handler

import "fmt"

func (h *Handler) Add() {
	fmt.Println("text panic")
	panic("implement me")
}
