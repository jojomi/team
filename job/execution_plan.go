package job

//go:generate go-enum --file "${GOFILE}"

// ENUM(skip, logDone, execute)
type ExecutionPlan uint8
