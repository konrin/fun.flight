package resiver

const (
	ResiverStateStopped = iota
	ResiverStateRunning
	ResiverStateConnecting

	ResiverTypeRaw   = "raw"
	ResiverTypeBeast = "beast"
)

type Resiver struct {
	state int
	cfg   ResiverConfig
}

func NewResiver(cfg ResiverConfig) *Resiver {
	return &Resiver{
		cfg:   cfg,
		state: ResiverStateStopped,
	}
}

func (r *Resiver) Name() string {
	return r.cfg.Name
}

func (r *Resiver) Type() string {
	return r.cfg.Type
}

func (r *Resiver) State() int {
	return r.state
}
