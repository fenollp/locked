package locked

// T is the root of a Lockfile
type T struct {
	At    string  `hcl:"at,optional"`
	Track []Track `hcl:"track,block"`
}

// Track defines a rewriting rule
type Track struct {
	Track    string     `hcl:"track,label"`
	Use      string     `hcl:"use"`
	Tracking []Tracking `hcl:"tracking,block"`
}

// Tracking represents the effect of a rule
type Tracking struct {
	Tracked string `hcl:"tracked,label"`
	At      string `hcl:"at"`
	Gives   string `hcl:"gives"`
}
