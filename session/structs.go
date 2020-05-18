package session

type FlagArguments struct {
	TimeOut     int
	TargetURL   string
	HardLimit   int
	DomainScope string
	NWorkers    int
	FileName    string
	StartZMQ    bool
}
