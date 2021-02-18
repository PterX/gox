package transaction

type Propagation int

const (
	PROPAGATION_REQUIRED Propagation = iota
	PROPAGATION_SUPPORTS
	PROPAGATION_MANDATORY
	PROPAGATION_REQUIRES_NEW
	PROPAGATION_NOT_SUPPORTED
	PROPAGATION_NEVER
	PROPAGATION_NESTED
)

const TIMEOUT_DEFAULT = -1

const COMMENT_NAME = "@Transactional"

const GO_FILE_SUFIX = ".go"
