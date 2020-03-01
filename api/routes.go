package api

import "net/http"

// Route holds the data which will be used while setting up router (API).
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routs implements array of Route.
type Routs []Route

var (
	// Routes holds all the endpoints that neuron exposes.
	Routes = Routs{

		Route{"Identity", "GET", "/", identity},
		Route{"GetIdentity", "GET", "/getidentity", getIdentity},
		Route{"CreateIdentity", "CREATE", "/createidentity", createIdentity},
		Route{"DeleteIdentity", "DELETE", "/deleteidentity", deleteIdentity},
		Route{"UpdateIdentity", "UPDATE", "/updateidentity", updateIdentity},
	}
)
