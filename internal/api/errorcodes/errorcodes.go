package errorcodes

const (
	// CodeOK Status code for successful request
	CodeOK string = "0"
	// CodeMaintenance Maintenanace mode or unknown fault
	CodeMaintenance string = "1000"
	// CodeCache Failed to connect to the cache (Redis)
	CodeCache string = "1200"
	// CodeCDN Failed to communicate with CDN service
	CodeCDN string = "1201"
	// CodeIdentity Failed to communicate with identity service
	CodeIdentity string = "1202"
	// CodeDBDiscovery Failed to communicate with DB discovery service
	CodeDBDiscovery string = "1203"
	// CodeDatabase Failed to connect to the database
	CodeDatabase string = "1003"
	// CodeUnknownRequest Unknown request or 'request' input missing
	CodeUnknownRequest = "1005"
	// CodeRequiredParameterMissing Required parameter is missing
	CodeRequiredParameterMissing = 1010
	// CodeUnauthenticated Status code when authentication fails
	CodeUnauthenticated string = "1051"
)
