package dbconfig

//var MASTER_URL = "http://20.204.29.18:8082/api"
//var NOTIFICATION_URL = "http://20.204.29.18:8089/iFIXNotification"

var MASTER_URL = "http://localhost:8082/api"
var NOTIFICATION_URL = "http://localhost:8089/iFIXNotification"
var FileUploadUrl = "https://iccmuat.ifixcloud.io/api/fileupload"

var DownloadFileURL ="https://iccmuat.ifixcloud.io/api/filedownload"
//var MASTER_URL = "http://10.5.2.6:8082/api"
//var NOTIFICATION_URL = "http://10.5.2.6:8089/iFIXNotification"

var TOKEN_EXPIRE_TIME = 3600
var JWT_KEY = []byte("ifix4@secret#token%^@")
