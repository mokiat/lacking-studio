package editor

type NavigatorPageChangedEvent struct {
	Editor *Model
}

type InspectorPageChangedEvent struct {
	Editor *Model
}

type SelectionChangedEvent struct {
	Editor *Model
}
