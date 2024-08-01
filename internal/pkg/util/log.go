package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

var Logger = pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)

const MessageNotIncoming = "Bukan Incoming"

func HandleAppError(err error, function string, step string, fatal bool) {
	if err != nil {
		if fatal {
			Logger.Fatal(fmt.Sprintf("Fatal error in function: %v", function), Logger.Args("Step", step, "Details", err))
			os.Exit(1)
		} else {
			Logger.Error(fmt.Sprintf("Error in function: %v", function), Logger.Args("Step", step, "Details", err))
		}
	}
}

func LoggerChannel(data interface{}, name, tenantID string) {
	nameLog := tenantID + " - " + name
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
		} else if tenantLog == tenantID {
			pterm.Info.Printfln("%s: %s", nameLog, dataJSON)
		}
	}
}
