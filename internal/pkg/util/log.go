package util

import (
	"encoding/json"
	"fmt"

	"github.com/pterm/pterm"
)

var Logger = pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

const MessageNotIncoming = "Bukan Incoming"

func HandleAppError(err error, function string, step string, fatal bool) error {
	if err != nil {
		if fatal {
			Logger.Fatal(fmt.Sprintf("Fatal error in function: %v", function), Logger.Args("Step", step, "Details", err))
			return err
		} else {
			Logger.Error(fmt.Sprintf("Error in function: %v", function), Logger.Args("Step", step, "Details", err))
		}
	}
	return nil
}

func LoggerChannel(data interface{}, name, tenantId string) {
	nameLog := tenantId + " - " + name
	logEnv := GodotEnv("LOG")
	tenantLog := GodotEnv("TENANT_LOG")

	if logEnv == "1" {
		dataJSON, err := json.Marshal(data)
		if err != nil {
			HandleAppError(err, "util", "LoggerChannel", false)
			return
		}

		if tenantLog == "ALL" {
			pterm.Info.Printfln("%s: %s", nameLog, dataJSON)
		} else if tenantLog == tenantId {
			pterm.Info.Printfln("%s: %s", nameLog, dataJSON)
		}
	}
}
