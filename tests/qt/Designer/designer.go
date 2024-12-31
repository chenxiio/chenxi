package main

// type CustomDragWidget struct {
// 	widgets.QWidget
// }

// func NewCustomDragWidget(parent widgets.QWidget_ITF) *CustomDragWidget {
// 	widget := &CustomDragWidget{}
// 	widget.SetAcceptDrops(true)
// 	widget.ConnectDragEnterEvent(widget.dragEnterEvent)
// 	widget.ConnectDropEvent(widget.dropEvent)
// 	widget.SetWindowTitle("Custom Drag Widget")
// 	widget.Resize2(300, 200)
// 	widget.Show()
// 	return widget
// }
// func (widget *CustomDragWidget) dragEnterEvent(event *widgets.QDragEnterEvent) {
// 	if event.MimeData().HasFormat("text/plain") {
// 		event.AcceptProposedAction()
// 	}
// }
// func (widget *CustomDragWidget) dropEvent(event *widgets.QDropEvent) {
// 	mimeType := "text/plain"
// 	data := event.MimeData().Data(mimeType)
// 	text := core.NewQByteArrayFromPointer(data.Data()).Data()
// 	widget.SetText(text)
// }
// func main() {
// 	widgets.NewQApplication(len(core.Qt__AA_X11InitThreads), []*core.QByteArray{})
// 	NewCustomDragWidget(nil)
// 	widgets.QApplication_Exec()
// }
