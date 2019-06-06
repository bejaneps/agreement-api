// Package handlers implements all routers needed for client's application requests
package handlers

import (
	"fmt"
	"net/http"

	"github.com/bejaneps/agreement-api/crud"
	"github.com/gin-gonic/gin"
)

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
		errorHandler(err, http.StatusOK, c)
		return
	}

	/* 2. Create document in google drive */
	file, err := crud.CreateDocument(usr.Email, usr.DocTitle)
	if file == nil || err != nil {
		errorHandler(err, http.StatusOK, c)
		return
	}

	docURL := fmt.Sprintf("https://docs.google.com/document/d/%s/", file.Id)

	/* 3. Add user document info to DB */
	doc, err := crud.AddUserDoc(file.Id, file.Title, docURL, usr.Email)
	if err != nil {
		errorHandler(err, http.StatusOK, c)
		return
	}

	c.JSON(http.StatusOK, &response{
		Success: true,
		Err:     "",
		Data:    doc,
	})
}

// TemplateCreateHandler creates a document from a template for a user
func TemplateCreateHandler(c *gin.Context) {
	/* 1. Get user and document data and unpack it to variable */
	usr := &userCreateTemplate{}

	err := c.BindJSON(usr)
	if err != nil {
		errorHandler(err, http.StatusOK, c)
		return
	}

	/* 2. Create a document with contents of template */
	file, err := crud.CreateTemplate(usr.Email, usr.TemplateID, usr.DocTitle)
	if file == nil || err != nil {
		errorHandler(err, http.StatusOK, c)
		return
	}

	docURL := fmt.Sprintf("https://docs.google.com/document/d/%s/", file.Id)

	/* 3. Add user document info to DB */
	doc, err := crud.AddUserDoc(file.Id, file.Title, docURL, usr.Email)
	if err != nil {
		errorHandler(err, http.StatusOK, c)
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
		errorHandler(err, http.StatusOK, c)
		return
	}

	doc, err := crud.GetUserDoc(usrPerm.DocID)
	if err != nil {
		errorHandler(err, http.StatusOK, c)
		return
	}

	if doc.Signed1 == 0 && doc.Signed2 == 0 {
		err = crud.SetPermission(usrPerm.DocID, usrPerm.Email, "reader")
		if err != nil {
			errorHandler(err, http.StatusOK, c)
			return
		}

		doc, err = crud.AddUserSign(usrPerm.Email, usrPerm.DocID)
		if err != nil {
			errorHandler(err, http.StatusOK, c)
			return
		}
	} else if doc.Signed1 == 1 {
		lmd, err := crud.LastModifiedDate(usrPerm.DocID)
		if err != nil {
			errorHandler(err, http.StatusOK, c)
			return
		}

		if lmd > doc.DateSigned1 {
			_, err := crud.RemoveUserSign(doc.Owner1, usrPerm.DocID)
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}

			err = crud.SetPermission(usrPerm.DocID, usrPerm.Email, "reader")
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}

			err = crud.SetPermission(usrPerm.DocID, doc.Owner1, "writer")
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}

			doc, err = crud.AddUserSign(usrPerm.Email, usrPerm.DocID)
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}
		} else {
			err = crud.SetPermission(usrPerm.DocID, usrPerm.Email, "reader")
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}

			doc, err = crud.AddUserSign(usrPerm.Email, usrPerm.DocID)
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}
		}
	} else if doc.Signed2 == 1 {
		lmd, err := crud.LastModifiedDate(usrPerm.DocID)
		if err != nil {
			errorHandler(err, http.StatusOK, c)
			return
		}

		if lmd > doc.DateSigned2 {
			_, err := crud.RemoveUserSign(doc.Owner2, usrPerm.DocID)
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}

			err = crud.SetPermission(usrPerm.DocID, usrPerm.Email, "reader")
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}

			err = crud.SetPermission(usrPerm.DocID, doc.Owner2, "writer")
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}

			doc, err = crud.AddUserSign(usrPerm.Email, usrPerm.DocID)
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}
		} else {
			err = crud.SetPermission(usrPerm.DocID, usrPerm.Email, "reader")
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}

			doc, err = crud.AddUserSign(usrPerm.Email, usrPerm.DocID)
			if err != nil {
				errorHandler(err, http.StatusOK, c)
				return
			}
		}
	}

	c.JSON(http.StatusOK, &response{
		Success: true,
		Err:     "",
		Data:    doc,
	})
	return
}

// DocListHandler sents the list of documents that belong to user
func DocListHandler(c *gin.Context) {
	uEmail := &userEmail{}

	err := c.BindJSON(uEmail)
	if err != nil {
		errorHandler(err, http.StatusOK, c)
		return
	}

	if uEmail.Email == "" {
		c.JSON(http.StatusOK, &response{
			Success: false,
			Err:     "empty email",
			Data:    nil,
		})
	}

	docList, err := crud.GetUserDocList(uEmail.Email)
	if err != nil {
		errorHandler(err, http.StatusOK, c)
		return
	}

	c.JSON(http.StatusOK, &response{
		Success: true,
		Err:     "",
		Data:    docList,
	})
}
