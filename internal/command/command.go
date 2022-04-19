package command

import (
	"fmt"
	"log"
	"math"

	"modules/internal/core"
	"modules/internal/vector"
)

type LogFunc func(string)

type LogCommand struct {
	command core.Command
	err     error
	logFunc LogFunc
}

func StdLogFunc(message string) {
	log.Println(message)
}

func (c LogCommand) Execute() error {
	message := fmt.Sprintf("%s got error: '%s'", getType(c.command), c.err)
	c.logFunc(message)
	return nil
}

type RepeatCommand struct {
	command core.Command
	attempt int
}

func (c RepeatCommand) Execute() error {
	return c.command.Execute()
}

func NewMoveCommand(m core.Movable) *MoveCommand {
	return &MoveCommand{
		m: m,
	}
}

type MoveCommand struct {
	m core.Movable
}

func (m *MoveCommand) Execute() error {
	position, err := m.m.GetPosition()
	if err != nil {
		return err
	}

	velocity, err := m.m.GetVelocity()
	if err != nil {
		return err
	}

	err = m.m.SetPosition(vector.Add(position, velocity))
	if err != nil {
		return err
	}

	return nil
}

func NewRotateCommand(r core.Rotatable) *RotateCommand {
	return &RotateCommand{
		r: r,
	}
}

type RotateCommand struct {
	r core.Rotatable
}

func (r *RotateCommand) Execute() error {
	direction, err := r.r.GetDirection()
	if err != nil {
		return err
	}

	angularVelocity, err := r.r.GetAngularVelocity()
	if err != nil {
		return err
	}

	n, err := r.r.GetDirectionsNumber()
	if err != nil {
		return err
	}

	directionNew := (direction + angularVelocity) % n
	err = r.r.SetDirection(directionNew)
	if err != nil {
		return err
	}

	return nil
}

func NewTurnVelocityCommand(object core.MovableRotatable) core.Command {
	return &TurnVelocityCommand{
		object: object,
	}
}

type TurnVelocityCommand struct {
	object core.MovableRotatable
}

func (c *TurnVelocityCommand) Execute() error {
	velocity, err := c.object.GetVelocity()
	if err != nil {
		return err
	}

	if len(velocity) != 2 {
		return ErrUnsupportedDimension
	}

	direction, err := c.object.GetDirection()
	if err != nil {
		return err
	}

	angularVelocity, err := c.object.GetAngularVelocity()
	if err != nil {
		return err
	}

	n, err := c.object.GetDirectionsNumber()
	if err != nil {
		return err
	}

	directionNew := (direction + angularVelocity) % n
	alpha := float64(directionNew) / float64(n) * 2 * math.Pi

	vx, vy := velocity[0], velocity[1]
	v := math.Sqrt(float64(vx*vx) + float64(vy*vy))

	vx = int(v * math.Cos(alpha))
	vy = int(v * math.Sin(alpha))
	velocity = vector.New([]int{vx, vy})

	err = c.object.SetVelocity(velocity)
	if err != nil {
		return err
	}

	return nil
}

type CheckFuelCommand struct {
	object core.FuelBurnable

	fuel        int
	consumption int
}

func NewCheckFuelCommand(object core.FuelBurnable) *CheckFuelCommand {
	return &CheckFuelCommand{
		object: object,
	}
}

func (c *CheckFuelCommand) Execute() (err error) {
	c.fuel, err = c.object.GetFuel()
	if err != nil {
		return err
	}

	c.consumption, err = c.object.GetConsumption()
	if err != nil {
		return err
	}

	if c.fuel < c.consumption {
		return ErrNotEnoughFuel
	}

	return nil
}

type BurnFuelCommand struct {
	*CheckFuelCommand
}

func NewBurnFuelCommand(object core.FuelBurnable) *BurnFuelCommand {
	return &BurnFuelCommand{
		CheckFuelCommand: NewCheckFuelCommand(object),
	}
}

func (c BurnFuelCommand) Execute() error {
	err := c.CheckFuelCommand.Execute()
	if err != nil {
		return err
	}

	err = c.object.SetFuel(c.fuel - c.consumption)
	if err != nil {
		return err
	}

	return nil
}

type MacroCommand struct {
	commands []core.Command
}

func NewMacroCommand(commands ...core.Command) core.Command {
	result := &MacroCommand{}
	result.commands = append(result.commands, commands...)
	return result
}

func (c MacroCommand) Execute() (err error) {
	for _, command := range c.commands {
		err = command.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewMoveWithFuelCommand(object core.MovableWithFuel) core.Command {
	checkCommand := NewCheckFuelCommand(object)
	moveCommand := NewMoveCommand(object)
	burnCommand := NewBurnFuelCommand(object)
	result := NewMacroCommand(checkCommand, moveCommand, burnCommand)
	return result
}

func NewRotateWithVelocityCommand(object core.Rotatable) core.Command {
	rotateCommand := NewRotateCommand(object)

	movableRotatable, isMovableRotatable := object.(core.MovableRotatable)
	if isMovableRotatable {
		turnCommand := NewTurnVelocityCommand(movableRotatable)
		return NewMacroCommand(rotateCommand, turnCommand)
	}

	return rotateCommand
}
