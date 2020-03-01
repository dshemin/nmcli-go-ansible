package nmcli

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
)

// Connection represent NetworkManager connection
type Connection struct {
	ConName string `json:"con_name"` // connection.id
	Ifname  string `json:"ifname"`   // connection.interface-name
	Type    string `json:"type"`     // connection.type
	Master  string `json:"master"`   // connection.master
}

const (
	fieldsForCheck = "connection.id,connection.type,connection.interface-name,connection.master"
	numberOfFields = 4
)

// ErrNotFound error which indicates that requested connection is not exists
var ErrNotFound = errors.New("connection not exists")

// FindConnectionByName find NM connection by name
// Return ErrNotFound error if connection with specified name not exists
func FindConnectionByName(ctx context.Context, name string) (Connection, error) {
	c := Connection{}

	out, err := nmcli(
		ctx,
		"-terse",
		"-mode",
		"tabular",
		"--fields",
		fieldsForCheck,
		"connection",
		"show",
		name,
	)

	if err != nil {
		exitErr := &exec.ExitError{}
		if errors.As(err, &exitErr) {
			// Very fragile check ...
			if exitErr.ProcessState.ExitCode() == 10 {
				return c, ErrNotFound
			}

			return c, err
		}

		return c, err
	}

	raw := make([]string, 0, numberOfFields)
	s := bufio.NewScanner(bytes.NewReader(out))
	for s.Scan() {
		raw = append(raw, s.Text())
	}

	c.ConName = raw[0]
	c.Type = raw[1]
	c.Ifname = raw[2]
	c.Master = raw[3]

	return c, nil
}

// EqualTo checks that current connection are equal to specified
func (c Connection) EqualTo(v Connection) bool {
	return c.ConName == v.ConName &&
		c.Type == v.Type &&
		c.Ifname == v.Ifname &&
		c.Master == v.Master
}

// Create creates this connection
func (c Connection) Create(ctx context.Context) error {
	args := []string{
		"connection",
		"add",
	}

	if c.ConName != "" {
		args = append(args, "con-name")
		args = append(args, c.ConName)
	}

	if c.Ifname != "" {
		args = append(args, "ifname")
		args = append(args, c.Ifname)
	}

	if c.Type != "" {
		args = append(args, "type")
		args = append(args, c.Type)
	}

	if c.Master != "" {
		args = append(args, "master")
		args = append(args, c.Master)
	}

	_, err := nmcli(ctx, args...)
	if err != nil {
		return fmt.Errorf("cannot create connection %#v: %w", c, err)
	}
	return nil
}

// Remove removes this connection
func (c Connection) Remove(ctx context.Context) error {
	_, err := nmcli(ctx, "connection", "delete", c.ConName)
	if err != nil {
		return fmt.Errorf("cannot delete connection %q: %w", c.ConName, err)
	}
	return nil
}
