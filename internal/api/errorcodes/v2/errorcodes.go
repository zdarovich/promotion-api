package errorcodes

const (
	// CodeOK Status code for successful request
	CodeOK = 0
	// CodeMaintenance Maintenanace mode or unknown fault
	CodeMaintenance = 2000
	// CodeCache Failed to connect to the cache (Redis)
	CodeCache = 2001
	// CodeCDN Failed to communicate with CDN service
	CodeCDN = 2002
	// CodeIdentity Failed to communicate with identity service
	CodeIdentity = 2003
	// CodeDBDiscovery Failed to communicate with DB discovery service
	CodeDBDiscovery = 2004
	// CodeDatabase Failed to connect to the database
	CodeDatabase = 2005
	// CodeDatabaseQuery Failed to parse sql
	CodeDatabaseQuery = 2006
	// CodeUnknownRequest Unknown request or 'request' input missing
	CodeUnknownRequest = 2007
	// CodeRequiredParameterMissing Required parameter is missing
	CodeRequiredParameterMissing = 2008
	// CodeUnauthenticated Status code when authentication fails
	CodeUnauthenticated = 2009
	// CodeNoViewRights Status when no viewing rights
	CodeNoViewRights = 2010
	// CodeDebugModeDisabled Status when debug is used but it has been disabled
	CodeDebugModeDisabled = 2011
)

// GetDescriptions returns error code descriptions
func GetDescriptions() map[int]string {

	return map[int]string{
		CodeOK:                       "No errors",
		CodeMaintenance:              "Application in maintenance mode",
		CodeCache:                    "Application cache failure",
		CodeCDN:                      "CDN service communication error",
		CodeIdentity:                 "Identity service communication error",
		CodeDBDiscovery:              "DB discovery service communication error",
		CodeDatabase:                 "DB communication error",
		CodeDatabaseQuery:            "DB data reading error",
		CodeUnknownRequest:           "Unknown request",
		CodeRequiredParameterMissing: "Required parameter missing",
		CodeUnauthenticated:          "Unable to authenticate the request",
		CodeNoViewRights:             "User has no access to the request",
	}
}

// GetDescription returns a description to the requested code
func GetDescription(code int) string {

	var descriptions = GetDescriptions()
	if description, ok := descriptions[code]; ok {
		return description
	}
	return ""
}
