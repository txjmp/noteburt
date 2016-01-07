package data

import "log"

/*
 A note's previd value is the id of the note located before it in sorted order.
 The 1st note's previd = Zeroid (all zeros id)

 User specified note ordering is enabled by previds.
 In the database, a bucket for each tab contains the previd value for each note in the tab
 	key=noteid, val=previd
 Storing previds separately prevents reading/writing full note records on changes to order.
 One note change may affect previd of multiple records.

 On the client side, a noteOrder map (javascript object) is created whenever the notes need
 to be retrieved in order. The noteOrder map key=previd, val=noteid. (see screen2.js)

 Each tab has a previdMap, tabid is unique across books, so bookid is not part of key.
 tabPrevids is a map of all the previdMaps
 The tab's previdMap is updated on: addNote, deleteNote, positionNote requests
 GetTabNotes uses the previdMap to populate the Previd field in each note

 Changes may affect more than the specified noteid.
 		Each previdMap method returns orderChanges containing all affected entries
		Database updates use orderChanges to keep the db in sync with each previdMap

 The database previd buckets are read when a book is 1st requested (see getbooktabs.go).
 Read requests are then handled using the previdMaps, not the database.
 Changes are made to both the previdMaps and the database.

 Position changes are accomplished with delete(from current position) add(to new position).
*/
type previdMap map[string]string    // key=noteid, val=previd (in sorted order)
var tabPrevids map[string]previdMap // key=tabid

func init() {
	tabPrevids = make(map[string]previdMap)
}

func (this previdMap) add(newid, newPrevid string) map[string]string {
	orderChanges := make(map[string]string)

	if len(this) == 0 && newPrevid != Zeroid {
		log.Fatal("previdMap.add, empty map, newPrevid parm not = zeroid")
	}
	// entry currently following newPrevid will now follow newid
	for noteid, notePrevid := range this {
		if notePrevid == newPrevid {
			this[noteid] = newid
			orderChanges[noteid] = newid
			break
		}
	}
	this[newid] = newPrevid
	orderChanges[newid] = newPrevid
	return orderChanges
}

func (this previdMap) delete(deleteid string) map[string]string {
	orderChanges := make(map[string]string)

	// entry currently following deleteid will now follow entry before deleteid
	for noteid, notePrevid := range this {
		if notePrevid == deleteid {
			this[noteid] = this[deleteid]
			orderChanges[noteid] = this[deleteid]
			break
		}
	}
	delete(this, deleteid)
	orderChanges[deleteid] = "delete"
	return orderChanges
}
