package configs

type PortCheck struct {
	port     int
	required bool
	label    string
}

var PlcFingerprint = []PortCheck{
	//these all follow the:  portnr, required? y/n, discription  template, add more here as needed

	{48898, true, "AMS router"},
	{8016, false, "TwinCAT TCP"},
	{443, true, "Web UI (HTTPS)"},
	{8015, false, "System service"},
}

var Subnet = "192.168.88." //leave last space of ip blank for entire subnet
