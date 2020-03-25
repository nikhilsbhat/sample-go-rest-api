package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// FetchIdentity helps in fetching the identities of the person
type FetchIdentity struct {
	Identity Identity
}

// Identity holds the identity of a person
type Identity struct {
	FirstName  string
	SecondName string
	Dob        string
}

// Error holds the error message for neuron UI
type Error struct {
	Error string
}

var (
	identities = make(map[int]Identity)
)

// Initialises basic identities during startup of the app.
func init() {
	conf, err := getConfig()
	if err != nil {
		fmt.Fprintf(os.Stdout, err.(error).Error())
		os.Exit(1)
	}
	if len(conf.DefaultIdentites) != 0 {
		for i, idnt := range conf.DefaultIdentites {
			identities[i] = idnt
		}
	} else {
		identities = map[int]Identity{
			1: Identity{FirstName: "John", SecondName: "Doe", Dob: "11/11/2000"},
			2: Identity{FirstName: "Bob", SecondName: "Builder", Dob: "10/10/2000"},
			3: Identity{FirstName: "Jag", SecondName: "Dragger", Dob: "9/09/2000"},
		}
	}
}

func getIdentity(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.MarshalIndent(Error{"Not received input in a valid format"}, "", "  ")
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var idn FetchIdentity

		err = json.Unmarshal(body, &idn)
		if err != nil {
			jsonval, _ := json.MarshalIndent(Error{"Unable to unmarshal the entered input. Provide input in valid format"}, "", "  ")
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			getIdentityResponse, err := idn.getIdentity()
			if err != nil {
				fmt.Fprintf(rw, "%v\n", err)
			} else {
				jsonval, _ := json.MarshalIndent(getIdentityResponse, "", " ")
				fmt.Fprintf(rw, "%v\n", "Retrieved identity can be found below")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func createIdentity(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.MarshalIndent(Error{"Not received input in a valid format"}, "", "  ")
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var idn FetchIdentity

		err = json.Unmarshal(body, &idn)
		if err != nil {
			jsonval, _ := json.MarshalIndent(Error{"Unable to unmarshal the entered input. Provide input in valid format"}, "", "  ")
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			createIdentityResponse, err := idn.createIdentity()
			if err != nil {
				fmt.Fprintf(rw, "%v\n", err)
			} else {
				jsonval, _ := json.MarshalIndent(createIdentityResponse, "", " ")
				fmt.Fprintf(rw, "%v\n", "Identity was created successfully, find the created identity below")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func updateIdentity(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.MarshalIndent(Error{"Not received input in a valid format"}, "", "  ")
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var idn FetchIdentity

		err = json.Unmarshal(body, &idn)
		if err != nil {
			jsonval, _ := json.MarshalIndent(Error{"Unable to unmarshal the entered input. Provide input in valid format"}, "", "  ")
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			updateIdentityResponse, err := idn.updateIdentity()
			if err != nil {
				fmt.Fprintf(rw, "%v\n", err)
			} else {
				jsonval, _ := json.MarshalIndent(updateIdentityResponse, "", " ")
				fmt.Fprintf(rw, "%v\n", "Identity was updated successfully, find the updated list below")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func deleteIdentity(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		jsonval, _ := json.MarshalIndent(Error{"Not received input in a valid format"}, "", "  ")
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	} else {
		var idn FetchIdentity

		err = json.Unmarshal(body, &idn)
		if err != nil {
			jsonval, _ := json.MarshalIndent(Error{"Unable to unmarshal the entered input. Provide input in valid format"}, "", "  ")
			fmt.Fprintf(rw, "%v\n", string(jsonval))
		} else {
			deleteIdentityResponse, err := idn.deleteIdentity()
			if err != nil {
				fmt.Fprintf(rw, "%v\n", err)
			} else {
				jsonval, _ := json.MarshalIndent(deleteIdentityResponse, "", "  ")
				fmt.Fprintf(rw, "%v\n", "Identity was deleted successfully, find the updated list below")
				fmt.Fprintf(rw, "%v\n", string(jsonval))
			}
		}
	}
}

func identity(rw http.ResponseWriter, req *http.Request) {
	idn := FetchIdentity{}
	identityResponse, err := idn.identity()
	if err != nil {
		fmt.Fprintf(rw, "%v\n", err)
	} else {
		jsonval, _ := json.MarshalIndent(identityResponse, "", "  ")
		fmt.Fprintf(rw, "%v\n", "List of identities found in database")
		fmt.Fprintf(rw, "%v\n", string(jsonval))
	}
}

// identity returns all the identities that are stored in the app so far.
func (i *FetchIdentity) identity() (map[int]Identity, error) {
	return identities, nil
}

// getIdentity retrives the identy that is seeked by the endpoint.
func (i *FetchIdentity) getIdentity() (Identity, error) {
	for _, name := range identities {
		if (name.FirstName == i.Identity.FirstName) || (name.SecondName == i.Identity.SecondName) || (name.Dob == i.Identity.Dob) {
			return name, nil
		}
	}
	return Identity{}, fmt.Errorf("Oops...!! Identity not found")
}

// createIdentity creates the identity and returns all the identities that has been managed by the app so far.
func (i *FetchIdentity) createIdentity() (Identity, error) {
	return i.Identity, nil
}

// updateIdentity updates the identity and returns the updated list of identities.
func (i *FetchIdentity) updateIdentity() (map[int]Identity, error) {
	length := len(identities)
	identities[(length + 1)] = Identity{FirstName: i.Identity.FirstName, SecondName: i.Identity.SecondName, Dob: i.Identity.Dob}
	return identities, nil
}

// deleteIdentity delets the identity and returns the updated list of identities after deletion.
func (i *FetchIdentity) deleteIdentity() (map[int]Identity, error) {
	for index, name := range identities {
		if (i.Identity.FirstName == name.FirstName) || (i.Identity.SecondName == name.SecondName) || (i.Identity.Dob == name.Dob) {
			delete(identities, index)
		}
		return identities, nil
	}

	return nil, fmt.Errorf("Oops...!! Unable to delete the Identity")
}
