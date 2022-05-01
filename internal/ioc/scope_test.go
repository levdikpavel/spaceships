package ioc

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"

	"modules/internal/core"
	"modules/internal/mock"
)

func TestScope(t *testing.T) {
	suite.Run(t, new(ScopeTestSuite))
}

type ScopeTestSuite struct {
	suite.Suite

	err error
}

func (s *ScopeTestSuite) SetupTest() {
	s.err = fmt.Errorf("some error")
}

func (s *ScopeTestSuite) TestRegister() {
	a := 10
	commandMock := mock.CommandMock{}
	commandMock.On("Execute").Return(func() error {
		a = 20
		return s.err
	})

	err := Resolve("IoC.Register", "my_command",
		func(params ...interface{}) interface{} {
			return &commandMock
		}).(core.Command).Execute()
	s.Require().NoError(err)

	err = Resolve("my_command").(core.Command).Execute()
	s.Require().ErrorIs(err, s.err)
	s.Require().Equal(20, a)

	commandMock.AssertExpectations(s.T())
}

func (s *ScopeTestSuite) TestScopes() {
	err := Resolve("Scopes.New", "scope1").(core.Command).Execute()
	s.Require().NoError(err)

	a := 10
	err = Resolve("IoC.Register", "my_command",
		func(params ...interface{}) interface{} {
			commandMock := mock.CommandMock{}
			commandMock.On("Execute").Return(func() error {
				a = 20
				return nil
			})
			return &commandMock
		}).(core.Command).Execute()
	s.Require().NoError(err)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := Resolve("Scopes.New", "scope2").(core.Command).Execute()
		s.Require().NoError(err)

		err = Resolve("IoC.Register", "my_command",
			func(params ...interface{}) interface{} {
				commandMock := mock.CommandMock{}
				commandMock.On("Execute").Return(func() error {
					a = 30
					return nil
				})
				return &commandMock
			}).(core.Command).Execute()
		s.Require().NoError(err)

		err = Resolve("Scopes.Current", "scope1").(core.Command).Execute()
		s.Require().NoError(err)

		err = Resolve("my_command").(core.Command).Execute()
		s.Require().NoError(err)
	}()

	wg.Wait()
	s.Require().Equal(20, a)

	err = Resolve("my_command").(core.Command).Execute()
	s.Require().NoError(err)
	s.Require().Equal(20, a)

	err = Resolve("Scopes.Current", "scope2").(core.Command).Execute()
	s.Require().NoError(err)

	err = Resolve("my_command").(core.Command).Execute()
	s.Require().NoError(err)
	s.Require().Equal(30, a)
}
