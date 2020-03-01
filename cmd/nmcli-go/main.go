package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dshemin/nmcli-go-ansible/internal/ansible"
	"github.com/dshemin/nmcli-go-ansible/internal/nmcli"
)

// connectionState represent available connection states
type connectionState string

const (
	absent  connectionState = "absent"
	present connectionState = "present"
)

type arguments struct {
	nmcli.Connection
	State connectionState `json:"state"`
}

func main() {
	expected, err := parseArguments(os.Args)
	if err != nil {
		returnResponse(ansible.FailResponsef("Cannot parse argument %#v: %s", os.Args, err))
	}
	returnResponse(bringToTheDesiredState(expected))
}

func parseArguments(args []string) (arguments, error) {
	expected := arguments{}

	if len(args) != 2 {
		return expected, errors.New("no argument file provided")
	}

	argsFile := args[1]

	f, err := os.Open(argsFile)
	if err != nil {
		return expected, fmt.Errorf("could not read configuration file %q: %s", argsFile, err)
	}

	err = json.NewDecoder(f).Decode(&expected)
	if err != nil {
		return expected, fmt.Errorf("configuration file %q not valid JSON: %s", argsFile, err)
	}

	if expected.State != absent && expected.State != present {
		return expected, fmt.Errorf("unknown state %q", expected.State)
	}

	return expected, nil
}

const execTimeout = 10 * time.Second

func bringToTheDesiredState(expected arguments) ansible.Response {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	actual, err := nmcli.FindConnectionByName(ctx, expected.ConName)
	if err != nil {
		if err != nmcli.ErrNotFound {
			return ansible.FailResponsef("Cannot find connection %q: %s", expected.ConName, err)
		}

		// Expected connection not exists.

		if expected.State == absent {
			// Connection not exists, but we expect this.
			return ansible.Response{}
		}

		// Connection not exists but should.
		// Create new one.
		err = expected.Create(ctx)
		if err != nil {
			return ansible.ErrorResponse(err)
		}
		return ansible.Response{Changed: true}
	}

	// Expected connection is exists. We get next possible cases:
	// * Connection should be removed. In that case we simple remove it and all.
	// * Connection should exist. In that case we should check that actual connection
	//   equals to expected and update it configuration if not.

	if expected.State == absent {
		err = actual.Remove(ctx)
		if err != nil {
			return ansible.FailResponsef("Cannot remove connection %#v: %s", actual, err)
		}

		return ansible.Response{Changed: true}
	}

	// Connection should exist.
	if expected.EqualTo(actual) {
		return ansible.Response{}
	}

	err = actual.Remove(ctx)
	if err != nil {
		return ansible.FailResponsef("Cannot remove connection %#v: %s", actual, err)
	}

	err = expected.Create(ctx)
	if err != nil {
		return ansible.FailResponsef("Cannot create connection %#v: %s", expected, err)
	}
	return ansible.Response{Changed: true}
}

//revive:disable:deep-exit It's safe to use `os.Exit` here
func returnResponse(r ansible.Response) {
	var response []byte
	var err error

	response, err = json.Marshal(r)
	if err != nil {
		response = []byte(`{"changed":false,"msg":"Invalid response object","failed":true}`)
	}
	fmt.Println(string(response))
	if r.Failed {
		os.Exit(1)
	}

	os.Exit(0)
}

//revive:enable:deep-exit
