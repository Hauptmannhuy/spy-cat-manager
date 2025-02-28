package main

import (
	"fmt"
	"net/http"
)

type crudExecutor[T any] struct {
	handler interface{}
}

type getHandler[T any] struct {
	fn func(id int) (T, error)
	r  *http.Request
	id string
}

type createHandler[T any] struct {
	create func(data T) error
	r      *http.Request
}

type updateHandler[T any] struct {
	update func(int, T) (T, error)
	r      *http.Request
	id     string
}

type deleteHandler[T any] struct {
	delete func(int) error
	r      *http.Request
	id     string
}
type getAllHandler[T any] struct {
	getAll func() ([]T, error)
}

func (crudExec *crudExecutor[T]) create() error {
	hanler := crudExec.handler.(createHandler[T])
	r := hanler.r
	var data T
	err := decodeBody(r.Body, &data)

	if err != nil {
		return serverError{
			reason: err.Error(),
			code:   400,
		}
	}
	validate, ok := any(data).(createValidator)
	if ok {
		valid := validate.checkCreateValidity()
		if !valid {
			return serverError{
				reason: "Data is not valid",
				code:   422,
			}
		}
	}
	err = hanler.create(data)
	if err != nil {
		return dbError{
			reason: err.Error(),
			code:   http.StatusUnprocessableEntity,
		}
	}
	return nil
}

func (crudExec crudExecutor[T]) getAll() ([]byte, error) {
	getAllHandler, _ := crudExec.handler.(getAllHandler[T])
	data, err := getAllHandler.getAll()
	if err != nil {
		return nil, dbError{
			reason: err.Error(),
			code:   400,
		}
	}
	json, err := encodeBody(data)
	return json, err
}

func (crudExec *crudExecutor[T]) get() ([]byte, error) {
	getHandler, _ := crudExec.handler.(getHandler[T])

	fmt.Println("handler id", getHandler.id)
	id := getURLvar(getHandler.r, getHandler.id)
	data, err := getHandler.fn(id)
	if err != nil {
		return nil, dbError{
			reason: err.Error(),
			code:   400,
		}
	}
	json, err := encodeBody(data)
	if err != nil {
		return nil, serverError{
			reason: err.Error(),
			code:   500,
		}
	}
	return json, nil
}

func (crudExec crudExecutor[T]) update() ([]byte, error) {
	handler := crudExec.handler.(updateHandler[T])
	var data T
	err := decodeBody(handler.r.Body, &data)
	if err != nil {
		return nil, serverError{
			reason: err.Error(),
			code:   500,
		}
	}

	validate, ok := any(data).(updateValidator)
	if ok {
		valid := validate.checkUpdateValidity()
		if !valid {
			return nil, serverError{
				reason: "Data is not valid",
				code:   422,
			}
		}
	}

	id := getURLvar(handler.r, handler.id)
	data, err = handler.update(id, data)

	if err != nil {
		return nil, dbError{
			reason: err.Error(),
			code:   400,
		}
	}

	json, err := encodeBody(data)
	if err != nil {
		return nil, serverError{
			reason: err.Error(),
			code:   500,
		}
	}
	return json, nil
}

func (crudExec *crudExecutor[T]) delete() error {
	handler := crudExec.handler.(deleteHandler[T])
	id := getURLvar(handler.r, "id")
	err := handler.delete(id)
	if err != nil {
		return dbError{
			reason: err.Error(),
			code:   400,
		}
	}
	return nil
}
