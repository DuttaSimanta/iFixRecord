//***************************//
// Package Name: Router
// Date Of Creation: 09/01/2021
// Authour Name: Moloy Mondal
// History: N/A
// Synopsis: Routes are defined here
// Functions: PostRequestHandler,main
// API: /addmstclient Inputs:  {"code": <string>,"name":<string>,"description":<string>,"ClientAuditFlg":<single character>}
// API: /getmstclientbyid Inputs:  {"id":<Number>}
// API: /delmstclientbyid Inputs:  {"id":<Number>}
// API: /updatemstclientbyid Inputs:  {"id":<Number>,code": <string>,"name":<string>,"description":<string>,"ClientAuditFlg":<single character>}
// API: /getmstclientall Inputs:  N/A
//
// API: /addmstcountry Inputs:  {"code": <string>,"name":<string>,"description":<string>,"ClientAuditFlg":<single character>}
// API: /getmstcountrybyid Inputs:  {"id":<Number>}
// API: /getmstcountryall Inputs:  {"id":<Number>}
// API: /delmstcountrybyid Inputs:  {"id":<Number>,code": <string>,"name":<string>,"description":<string>,"ClientAuditFlg":<single character>}
// API: /updatemstcountrybyid Inputs:  N/A
// Global Variable: N/A
// Version: 1.0.0
//***************************//

package router

import (
	"ifixRecord/ifix/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"404 Not Found",
		"POST",
		"/",
		ThrowBlankResponse,
	},
	Route{
		"heartbeat",
		"POST",
		"/heartbeat",
		handlers.GetifixRecordDockerStatus,
	},
	Route{
		"insertbulkrecordtermvalues",
		"POST",
		"/insertbulkrecordtermvalues",
		handlers.InsertBulkRecordTermvalues,
	},
	Route{
		"updatesladuetime",
		"POST",
		"/updatesladuetime",
		handlers.UpdateSlaDueTime,
	},
	//The below one is added for external api 
	Route{
        "extarnalrecordattachment",
        "POST",
        "/extarnalrecordattachment",
        handlers.UpdateExternalRecordAttachment,
    },
    Route{
		"ifixebonding",
		"POST",
		"/ifixebonding",
		handlers.EbondingTicket,
	},

	Route{
		"createrecord",
		"POST",
		"/createrecord",
		handlers.AddRecordData,
	},

	Route{
		"getrecorddata",
		"POST",
		"/getrecorddata",
		handlers.GetRecordcreatedata,
	},
	Route{
		"getrecordtypedata",
		"POST",
		"/getrecordtypedata",
		handlers.GetRecordtypedata,
	},
	Route{
		"getrecordcatchilddata",
		"POST",
		"/getrecordcatchilddata",
		handlers.GetRecordcatchild,
	},
	Route{
		"getrecordpriority",
		"POST",
		"/getrecordpriority",
		handlers.GetRecordprioritydata,
	}, Route{
		"getadditionalfields",
		"POST",
		"/getadditionalfields",
		handlers.GetAdditionalFields,
	}, Route{
		"getdynamicqueryresult",
		"POST",
		"/getdynamicqueryresult",
		handlers.GetDynamicQueryResult,
	},
	Route{
		"getdynamicquerycountresult",
		"POST",
		"/getdynamicquerycountresult",
		handlers.GetDynamicCountQueryResult,
	}, Route{
		"getrecorddetails",
		"POST",
		"/getrecorddetails",
		handlers.GetRecordDetails,
	}, Route{
		"getrecordcategory",
		"POST",
		"/getrecordcategory",
		handlers.GetRecordCatDetails,
	},
	Route{
		"inserttermvalue",
		"POST",
		"/inserttermvalue",
		handlers.InsertTermvalue,
	},
	Route{
		"gettermvalues",
		"POST",
		"/gettermvalues",
		handlers.GetAllTermvalues,
	},
	Route{
		"gettermvaluesbytermid",
		"POST",
		"/gettermvaluesbytermid",
		handlers.GetTermvaluesbytermid,
	},
	Route{
		"getcommontermnames",
		"POST",
		"/getcommontermnames",
		handlers.GetTermnames,
	},
	Route{
		"getcommontermnamesbystate",
		"POST",
		"/getcommontermnamesbystate",
		handlers.GetTermnamesbystate,
	},
	Route{
		"updaterecordstatus",
		"POST",
		"/updaterecordstatus",
		handlers.Updaterecordstatusvalue,
	},
	Route{
		"insertmultipletermvalue",
		"POST",
		"/insertmultipletermvalue",
		handlers.InsertMultipleTermvalue,
	}, Route{
		"getrecordassetbyid",
		"POST",
		"/getrecordassetbyid",
		handlers.GetRecordAssetByID,
	}, Route{
		"addassetwithrecord",
		"POST",
		"/addassetwithrecord",
		handlers.AddAssetWithRecord,
	},
	Route{
		"deleteassetfromrecord",
		"POST",
		"/deleteassetfromrecord",
		handlers.DeleteAssetFromRecord,
	}, Route{
		"getassettypesbyrecordid",
		"POST",
		"/getassettypesbyrecordid",
		handlers.GetAssetTypeByRecordID,
	},

	Route{
		"getSlatabvalues",
		"POST",
		"/getSlatabvalues",
		handlers.GetSLATabvalues,
	},

	Route{
		"getSlaresolutionremain",
		"POST",
		"/getSlaresolutionremain",
		handlers.GetSLAResolution,
	}, Route{
		"getmiscdatabyrecordid",
		"POST",
		"/getmiscdatabyrecordid",
		handlers.GetMiscDataByRecordID,
	}, Route{
		"getrecorddetailsbyno",
		"POST",
		"/getrecorddetailsbyno",
		handlers.GetRecordDetailsByNo,
	}, Route{
		"savechildrecord",
		"POST",
		"/savechildrecord",
		handlers.SaveChildRecord,
	}, Route{
		"getchildrecordbyparent",
		"POST",
		"/getchildrecordbyparent",
		handlers.GetChildRecordsBYParentID,
	},
	Route{
		"updatepriority",
		"POST",
		"/updatepriority",
		handlers.Updaterecordpriority,
	},

	Route{
		"recordcount",
		"POST",
		"/recordcount",
		handlers.GetRecordcount,
	},
	Route{
		"recentrecords",
		"POST",
		"/recentrecords",
		handlers.Getrecentrecords,
	},

	Route{
		"recordlogs",
		"POST",
		"/recordlogs",
		handlers.GetRecordlogs,
	}, Route{
		"getadditionalfieldsbytypecat",
		"POST",
		"/getadditionalfieldsbytypecat",
		handlers.GetAdditionalFieldsBYTypeCat,
	}, Route{
		"removechildrecord",
		"POST",
		"/removechildrecord",
		handlers.RemoveChildRecord,
	},
	Route{
		"frequentrecords",
		"POST",
		"/frequentrecords",
		handlers.Getfrequentissues,
	},

	Route{
		"getparentrecordid",
		"POST",
		"/getparentrecordid",
		handlers.GetParentrecord,
	},

	Route{
		"updatecategory",
		"POST",
		"/updatecategory",
		handlers.Updaterecordcategory,
	},

	Route{
		"getactivitynms",
		"POST",
		"/getactivitynms",
		handlers.GetActivitymstnames,
	},
	Route{
		"newactivitylogs",
		"POST",
		"/newactivitylogs",
		handlers.GetNewActivitylogs,
	},
	Route{
		"searchactivitylogs",
		"POST",
		"/searchactivitylogs",
		handlers.Activitylogsearch,
	},
	Route{
		"updateadditionalfields",
		"POST",
		"/updateadditionalfields",
		handlers.Updateadditionalfields,
	},
	Route{
		"pendingvendortermsvalue",
		"POST",
		"/pendingvendortermsvalue",
		handlers.GetPendingstatustermvalue,
	},
	Route{
		"getattachedfiles",
		"POST",
		"/getattachedfiles",
		handlers.GetAttachmentfiles,
	},

	Route{
		"updatedoccount",
		"POST",
		"/updatedoccount",
		handlers.Updatedocumentcount,
	},

	Route{
		"gettermnamebyseq",
		"POST",
		"/gettermnamebyseq",
		handlers.GetTermnamesbyseq,
	}, Route{
		"getassetbyrecordidnfieldname",
		"POST",
		"/getassetbyrecordidnfieldname",
		handlers.GetAssetfieldSpecificDataBYRecordID,
	},

	Route{
		"customervisiblecomment",
		"POST",
		"/customervisiblecomment",
		handlers.Customervisiblecomment,
	},

	Route{
		"deleteattachment",
		"POST",
		"/deleteattachment",
		handlers.Deleteattchfile,
	},

	Route{
		"gettermvaluebyseq",
		"POST",
		"/gettermvaluebyseq",
		handlers.GetTermvaluebyseq,
	},

	Route{
		"getparencollaborationtchildlogs",
		"POST",
		"/getparencollaborationtchildlogs",
		handlers.Parentchildcollaborationlogs,
	},
	Route{
		"addparentfromchild",
		"POST",
		"/addparentfromchild",
		handlers.Addparentfromchild,
	},

	Route{
		"getparentrecorddetails",
		"POST",
		"/getparentrecorddetails",
		handlers.GetParentRecordDetailsByNo,
	},
	Route{
		"childsearchcriteria",
		"POST",
		"/childsearchcriteria",
		handlers.ChildRecordSearchCriteria,
	}, Route{
		"getassetdetailsbyid",
		"POST",
		"/getassetdetailsbyid",
		handlers.GetAssetDetailsByAssetID,
	}, Route{
		"getallassettypendetailsbyrecordid",
		"POST",
		"/getallassettypendetailsbyrecordid",
		handlers.GetAllAssetTypeNDetailsByRecordID,
	},
	Route{
		"createexternalrecord",
		"POST",
		"/createexternalrecord",
		handlers.CreateExternalRecord,
	},

	Route{
		"getrecordvaluesbyno",
		"POST",
		"/getrecordvaluesbyno",
		handlers.GetExternalRecordDetailsNo,
	},
	Route{
		"getrecordvaluesbydate",
		"POST",
		"/getrecordvaluesbydate",
		handlers.GetExternalRecordDetailsbyDate,
	},

	Route{
		"externalrecordstatusupdate",
		"POST",
		"/externalrecordstatusupdate",
		handlers.CreateExternalRecordStatusUpdate,
	},

	Route{
		"externalrecordassignegrpupdate",
		"POST",
		"/externalrecordassignegrpupdate",
		handlers.CreateExternalRecordGrpUpdate,
	},

	Route{
		"externalrecordassigneuserupdate",
		"POST",
		"/externalrecordassigneuserupdate",
		handlers.CreateExternalRecordUserUpdate,
	},

	Route{
		"externalrecordcommentupdate",
		"POST",
		"/externalrecordcommentupdate",
		handlers.ExternalRecordCommentUpdate,
	}, Route{
		"getcategorybylastid",
		"POST",
		"/getcategorybylastid",
		handlers.GetCatByLastID,
	},Route{
		"recordgridresult",
		"POST",
		"/recordgridresult",
		handlers.RecordGridResult,
	},
	Route{
		"getrecordid",
		"POST",
		"/getrecordid",
		handlers.GetRecordID,
	},Route{
		"recordfilteradd",
		"POST",
		"/recordfilteradd",
		handlers.RecordFilterAdd,
	}, Route{
		"recordfilterlist",
		"POST",
		"/recordfilterlist",
		handlers.RecordFilterList,
	}, Route{
		"recordfilterdelete",
		"POST",
		"/recordfilterdelete",
		handlers.RecordFilterDelete,
	}, Route{
		"recordfullresult",
		"POST",
		"/recordfullresult",
		handlers.RecordGridResultOnly,
	},
	Route{
		"gettabterms",
		"POST",
		"/gettabterms",
		handlers.GetTabTermnames,
	},

	Route{
		"gettabtermvalues",
		"POST",
		"/gettabtermvalues",
		handlers.GetTabTermvalues,
	},
	Route{
		"removerecordlink",
		"POST",
		"/removerecordlink",
		handlers.Removelinkrecord,
	},
	Route{
		"getlinkrecorddetails",
		"POST",
		"/getlinkrecorddetails",
		handlers.GetLinkRecorddetails,
	},
	Route{
		"saverecordlink",
		"POST",
		"/saverecordlink",
		handlers.Savelinkrecord,
	},

	Route{
		"getparentrecordinfo",
		"POST",
		"/getparentrecordinfo",
		handlers.GetParentRecordInfo,
	},
	Route{
		"updaterecordasset",
		"POST",
		"/updaterecordasset",
		handlers.UpdateRecordAsset,
	},
	Route{
		"Inserrecordasset",
		"POST",
		"/insertrecordasset",
		handlers.InsertRecordAsset,
	},
	Route{
		"getAssetHistroyByAssetID",
		"POST",
		"/getassethistorybyassetid",
		handlers.GetAssetHistroyByAssetID,
	},
	Route{
		"getrecorddetailsbynoforlinkrecord",
		"POST",
		"/getrecorddetailsbynoforlinkrecord",
		handlers.GetRecordDetailsByNoForlinkrecord,
	},
	Route{
		"getadditionalinfobasedoncategory",
		"POST",
		"/getadditionalinfobasedoncategory",
		handlers.GetAdditionalInfoBasedonCategory,
	},

	Route{
		"getsladuetimecalculate",
		"POST",
		"/getsladuetimecalculate",
		handlers.GetSLADuetimeCalculation,
	},

	Route{
		"getrecordpermissionbyno",
		"POST",
		"/getrecordpermissionbyno",
		handlers.GetRecordAccesspermissionByNo,
	},

	Route{
		"getparentrecordidforIM",
		"POST",
		"/getparentrecordidforIM",
		handlers.GetParentrecordForIM,
	},

	Route{
		"recordfilterupdate",
		"POST",
		"/recordfilterupdate",
		handlers.RecordFilterUpdate,
	},

	Route{
		"updatevendorticketid",
		"POST",
		"/updatevendorticketid",
		handlers.UpdateVendorTickeID,
	},
}



