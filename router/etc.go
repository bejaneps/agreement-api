// Package routers implements all routers needed for client's application requests
package routers

type response struct {
	Success bool        `json:"success"`
	Err     string      `json:"error"`
	Data    interface{} `json:"data"`
}

// userCreateDoc is a user info, needed when creating a document
type userCreateDoc struct {
	Email      string `json:"email,omitempty"`
	DocTitle   string `json:"doc_title"`
	TemplateID string `json:"template_id"`
}

// userCreateTemplate is a template info from which document will be created
type userCreateTemplate struct {
	Email      string `json:"email,omitempty"`
	TemplateID string `json:"template_id,omitempty"`
	DocTitle   string `json:"doc_title,omitempty"`
}

// userPermission is a struct representing user email and id for giving read write perm to this user
type userPermission struct {
	Email string `json:"email,omitempty"`
	DocID string `json:"doc_id,omitempty"`
}

// userEmail is a user email from a DB
type userEmail struct {
	Email string `json:"email"`
}
