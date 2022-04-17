package ioc

import (
	"fmt"
	"sync"

	"modules/internal/command"
)

type createCommand func(params ...interface{}) command.Command

type scope struct {
	registry sync.Map
	parent   *scope
}

func (s *scope) resolve(key string, params ...interface{}) command.Command {
	if s == nil {
		return nil
	}

	val, ok := s.registry.Load(key)
	if !ok {
		return s.parent.resolve(key, params...)
	}

	create := val.(createCommand)
	return create(params...)
}

type registerCommand struct {
	scope  *scope
	name   string
	create createCommand
}

func (c *registerCommand) Execute() error {
	c.scope.registry.Store(c.name, c.create)
	return nil
}

var scopes = func() scopesRegistry {
	defaultScope := createScope(nil)

	return scopesRegistry{
		currentScope: defaultScope,
		scopes:       make(map[string]*scope),
	}
}()

type scopesRegistry struct {
	currentScope *scope
	scopes       map[string]*scope
}

func createScope(parent *scope) *scope {
	result := &scope{parent: parent}

	result.registry.Store("IoC.Register", func(params ...interface{}) command.Command {
		name := params[0].(string)
		create := params[1].(createCommand)
		return &registerCommand{
			scope:  result,
			name:   name,
			create: create,
		}
	})

	result.registry.Store("Scopes.New", func(params ...interface{}) command.Command {
		return &setScopeCommand{
			scopeID:   params[0].(string),
			createNew: true,
		}
	})

	result.registry.Store("Scopes.Current", func(params ...interface{}) command.Command {
		return &setScopeCommand{
			scopeID:   params[0].(string),
			createNew: false,
		}
	})

	return result
}

type setScopeCommand struct {
	scopeID   string
	createNew bool
}

var ErrNoSuchScope = fmt.Errorf("no such scope")

func (c *setScopeCommand) Execute() error {
	newScope, ok := scopes.scopes[c.scopeID]
	if ok {
		scopes.currentScope = newScope
		return nil
	}

	if c.createNew {
		newScope := createScope(scopes.currentScope)
		scopes.scopes[c.scopeID] = newScope
		scopes.currentScope = newScope
		return nil
	}

	return ErrNoSuchScope
}

func Resolve(key string, params ...interface{}) command.Command {
	return scopes.currentScope.resolve(key, params...)
}
