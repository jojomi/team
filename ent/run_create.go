// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/jojomi/team/ent/run"
)

// RunCreate is the builder for creating a Run entity.
type RunCreate struct {
	config
	mutation *RunMutation
	hooks    []Hook
}

// SetJob sets the "job" field.
func (rc *RunCreate) SetJob(s string) *RunCreate {
	rc.mutation.SetJob(s)
	return rc
}

// SetStart sets the "start" field.
func (rc *RunCreate) SetStart(t time.Time) *RunCreate {
	rc.mutation.SetStart(t)
	return rc
}

// SetNillableStart sets the "start" field if the given value is not nil.
func (rc *RunCreate) SetNillableStart(t *time.Time) *RunCreate {
	if t != nil {
		rc.SetStart(*t)
	}
	return rc
}

// SetEnd sets the "end" field.
func (rc *RunCreate) SetEnd(t time.Time) *RunCreate {
	rc.mutation.SetEnd(t)
	return rc
}

// SetNillableEnd sets the "end" field if the given value is not nil.
func (rc *RunCreate) SetNillableEnd(t *time.Time) *RunCreate {
	if t != nil {
		rc.SetEnd(*t)
	}
	return rc
}

// SetStatus sets the "status" field.
func (rc *RunCreate) SetStatus(r run.Status) *RunCreate {
	rc.mutation.SetStatus(r)
	return rc
}

// SetLog sets the "log" field.
func (rc *RunCreate) SetLog(s string) *RunCreate {
	rc.mutation.SetLog(s)
	return rc
}

// SetNillableLog sets the "log" field if the given value is not nil.
func (rc *RunCreate) SetNillableLog(s *string) *RunCreate {
	if s != nil {
		rc.SetLog(*s)
	}
	return rc
}

// SetID sets the "id" field.
func (rc *RunCreate) SetID(u uuid.UUID) *RunCreate {
	rc.mutation.SetID(u)
	return rc
}

// Mutation returns the RunMutation object of the builder.
func (rc *RunCreate) Mutation() *RunMutation {
	return rc.mutation
}

// Save creates the Run in the database.
func (rc *RunCreate) Save(ctx context.Context) (*Run, error) {
	var (
		err  error
		node *Run
	)
	rc.defaults()
	if len(rc.hooks) == 0 {
		if err = rc.check(); err != nil {
			return nil, err
		}
		node, err = rc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*RunMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = rc.check(); err != nil {
				return nil, err
			}
			rc.mutation = mutation
			if node, err = rc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(rc.hooks) - 1; i >= 0; i-- {
			if rc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = rc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, rc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (rc *RunCreate) SaveX(ctx context.Context) *Run {
	v, err := rc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rc *RunCreate) Exec(ctx context.Context) error {
	_, err := rc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rc *RunCreate) ExecX(ctx context.Context) {
	if err := rc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (rc *RunCreate) defaults() {
	if _, ok := rc.mutation.Start(); !ok {
		v := run.DefaultStart()
		rc.mutation.SetStart(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (rc *RunCreate) check() error {
	if _, ok := rc.mutation.Job(); !ok {
		return &ValidationError{Name: "job", err: errors.New(`ent: missing required field "Run.job"`)}
	}
	if _, ok := rc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "Run.status"`)}
	}
	if v, ok := rc.mutation.Status(); ok {
		if err := run.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "Run.status": %w`, err)}
		}
	}
	return nil
}

func (rc *RunCreate) sqlSave(ctx context.Context) (*Run, error) {
	_node, _spec := rc.createSpec()
	if err := sqlgraph.CreateNode(ctx, rc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	return _node, nil
}

func (rc *RunCreate) createSpec() (*Run, *sqlgraph.CreateSpec) {
	var (
		_node = &Run{config: rc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: run.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: run.FieldID,
			},
		}
	)
	if id, ok := rc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := rc.mutation.Job(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: run.FieldJob,
		})
		_node.Job = value
	}
	if value, ok := rc.mutation.Start(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: run.FieldStart,
		})
		_node.Start = value
	}
	if value, ok := rc.mutation.End(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: run.FieldEnd,
		})
		_node.End = value
	}
	if value, ok := rc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: run.FieldStatus,
		})
		_node.Status = value
	}
	if value, ok := rc.mutation.Log(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: run.FieldLog,
		})
		_node.Log = value
	}
	return _node, _spec
}

// RunCreateBulk is the builder for creating many Run entities in bulk.
type RunCreateBulk struct {
	config
	builders []*RunCreate
}

// Save creates the Run entities in the database.
func (rcb *RunCreateBulk) Save(ctx context.Context) ([]*Run, error) {
	specs := make([]*sqlgraph.CreateSpec, len(rcb.builders))
	nodes := make([]*Run, len(rcb.builders))
	mutators := make([]Mutator, len(rcb.builders))
	for i := range rcb.builders {
		func(i int, root context.Context) {
			builder := rcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*RunMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, rcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, rcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, rcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (rcb *RunCreateBulk) SaveX(ctx context.Context) []*Run {
	v, err := rcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (rcb *RunCreateBulk) Exec(ctx context.Context) error {
	_, err := rcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (rcb *RunCreateBulk) ExecX(ctx context.Context) {
	if err := rcb.Exec(ctx); err != nil {
		panic(err)
	}
}