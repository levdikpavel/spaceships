package ioc

import (
	"fmt"
	"sync"

	"github.com/timandy/routine"

	"modules/internal/command"
)

type createCommand func(params ...interface{}) command.Command

type scope struct {
	commandRegistry sync.Map
	parent          *scope
}

func (s *scope) resolve(key string, params ...interface{}) command.Command {
	if s == nil {
		return nil
	}

	val, ok := s.commandRegistry.Load(key)
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
	c.scope.commandRegistry.Store(c.name, c.create)
	return nil
}

var scopes = func() scopesRegistry {
	defaultScope := createScope(nil)

	return scopesRegistry{
		defaultScope: defaultScope,
	}
}()

type scopesRegistry struct {
	defaultScope *scope
	scopesByName sync.Map
	scopesByGID  sync.Map
}

func createScope(parent *scope) *scope {
	result := &scope{parent: parent}

	result.commandRegistry.Store("IoC.Register", func(params ...interface{}) command.Command {
		name := params[0].(string)
		create := params[1].(createCommand)
		return &registerCommand{
			scope:  result,
			name:   name,
			create: create,
		}
	})

	return result
}

type setScopeCommand struct {
	goroutineID int64
	scopeName   string
	createNew   bool
}

var ErrNoSuchScope = fmt.Errorf("no such scope")

func (c *setScopeCommand) Execute() error {
	newScope, ok := scopes.scopesByName.Load(c.scopeName)
	if ok {
		scopes.scopesByGID.Store(c.goroutineID, newScope)
		return nil
	}

	if c.createNew {
		newScope := createScope(getCurrentScope(c.goroutineID))
		scopes.scopesByName.Store(c.scopeName, newScope)
		scopes.scopesByGID.Store(c.goroutineID, newScope)
		return nil
	}

	return ErrNoSuchScope
}

func getCurrentScope(gid int64) *scope {
	currentScope, ok := scopes.scopesByGID.Load(gid)
	if ok {
		return currentScope.(*scope)
	}

	return scopes.defaultScope
}

func Resolve(key string, params ...interface{}) command.Command {
	gid := routine.Goid()

	switch key {
	case "Scopes.New":
		return &setScopeCommand{
			goroutineID: gid,
			scopeName:   params[0].(string),
			createNew:   true,
		}
	case "Scopes.Current":
		return &setScopeCommand{
			goroutineID: gid,
			scopeName:   params[0].(string),
			createNew:   false,
		}
	default:
		currentScope := getCurrentScope(gid)
		return currentScope.resolve(key, params...)
	}
}
