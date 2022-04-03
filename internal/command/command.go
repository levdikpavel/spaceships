package command

import (
	"fmt"
	"log"
	"math"

	"modules/internal/vector"
)

type LogFunc func(string)

type LogCommand struct {
	command Command
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
	command Command
	attempt int
}

func (c RepeatCommand) Execute() error {
	return c.command.Execute()
}

func NewMoveCommand(m Movable) *MoveCommand {
	return &MoveCommand{
		m: m,
	}
}

type MoveCommand struct {
	m Movable
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

func NewRotateCommand(r Rotatable) *RotateCommand {
	return &RotateCommand{
		r: r,
	}
}

type RotateCommand struct {
	r Rotatable
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

func NewTurnVelocityCommand(object MovableRotatable) Command {
	return &TurnVelocityCommand{
		object: object,
	}
}

type TurnVelocityCommand struct {
	object MovableRotatable
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
	object FuelBurnable

	fuel        int
	consumption int
}

func NewCheckFuelCommand(object FuelBurnable) *CheckFuelCommand {
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

func NewBurnFuelCommand(object FuelBurnable) *BurnFuelCommand {
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
	commands []Command
}

func NewMacroCommand(commands ...Command) Command {
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

func NewMoveWithFuelCommand(object MovableWithFuel) Command {
	checkCommand := NewCheckFuelCommand(object)
	moveCommand := NewMoveCommand(object)
	burnCommand := NewBurnFuelCommand(object)
	result := NewMacroCommand(checkCommand, moveCommand, burnCommand)
	return result
}

func NewRotateWithVelocityCommand(object Rotatable) Command {
	rotateCommand := NewRotateCommand(object)

	movableRotatable, isMovableRotatable := object.(MovableRotatable)
	if isMovableRotatable {
		turnCommand := NewTurnVelocityCommand(movableRotatable)
		return NewMacroCommand(rotateCommand, turnCommand)
	}

	return rotateCommand
}
