package v1

//ActCreated is an act indicates that create action was finished
type ActCreated struct {
	Created bool `json:"created"`
}


//ActUpdated is an act indicates that update action was finished
type ActUpdated struct {
	Updated bool `json:"updated"`
}

//ActDeleted is an act indicates that delete action was finished
type ActDeleted struct {
	Deleted bool `json:"deleted"`
}
