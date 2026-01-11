package config

// Root config
type Config struct {
	Version  string             `yaml:"version"`
	Project  Project            `yaml:"project"`
	Database Database           `yaml:"database"`
	Domains  map[string]Domain  `yaml:"domains"`
	Policies []Policy           `yaml:"policies,omitempty"`
	Docs     *Documentation     `yaml:"documentation,omitempty"`
}

// Project metadata
type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Environment string `yaml:"environment"`
}

// Database connection info
type Database struct {
	Engine      string      `yaml:"engine"`
	Host        string      `yaml:"host"`
	Port        int         `yaml:"port"`
	Name        string      `yaml:"name"`
	SSL         bool        `yaml:"ssl"`
	Credentials Credentials `yaml:"credentials"`
}

type Credentials struct {
	User        string `yaml:"user"`
	PasswordEnv string `yaml:"password_env"`
}

// Domain definition
type Domain struct {
	Description string            `yaml:"description,omitempty"`
	Schema      string            `yaml:"schema"`
	Owner       string            `yaml:"owner"`
	Database    *DomainDatabase   `yaml:"database,omitempty"`
	Owns        *OwnedResources   `yaml:"owns,omitempty"`
	Reads       map[string]Read   `yaml:"reads,omitempty"`
	Access      *AccessControl   `yaml:"access,omitempty"`
	Sync        *SyncConfig      `yaml:"sync,omitempty"`
}

// Optional domain-specific database
type DomainDatabase struct {
	Name string `yaml:"name"`
}

// Owned tables
type OwnedResources struct {
	Tables map[string]Table `yaml:"tables"`
}

// Table definition
type Table struct {
	Description string             `yaml:"description,omitempty"`
	Columns     map[string]Column  `yaml:"columns"`
	Indexes     []Index            `yaml:"indexes,omitempty"`
}

// Column definition
type Column struct {
	Type       string `yaml:"type"`
	Nullable   bool   `yaml:"nullable,omitempty"`
	PrimaryKey bool   `yaml:"primary_key,omitempty"`
	Unique     bool   `yaml:"unique,omitempty"`
}

// Index definition
type Index struct {
	Columns []string `yaml:"columns"`
	Unique  bool     `yaml:"unique,omitempty"`
}

// Read model / view
type Read struct {
	From        string       `yaml:"from"`
	Columns     []string     `yaml:"columns"`
	Access      AccessMode  `yaml:"access,omitempty"`
	Consistency Consistency `yaml:"consistency,omitempty"`
	Materialized bool       `yaml:"materialized,omitempty"`
}

type AccessMode struct {
	Mode string `yaml:"mode"`
}

type Consistency struct {
	Type string `yaml:"type"`
}

// Explicit access override
type AccessControl struct {
	Users []UserAccess `yaml:"users"`
}

type UserAccess struct {
	Name       string             `yaml:"name"`
	Privileges []ResourcePrivilege `yaml:"privileges"`
}

type ResourcePrivilege struct {
	Resource string   `yaml:"resource"`
	Actions  []string `yaml:"actions"`
}

// Sync / CDC (future)
type SyncConfig struct {
	Enabled bool       `yaml:"enabled"`
	Source  SyncSource `yaml:"source"`
	Mode    string     `yaml:"mode"`
	Target  SyncTarget `yaml:"target"`
}

type SyncSource struct {
	Domain string `yaml:"domain"`
	Table  string `yaml:"table"`
}

type SyncTarget struct {
	Database string `yaml:"database"`
	Table    string `yaml:"table"`
}

// Policies (future)
type Policy struct {
	Name      string `yaml:"name"`
	Domain    string `yaml:"domain"`
	Type      string `yaml:"type"`
	Condition string `yaml:"condition"`
}

// Documentation output
type Documentation struct {
	Generate bool   `yaml:"generate"`
	Format   string `yaml:"format"`
	Output   string `yaml:"output"`
}
