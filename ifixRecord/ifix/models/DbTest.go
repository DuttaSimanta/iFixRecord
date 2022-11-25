package models

import (
	"fmt"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func APIcall() {
	runtime.GOMAXPROCS(8)
	start := time.Now()
	var wg sync.WaitGroup
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		//go Worker(i, &wg)
		go RecordWorker(i, &wg)
	}
	wg.Wait()
	time.Sleep(time.Second)
	end := time.Now()
	diff := end.Sub(start).Milliseconds()
	fmt.Printf("%d", diff-1000)

}

func Worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	t := strconv.Itoa(id)
	reqbd := entities.RecordcommonEntity{}
	reqbd.ClientID = 2
	reqbd.Mstorgnhirarchyid = 2
	reqbd.RecordID = 868
	reqbd.RecordstageID = 958
	reqbd.Termseq = 22
	reqbd.Recorddifftypeid = 2
	reqbd.Recorddiffid = 4
	reqbd.Usergroupid = 16
	reqbd.ForuserID = 12
	reqbd.Termvalue = "aaaa->" + t
	reqbd.Userid = 12
	insertID, _, err, _ := InsertRecordTermvalues(&reqbd)
	if err != nil {
		logger.Log.Println("Error is  ------> ", err)
	}
	logger.Log.Println("Insert id is  ------> ", insertID)
}

func RecordWorker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	t := strconv.Itoa(id)
	//logger.Log.Println("Error is  ------> ", t)
	reqbd := entities.RecordEntity{}
	reqbd.ClientID = 2
	reqbd.Mstorgnhirarchyid = 2
	reqbd.Recordname = "Test IM" + t
	reqbd.Recordesc = "Test IM" + t
	reqbd.Userid = 5
	reqbd.Usergroupid = 98
	reqbd.Originaluserid = 5
	reqbd.Originalusergroupid = 98
	reqbd.Workingcatlabelid = 34
	reqbd.Requestername = "Akash Jadhav"
	reqbd.Requesteremail = "jadhav.akash2@tcs.com"
	reqbd.Requestermobile = "8097938171"
	reqbd.Requesterlocation = "Yantra Park"
	reqbd.Source = "Call"
	reqbd.CreateduserID = 5
	reqbd.CreatedusergroupID = 98
	//reqbd.Lastlevelcatid	=
	reqbd.ParentID = 0

	addl := []entities.RecordAdditional{}
	addl1 := entities.RecordAdditional{}
	addl1.ID = 7
	addl1.Termsid = 20
	addl = append(addl, addl1)

	addl2 := entities.RecordAdditional{}
	addl2.ID = 9
	addl2.Termsid = 21
	addl2.Val = ""
	addl = append(addl, addl2)

	addl3 := entities.RecordAdditional{}
	addl3.ID = 8
	addl3.Termsid = 22
	addl3.Val = ""
	addl = append(addl, addl3)

	reqbd.Additionalfields = addl

	rdd := []entities.RecordData{}
	rdd1 := entities.RecordData{}
	rdd1.ID = 32
	rdd1.Val = 496
	rdd = append(rdd, rdd1)
	rdd2 := entities.RecordData{}
	rdd2.ID = 33
	rdd2.Val = 497
	rdd = append(rdd, rdd2)
	rdd3 := entities.RecordData{}
	rdd3.ID = 34
	rdd3.Val = 498
	rdd = append(rdd, rdd3)
	rdd4 := entities.RecordData{}
	rdd4.ID = 35
	rdd4.Val = 499
	rdd = append(rdd, rdd4)
	rdd5 := entities.RecordData{}
	rdd5.ID = 36
	rdd5.Val = 509
	rdd = append(rdd, rdd5)

	rss := []entities.RecordSet{}
	rss1 := entities.RecordSet{}
	rss1.ID = 1
	rss1.Type = rdd
	rss = append(rss, rss1)

	rss2 := entities.RecordSet{}
	rss2.ID = 2
	rss2.Val = 4
	rss = append(rss, rss2)

	rss3 := entities.RecordSet{}
	rss3.ID = 5
	rss3.Val = 473
	rss = append(rss, rss3)

	rss4 := entities.RecordSet{}
	rss4.ID = 3
	rss4.Val = 8
	rss = append(rss, rss4)
	reqbd.RecordSets = rss
	//logger.Log.Println("Error is  --22222222222222222----> ", reqbd)
	recordnumber, _, modelResponseError := CreateRecordModel(&reqbd)
	if modelResponseError != nil {
		logger.Log.Println("Error is  ------> ", modelResponseError)
	}
	logger.Log.Println("Insert id is  ------> ", recordnumber)

}
