package adapters

import (
	"modules/internal/core"
	"modules/internal/ioc"
	"modules/internal/vector"
)

type MovableAdapter struct{
   obj core.UObject
}

func NewMovableAdapter(obj core.UObject) *MovableAdapter {
	return &MovableAdapter{
		obj: obj,
	}
}

func (a *MovableAdapter) GetPosition() (vector.Vector, error) {
	return ioc.Resolve("Operations.Movable:Position.get", a.obj).(vector.Vector), nil
}

func (a *MovableAdapter) GetVelocity() (vector.Vector, error) {
	return ioc.Resolve("Operations.Movable:Velocity.get", a.obj).(vector.Vector), nil
}

func (a *MovableAdapter) SetPosition(newValue vector.Vector) (error) {
	return ioc.Resolve("Operations.Movable:Position.set", a.obj, newValue).(core.Command).Execute()
}

type RotatableAdapter struct{
   obj core.UObject
}

func NewRotatableAdapter(obj core.UObject) *RotatableAdapter {
	return &RotatableAdapter{
		obj: obj,
	}
}

func (a *RotatableAdapter) GetDirection() (int, error) {
	return ioc.Resolve("Operations.Rotatable:Direction.get", a.obj).(int), nil
}

func (a *RotatableAdapter) GetAngularVelocity() (int, error) {
	return ioc.Resolve("Operations.Rotatable:AngularVelocity.get", a.obj).(int), nil
}

func (a *RotatableAdapter) SetDirection(newValue int) (error) {
	return ioc.Resolve("Operations.Rotatable:Direction.set", a.obj, newValue).(core.Command).Execute()
}

func (a *RotatableAdapter) GetDirectionsNumber() (int, error) {
	return ioc.Resolve("Operations.Rotatable:DirectionsNumber.get", a.obj).(int), nil
}


func init() {
	_ = ioc.Resolve("IoC.Register", "Adapter",
		func(params ...interface{}) interface{} {
			adapterType := params[0].(string)
			switch adapterType {
			case "core.Movable":
				return NewMovableAdapter(params[1].(core.UObject))
			case "core.Rotatable":
				return NewRotatableAdapter(params[1].(core.UObject))
			default:
				panic("unknown adapter type" + adapterType)
			}
	}).(core.Command).Execute()
}
