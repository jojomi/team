// Code generated by ent, DO NOT EDIT.

package run

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/jojomi/team/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Run {
	return predicate.Run(sql.FieldLTE(FieldID, id))
}

// Job applies equality check predicate on the "job" field. It's identical to JobEQ.
func Job(v string) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldJob, v))
}

// Start applies equality check predicate on the "start" field. It's identical to StartEQ.
func Start(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldStart, v))
}

// End applies equality check predicate on the "end" field. It's identical to EndEQ.
func End(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldEnd, v))
}

// Log applies equality check predicate on the "log" field. It's identical to LogEQ.
func Log(v string) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldLog, v))
}

// JobEQ applies the EQ predicate on the "job" field.
func JobEQ(v string) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldJob, v))
}

// JobNEQ applies the NEQ predicate on the "job" field.
func JobNEQ(v string) predicate.Run {
	return predicate.Run(sql.FieldNEQ(FieldJob, v))
}

// JobIn applies the In predicate on the "job" field.
func JobIn(vs ...string) predicate.Run {
	return predicate.Run(sql.FieldIn(FieldJob, vs...))
}

// JobNotIn applies the NotIn predicate on the "job" field.
func JobNotIn(vs ...string) predicate.Run {
	return predicate.Run(sql.FieldNotIn(FieldJob, vs...))
}

// JobGT applies the GT predicate on the "job" field.
func JobGT(v string) predicate.Run {
	return predicate.Run(sql.FieldGT(FieldJob, v))
}

// JobGTE applies the GTE predicate on the "job" field.
func JobGTE(v string) predicate.Run {
	return predicate.Run(sql.FieldGTE(FieldJob, v))
}

// JobLT applies the LT predicate on the "job" field.
func JobLT(v string) predicate.Run {
	return predicate.Run(sql.FieldLT(FieldJob, v))
}

// JobLTE applies the LTE predicate on the "job" field.
func JobLTE(v string) predicate.Run {
	return predicate.Run(sql.FieldLTE(FieldJob, v))
}

// JobContains applies the Contains predicate on the "job" field.
func JobContains(v string) predicate.Run {
	return predicate.Run(sql.FieldContains(FieldJob, v))
}

// JobHasPrefix applies the HasPrefix predicate on the "job" field.
func JobHasPrefix(v string) predicate.Run {
	return predicate.Run(sql.FieldHasPrefix(FieldJob, v))
}

// JobHasSuffix applies the HasSuffix predicate on the "job" field.
func JobHasSuffix(v string) predicate.Run {
	return predicate.Run(sql.FieldHasSuffix(FieldJob, v))
}

// JobEqualFold applies the EqualFold predicate on the "job" field.
func JobEqualFold(v string) predicate.Run {
	return predicate.Run(sql.FieldEqualFold(FieldJob, v))
}

// JobContainsFold applies the ContainsFold predicate on the "job" field.
func JobContainsFold(v string) predicate.Run {
	return predicate.Run(sql.FieldContainsFold(FieldJob, v))
}

// StartEQ applies the EQ predicate on the "start" field.
func StartEQ(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldStart, v))
}

// StartNEQ applies the NEQ predicate on the "start" field.
func StartNEQ(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldNEQ(FieldStart, v))
}

// StartIn applies the In predicate on the "start" field.
func StartIn(vs ...time.Time) predicate.Run {
	return predicate.Run(sql.FieldIn(FieldStart, vs...))
}

// StartNotIn applies the NotIn predicate on the "start" field.
func StartNotIn(vs ...time.Time) predicate.Run {
	return predicate.Run(sql.FieldNotIn(FieldStart, vs...))
}

// StartGT applies the GT predicate on the "start" field.
func StartGT(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldGT(FieldStart, v))
}

// StartGTE applies the GTE predicate on the "start" field.
func StartGTE(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldGTE(FieldStart, v))
}

// StartLT applies the LT predicate on the "start" field.
func StartLT(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldLT(FieldStart, v))
}

// StartLTE applies the LTE predicate on the "start" field.
func StartLTE(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldLTE(FieldStart, v))
}

// StartIsNil applies the IsNil predicate on the "start" field.
func StartIsNil() predicate.Run {
	return predicate.Run(sql.FieldIsNull(FieldStart))
}

// StartNotNil applies the NotNil predicate on the "start" field.
func StartNotNil() predicate.Run {
	return predicate.Run(sql.FieldNotNull(FieldStart))
}

// EndEQ applies the EQ predicate on the "end" field.
func EndEQ(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldEnd, v))
}

// EndNEQ applies the NEQ predicate on the "end" field.
func EndNEQ(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldNEQ(FieldEnd, v))
}

// EndIn applies the In predicate on the "end" field.
func EndIn(vs ...time.Time) predicate.Run {
	return predicate.Run(sql.FieldIn(FieldEnd, vs...))
}

// EndNotIn applies the NotIn predicate on the "end" field.
func EndNotIn(vs ...time.Time) predicate.Run {
	return predicate.Run(sql.FieldNotIn(FieldEnd, vs...))
}

// EndGT applies the GT predicate on the "end" field.
func EndGT(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldGT(FieldEnd, v))
}

// EndGTE applies the GTE predicate on the "end" field.
func EndGTE(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldGTE(FieldEnd, v))
}

// EndLT applies the LT predicate on the "end" field.
func EndLT(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldLT(FieldEnd, v))
}

// EndLTE applies the LTE predicate on the "end" field.
func EndLTE(v time.Time) predicate.Run {
	return predicate.Run(sql.FieldLTE(FieldEnd, v))
}

// EndIsNil applies the IsNil predicate on the "end" field.
func EndIsNil() predicate.Run {
	return predicate.Run(sql.FieldIsNull(FieldEnd))
}

// EndNotNil applies the NotNil predicate on the "end" field.
func EndNotNil() predicate.Run {
	return predicate.Run(sql.FieldNotNull(FieldEnd))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v Status) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v Status) predicate.Run {
	return predicate.Run(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...Status) predicate.Run {
	return predicate.Run(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...Status) predicate.Run {
	return predicate.Run(sql.FieldNotIn(FieldStatus, vs...))
}

// LogEQ applies the EQ predicate on the "log" field.
func LogEQ(v string) predicate.Run {
	return predicate.Run(sql.FieldEQ(FieldLog, v))
}

// LogNEQ applies the NEQ predicate on the "log" field.
func LogNEQ(v string) predicate.Run {
	return predicate.Run(sql.FieldNEQ(FieldLog, v))
}

// LogIn applies the In predicate on the "log" field.
func LogIn(vs ...string) predicate.Run {
	return predicate.Run(sql.FieldIn(FieldLog, vs...))
}

// LogNotIn applies the NotIn predicate on the "log" field.
func LogNotIn(vs ...string) predicate.Run {
	return predicate.Run(sql.FieldNotIn(FieldLog, vs...))
}

// LogGT applies the GT predicate on the "log" field.
func LogGT(v string) predicate.Run {
	return predicate.Run(sql.FieldGT(FieldLog, v))
}

// LogGTE applies the GTE predicate on the "log" field.
func LogGTE(v string) predicate.Run {
	return predicate.Run(sql.FieldGTE(FieldLog, v))
}

// LogLT applies the LT predicate on the "log" field.
func LogLT(v string) predicate.Run {
	return predicate.Run(sql.FieldLT(FieldLog, v))
}

// LogLTE applies the LTE predicate on the "log" field.
func LogLTE(v string) predicate.Run {
	return predicate.Run(sql.FieldLTE(FieldLog, v))
}

// LogContains applies the Contains predicate on the "log" field.
func LogContains(v string) predicate.Run {
	return predicate.Run(sql.FieldContains(FieldLog, v))
}

// LogHasPrefix applies the HasPrefix predicate on the "log" field.
func LogHasPrefix(v string) predicate.Run {
	return predicate.Run(sql.FieldHasPrefix(FieldLog, v))
}

// LogHasSuffix applies the HasSuffix predicate on the "log" field.
func LogHasSuffix(v string) predicate.Run {
	return predicate.Run(sql.FieldHasSuffix(FieldLog, v))
}

// LogIsNil applies the IsNil predicate on the "log" field.
func LogIsNil() predicate.Run {
	return predicate.Run(sql.FieldIsNull(FieldLog))
}

// LogNotNil applies the NotNil predicate on the "log" field.
func LogNotNil() predicate.Run {
	return predicate.Run(sql.FieldNotNull(FieldLog))
}

// LogEqualFold applies the EqualFold predicate on the "log" field.
func LogEqualFold(v string) predicate.Run {
	return predicate.Run(sql.FieldEqualFold(FieldLog, v))
}

// LogContainsFold applies the ContainsFold predicate on the "log" field.
func LogContainsFold(v string) predicate.Run {
	return predicate.Run(sql.FieldContainsFold(FieldLog, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Run) predicate.Run {
	return predicate.Run(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Run) predicate.Run {
	return predicate.Run(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Run) predicate.Run {
	return predicate.Run(sql.NotPredicates(p))
}
