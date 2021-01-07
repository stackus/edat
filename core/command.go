package core

import (
	"fmt"
	"reflect"
)

// Command interface
type Command interface {
	CommandName() string
}

// SerializeCommand serializes commands with a registered marshaller
func SerializeCommand(v Command) ([]byte, error) {
	return marshal(v.CommandName(), v)
}

// DeserializeCommand deserializes the command data using a registered marshaller returning a *Command
func DeserializeCommand(commandName string, data []byte) (Command, error) {
	cmd, err := unmarshal(commandName, data)
	if err != nil {
		return nil, err
	}

	if cmd != nil {
		if _, ok := cmd.(Command); !ok {
			return nil, fmt.Errorf("`%s` was registered but not registered as a command", commandName)
		}
	}

	return cmd.(Command), nil
}

// RegisterCommands registers one or more commands with a registered marshaller
//
// Register commands using any form desired "&MyCommand{}", "MyCommand{}", "(*MyCommand)(nil)"
//
// Commands must be registered after first registering a marshaller you wish to use
func RegisterCommands(commands ...Command) {
	for _, command := range commands {
		if v := reflect.ValueOf(command); v.Kind() == reflect.Ptr && v.Pointer() == 0 {
			commandName := reflect.Zero(reflect.TypeOf(command).Elem()).Interface().(Command).CommandName()
			registerType(commandName, command)
		} else {
			registerType(command.CommandName(), command)
		}
	}
}
