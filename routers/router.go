// Package routers implements all routers needed for client's application requests
package routers

import (
	"fmt"
	"net/http"

	"google.golang.org/api/drive/v2"

	"github.com/bejaneps/agreement-api/crud"
	"github.com/gin-gonic/gin"
)

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

// userPermission is a struct representing user email and id for giving read write perm to this user
type userPermission struct {
	Email string `json:"email,omitempty"`
	DocID string `json:"doc_id,omitempty"`
}

// userEmail is a user email from a DB
type userEmail struct {
	Email string `json:"email"`
}

func errorHandler(err error, errCode uint, c *gin.Context) {
	c.AbortWithStatusJSON(int(errCode), &response{
		Success: false,
		Err:     err.Error(),
		Data:    "",
	})
}

// DocCreateHandler handles upcoming post requests for creation of document
func DocCreateHandler(c *gin.Context) {
	/* 1. Get user and document data and unpack it to variable */
	usr := &userCreateDoc{}

	err := c.BindJSON(usr)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	/* 2. Create a document from data */
	file := &drive.File{}

	if usr.TemplateID != "" {
		file, err = crud.CreateTemplate(usr.Email, usr.TemplateID)
		if file == nil || err != nil {
			errorHandler(err, http.StatusInternalServerError, c)
			return
		}
	} else {
		file, err = crud.CreateDocument(usr.Email, usr.DocTitle)
		if file == nil || err != nil {
			errorHandler(err, http.StatusInternalServerError, c)
			return
		}
	}

	docURL := fmt.Sprintf("https://docs.google.com/document/d/%s/", file.Id)

	/* 3. Add user document info to DB */
	doc, err := crud.AddUserDoc(file.Id, file.Title, docURL, usr.Email)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	c.JSON(http.StatusOK, &response{
		Success: true,
		Err:     "",
		Data:    doc,
	})
}

// DocPermHandler gives write permission to a user
func DocPermHandler(c *gin.Context) {
	usrPerm := &userPermission{}

	err := c.BindJSON(usrPerm)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	err = crud.SetPermission(usrPerm.DocID, usrPerm.Email, "writer")
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	doc, err := crud.AddDocOwner(usrPerm.Email, usrPerm.DocID)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	c.JSON(http.StatusOK, &response{
		Success: true,
		Err:     "",
		Data:    doc,
	})
}

// DocSignHandler removes write permission from a user
func DocSignHandler(c *gin.Context) {
	usrPerm := &userPermission{}

	err := c.BindJSON(usrPerm)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	sign, err := crud.CheckUserChanges(usrPerm.DocID, usrPerm.Email)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	if sign {
		doc, err := crud.RemoveUserSign(usrPerm.Email, usrPerm.DocID)
		if err != nil {
			errorHandler(err, http.StatusInternalServerError, c)
			return
		}

		var usrDoc *crud.Document
		if doc.Signed1 == 0 {
			err = crud.SetPermission(usrPerm.DocID, doc.Owner1, "writer")
			if err != nil {
				errorHandler(err, http.StatusInternalServerError, c)
				return
			}

			err = crud.SetPermission(usrPerm.DocID, doc.Owner2, "reader")
			if err != nil {
				errorHandler(err, http.StatusInternalServerError, c)
				return
			}

			usrDoc, err = crud.AddUserSign(doc.Owner2, usrPerm.DocID)
			if err != nil {
				errorHandler(err, http.StatusInternalServerError, c)
				return
			}
		} else if doc.Signed2 == 0 {
			err = crud.SetPermission(usrPerm.DocID, doc.Owner2, "writer")
			if err != nil {
				errorHandler(err, http.StatusInternalServerError, c)
				return
			}

			err = crud.SetPermission(usrPerm.DocID, doc.Owner1, "reader")
			if err != nil {
				errorHandler(err, http.StatusInternalServerError, c)
				return
			}

			usrDoc, err = crud.AddUserSign(doc.Owner1, usrPerm.DocID)
			if err != nil {
				errorHandler(err, http.StatusInternalServerError, c)
				return
			}
		}

		c.JSON(http.StatusOK, &response{
			Success: true,
			Err:     "",
			Data:    usrDoc,
		})
		return
	}

	err = crud.SetPermission(usrPerm.DocID, usrPerm.Email, "reader")
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	file, err := crud.AddUserSign(usrPerm.Email, usrPerm.DocID)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	c.JSON(http.StatusOK, &response{
		Success: true,
		Err:     "",
		Data:    file,
	})
	return
}

// DocListHandler sends the list of documents that belong to user
func DocListHandler(c *gin.Context) {
	uEmail := &userEmail{}

	err := c.BindJSON(uEmail)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	docList, err := crud.GetUserDocList(uEmail.Email)
	if err != nil {
		errorHandler(err, http.StatusInternalServerError, c)
		return
	}

	c.JSON(http.StatusOK, &response{
		Success: true,
		Err:     "",
		Data:    docList,
	})
}
