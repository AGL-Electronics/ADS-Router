package router

// Header represents the ADS header structure.
type Header struct {
	Length     uint32 // Length of the header and data
	InvokeID   uint32 // Invoke ID for the request
	Flags      uint32 // Flags for the request
	Reserved   uint32 // Reserved for future use
	SourceNet  uint32 // Source network number
	SourcePort uint32 // Source port number
	TargetNet  uint32 // Target network number
	TargetPort uint32 // Target port number

	// Data is the payload of the message, which can be of variable length.
	// 	Data []byte // This is not included in the struct, as it is variable length
}

func Rewriteheaderincoming(*Header) {
	// Rewrite the header fields as needed
	Header.Length = 0     // Set to 0 for now, will be set later
	Header.InvokeID = 0   // Set to 0 for now, will be set later
	Header.Flags = 0      // Set to 0 for now, will be set later
	Header.Reserved = 0   // Set to 0 for now, will be set later
	Header.SourceNet = 0  // Set to 0 for now, will be set later
	Header.SourcePort = 0 // Set to 0 for now, will be set later
	Header.TargetNet = 0  // Set to 0 for now, will be set later
	Header.TargetPort = 0 // Set to 0 for now, will be set later
}

func Rewriteheaderoutgoing(*Header) {
	// Rewrite the header fields as needed
	Header.Length = 0     // Set to 0 for now, will be set later
	Header.InvokeID = 0   // Set to 0 for now, will be set later
	Header.Flags = 0      // Set to 0 for now, will be set later
	Header.Reserved = 0   // Set to 0 for now, will be set later
	Header.SourceNet = 0  // Set to 0 for now, will be set later
	Header.SourcePort = 0 // Set to 0 for now, will be set later
	Header.TargetNet = 0  // Set to 0 for now, will be set later
	Header.TargetPort = 0 // Set to 0 for now, will be set later
}
