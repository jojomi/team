// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/jojomi/team/ent/run"
	"github.com/jojomi/team/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	runFields := schema.Run{}.Fields()
	_ = runFields
	// runDescStart is the schema descriptor for start field.
	runDescStart := runFields[2].Descriptor()
	// run.DefaultStart holds the default value on creation for the start field.
	run.DefaultStart = runDescStart.Default.(func() time.Time)
}
