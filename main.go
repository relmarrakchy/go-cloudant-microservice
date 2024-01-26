package main

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/IBM/go-sdk-core/v5/core"
)

func main() {
	//Create the authenticator
	authenticator := &core.IamAuthenticator{
		ApiKey: "---------------", //API_KEY provided in service credentials
	}

	options := &cloudantv1.CloudantV1Options{
		Authenticator: authenticator,
		//the URL provided in service credentials
		URL: "-------------------------------------",
	}

	myService, err := cloudantv1.NewCloudantV1(options)
	if err != nil {
		panic(err)
	}

	databaseName := "testing"

	//============ Creating a database======================================

	// putDatabaseResult, putDatabaseResponse, err := myService.PutDatabase(
	// 	myService.NewPutDatabaseOptions(databaseName),
	// )

	// if err != nil {
	// 	if putDatabaseResult != nil && putDatabaseResponse.StatusCode == 412 {
	// 		fmt.Printf("Cannot create \"%s\" database, it already exists. \n", databaseName)
	// 	} else {
	// 		panic(err)
	// 	}
	// }

	// if putDatabaseResponse != nil && putDatabaseResult != nil && *putDatabaseResult.Ok {
	// 	fmt.Printf("\"%s\" database created.\n", databaseName)
	// }

	//===============Adding a document=========================================
	documentID := "DOC01"

	docToAdd := cloudantv1.Document{
		ID: &documentID,
	}

	docToAdd.SetProperty("name", "Reda ELMARRAKCHY")
	docToAdd.SetProperty("age", 19)

	creationOptions := myService.NewPostDocumentOptions(
		databaseName,
	).SetDocument(&docToAdd)

	cerationResponse, _, err := myService.PostDocument(creationOptions)

	if err != nil {
		panic(err)
	}

	docToAdd.Rev = cerationResponse.Rev
	fmt.Println(docToAdd.ID)

	//=========================Reading a doc================================================
	doc, _, err := myService.GetDocument(myService.NewGetDocumentOptions(databaseName, "DOC01"))
	if err != nil {
		panic(err)
	}
	docBuffer, _ := json.MarshalIndent(doc, "", " ")

	fmt.Println(string(docBuffer))

	//===========================Updationg doc================================
	docToUp, fetchingResponse, err := myService.GetDocument(myService.NewGetDocumentOptions(databaseName, *doc.ID))
	if err != nil {
		if fetchingResponse.StatusCode == 404 {
			fmt.Printf("The doc with the id \"%s\" not exist!", *doc.ID)
		} else {
			fmt.Println("Internal server error !")
		}
	}

	if docToUp != nil {
		docToUp.SetProperty("Address", "Add Test")
		delete(docToUp.GetProperties(), "joined")

		updatedDoc, _, err := myService.PostDocument(myService.NewPostDocumentOptions(databaseName).SetDocument(docToUp))
		if err != nil {
			panic(err)
		}
		upBuff, _ := json.MarshalIndent(updatedDoc, "", "")
		fmt.Println(string(upBuff)) //The updated version return just the id and the rev value with a index ok = true to confirm the changes
	}

	//=============================Deleting a doc=========================================================================
	docToDel, fetchingResponse, err := myService.GetDocument(myService.NewGetDocumentOptions(databaseName, *doc.ID))
	if err != nil {
		if fetchingResponse.StatusCode == 404 {
			fmt.Printf("The doc with the id \"%s\" not exist!", *doc.ID)
		} else {
			fmt.Println("Internal server error !")
		}
	}

	if docToDel != nil {
		delResult, _, err := myService.DeleteDocument(myService.NewDeleteDocumentOptions(databaseName, *docToDel.ID).SetRev(*docToDel.Rev))
		if err != nil {
			fmt.Println("Internal server error !")
		}

		if *delResult.Ok {
			fmt.Println("Deletion process done !")
		}
	}

}
