package chenxi

type IoEventTypes byte

const (
	DEVICE_STATUS           = "device-status"
	DEVICE_PROPERTY         = "device-property"
	DEVICE_VALUES           = "device-values"
	DEVICE_BROWSE           = "device-browse"
	DEVICE_NODE_ATTRIBUTE   = "device-node-attribute"
	DEVICE_WEBAPI_REQUEST   = "device-webapi-request"
	DEVICE_TAGS_REQUEST     = "device-tags-request"
	DEVICE_TAGS_SUBSCRIBE   = "device-tags-subscribe"
	DEVICE_TAGS_UNSUBSCRIBE = "device-tags-unsubscribe"
	DAQ_QUERY               = "daq-query"
	DAQ_RESULT              = "daq-result"
	DAQ_ERROR               = "daq-error"
	ALARMS_STATUS           = "alarms-status"
	HOST_INTERFACES         = "host-interfaces"
	SCRIPT_CONSOLE          = "script-console"
	SCRIPT_COMMAND          = "script-command"
)

// var (
// 	socketioServer *socketio.SocketioServer
// 	once           sync.Once
// )

// func initSocketIOServer() {
// 	// 初始化 socketioServer 代码
// 	socketioServer = socketio.NewServer()
// 	mux1.Handle("/socket.io/", socketioServer)
// 	//mux1.Handle("/", http.FileServer(http.Dir("../asset")))
// }
// func GetSocketIOServer() *socketio.SocketioServer {
// 	once.Do(initSocketIOServer)
// 	return socketioServer
// }
