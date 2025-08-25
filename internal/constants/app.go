package constants

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Header string

const (
	ContentType     Header = "Content-Type"
	ContentEncoding Header = "Content-Encoding"
)

type ContentTypes string

const (
	Json ContentTypes = "application/json"
)

type Compression string

const (
	Gzip Compression = "gzip"
)

type Roles string

const (
	RoleSessionUser   Roles = "session_user"   // With session
	RoleAnonymousUser Roles = "anonymous_user" // without session
	RoleSystemAdmin   Roles = "system_admin"
	Roleadmin         Roles = "admin"
	Roleuser          Roles = "user"
)

type Status string

const (
	StatusSuccess         Status = "success"
	StatusFail            Status = "fail"
	StatusPending         Status = "pending"
	StatusActive          Status = "active"   // The entity is currently enabled, usable, or functioning as expected
	StatusInactive        Status = "inactive" // // Temporarily disabled or not currently in use, but not due to any problem or penalty
	StatusTemporaryLocked Status = "temporary_locked"
	StatusSuspended       Status = "suspended" // Temporarily restricted due to some violation, risk, or admin action, but possibly reversible.
	StatusSoftDeleted     Status = "soft_deleted"
	StatusBlocked         Status = "blocked" // Permanently or strongly restricted — typically by system or admin.
)

type Operations string

const (
	Create Operations = "create"
	Update Operations = "update"
)

const (
	Authorization string = "Authorization"
	AuthType      string = "Bearer "
)
