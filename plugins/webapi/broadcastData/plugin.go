package broadcastData

import (
	"net/http"

	"github.com/iotaledger/goshimmer/packages/binary/messagelayer/payload"
	"github.com/iotaledger/goshimmer/plugins/messagelayer"
	"github.com/iotaledger/goshimmer/plugins/webapi"
	"github.com/iotaledger/hive.go/logger"
	"github.com/iotaledger/hive.go/marshalutil"
	"github.com/iotaledger/hive.go/node"
	"github.com/labstack/echo"
)

var PLUGIN = node.NewPlugin("WebAPI broadcastData Endpoint", node.Enabled, configure)
var log *logger.Logger

func configure(plugin *node.Plugin) {
	log = logger.NewLogger("API-broadcastData")
	webapi.Server.POST("broadcastData", broadcastData)
}

// broadcastData creates a message of the given payload and
// broadcasts it to the node's neighbors. It returns the message ID if successful.
func broadcastData(c echo.Context) error {
	var request Request
	if err := c.Bind(&request); err != nil {
		log.Info(err.Error())
		return c.JSON(http.StatusBadRequest, Response{Error: err.Error()})
	}

	//TODO: to check max payload size allowed, if exceeding return an error

	marshalUtil := marshalutil.New(request.Data)
	parsedPayload, err := payload.Parse(marshalUtil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{Error: "Not a valid Payload"})
	}
	tx := messagelayer.MessageFactory.IssuePayload(parsedPayload)

	return c.JSON(http.StatusOK, Response{Id: tx.Id().String()})
}

type Response struct {
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

type Request struct {
	Data []byte `json:"data"`
}
