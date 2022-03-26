package command

type CompositeErrorHandler struct {
	handlers map[string]ErrorHandler
}
