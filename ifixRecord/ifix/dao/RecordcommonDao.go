package dao

import (
	"database/sql"
	"errors"
	"ifixRecord/ifix/entities"
	"ifixRecord/ifix/logger"
)

var commontabvalues = "SELECT a.id as ID,a.recordtermid as TermID,COALESCE(a.recordtrackvalue,'NA') as Termvalue,COALESCE(a.foruserid,'NA') as ForuserID,b.termname as Tername,COALESCE(c.name,'') recorddiffname,COALESCE(CONCAT(d.firstname,' ',d.lastname),'NA') Username,a.createddate Createddate  FROM mstsupportgrptermmap e,mstgroupmember f,mstclientuser d,trnreordtracking a JOIN mstrecordterms b ON a.referenceid IS NULL AND a.recordtermid=b.id LEFT JOIN (SELECT mststateterm.clientid,mststateterm.mstorgnhirarchyid,mststateterm.recordtermid,mstrecorddifferentiation.id,mstrecorddifferentiation.name FROM mststateterm JOIN mstrecorddifferentiation ON mststateterm.recorddiffid=mstrecorddifferentiation.id AND mststateterm.recorddifftypeid=mstrecorddifferentiation.recorddifftypeid AND mststateterm.clientid=mstrecorddifferentiation.clientid AND mststateterm.mstorgnhirarchyid=mstrecorddifferentiation.mstorgnhirarchyid JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid=mstrecorddifferentiationtype.id WHERE mststateterm.activeflg=1 AND mststateterm.deleteflg=0 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0 AND mstrecorddifferentiationtype.seqno=2) c ON c.clientid=a.clientid AND c.mstorgnhirarchyid=a.mstorgnhirarchyid AND c.recordtermid=a.recordtermid WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.deleteflg=0 AND a.activeflg=1 and a.createdbyid=d.id and f.userid=? AND f.groupid=e.grpid AND e.termid=b.id AND e.deleteflg=0 AND e.activeflg=1"
var termvalueagainsttermid = "SELECT a.id as ID,a.recordtermid as TermID,COALESCE(a.recordtrackvalue,'NA') as Termvalue,COALESCE(a.foruserid,'NA') as ForuserID,b.termname as Tername,COALESCE(c.name,'') recorddiffname,COALESCE(CONCAT(d.firstname,' ',d.lastname),'NA') Username,a.createddate Createddate  FROM mstclientuser d,trnreordtracking a JOIN mstrecordterms b ON a.referenceid IS NULL AND a.recordtermid=b.id LEFT JOIN (SELECT mststateterm.clientid,mststateterm.mstorgnhirarchyid,mststateterm.recordtermid,mstrecorddifferentiation.id,mstrecorddifferentiation.name FROM mststateterm JOIN mstrecorddifferentiation ON mststateterm.recorddiffid=mstrecorddifferentiation.id AND mststateterm.recorddifftypeid=mstrecorddifferentiation.recorddifftypeid AND mststateterm.clientid=mstrecorddifferentiation.clientid AND mststateterm.mstorgnhirarchyid=mstrecorddifferentiation.mstorgnhirarchyid JOIN mstrecorddifferentiationtype ON mstrecorddifferentiation.recorddifftypeid=mstrecorddifferentiationtype.id WHERE mststateterm.activeflg=1 AND mststateterm.deleteflg=0 AND mstrecorddifferentiation.activeflg=1 AND mstrecorddifferentiation.deleteflg=0 AND mstrecorddifferentiationtype.seqno=2) c ON c.clientid=a.clientid AND c.mstorgnhirarchyid=a.mstorgnhirarchyid AND c.recordtermid=a.recordtermid WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.deleteflg=0 AND a.activeflg=1 and a.createdbyid=d.id AND a.recordtermid=?"
var termnames = "select distinct a.id as ID,a.termname as Termname,b.recordtermvalue as Recordtermvalue,b.iscompulsory as Iscompulsory,c.termtypename as Termtypename,a.termtypeid as Termtypeid from mstrecordterms a,mststateterm b,msttermtype c,mstsupportgrptermmap d,mstgroupmember e where b.clientid=? and b.mstorgnhirarchyid=? and b.recorddifftypeid=? and b.recorddiffid=? and a.id=b.recordtermid and a.termtypeid=c.id AND a.deleteflg=0 AND a.activeflg=1 AND e.userid=? AND e.groupid=d.grpid AND d.termid=a.id AND d.deleteflg=0 AND d.activeflg=1 AND b.deleteflg = 0 AND b.activeflg = 1 order by a.seq"

//var termnamesbystate = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,b.iscompulsory as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid from mststateterm a,mststateterm b,mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstgroupmember f where a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and b.recorddifftypeid=? and b.recorddiffid=? and a.recordtermid=b.recordtermid and b.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and b.deleteflg=0 and b.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND f.userid=? AND f.groupid=e.grpid AND e.termid=c.id AND e.deleteflg=0 and e.activeflg=1"
var termnamesbystate = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,b.iscompulsory as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,c.seq from mststateterm a,mststateterm b,mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstgroupmember f where a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and b.recorddifftypeid=? and b.recorddiffid=? and a.recordtermid=b.recordtermid and b.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and b.deleteflg=0 and b.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND f.userid=? AND f.groupid=e.grpid AND e.termid=c.id AND e.deleteflg=0 and e.activeflg=1 AND c.id not in (SELECT recordtermid FROM mstrecordfield WHERE clientid=? AND mstorgnhirarchyid=? AND activeflg=1 AND deleteflg=0) order by c.seq"

var getchildids = "select childrecordid from mstparentchildmap where parentrecordid=? and clientid=? and mstorgnhirarchyid=? and recorddifftypeid=? and recorddiffid=? and deleteflg=0 and activeflg=1 and isattached='Y'"

//var getchildids = "SELECT a.recordid FROM maprecordtorecorddifferentiation a, mstrecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddifftypeid=3 AND a.islatest=1 AND a.recorddiffid=b.id AND b.deleteflg=0 and b.activeflg=1 AND b.seqno NOT IN (3,11,14) AND a.recordid IN (select childrecordid from mstparentchildmap where parentrecordid=? and clientid=? and mstorgnhirarchyid=? and recorddifftypeid=? and recorddiffid=? and deleteflg=0 and activeflg=1 and isattached='Y')"
var termreleation = "select id as ID from mstrcordtremswisereleationconfig where clientid=? AND mstorgnhirarchyid=? AND recorddifftypeid=? AND recorddiffid=? AND termsid=? and deleteflg=0 and activeflg=1"

var pcount = "SELECT count(id) count FROM maprecordtorecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND recorddifftypeid=5"
var pvacount = "SELECT count(b.id) count FROM mstrecorddifferentiation a,maprecordtorecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddifftypeid=3 AND a.seqno=4 AND a.id = b.recorddiffid AND b.recordid=?"

//var recent = "SELECT a.id,a.code,a.recordtitle,a.createdatetime,c.name,c.seqno FROM trnrecord a,maprecordtorecorddifferentiation b,mstrecorddifferentiation c,mstrecordactivitylogs d where (a.userid=? OR a.originaluserid=?) AND d.recordid = a.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.id=b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND b.recorddiffid=c.id and d.id in (select max(x.id) from mstrecordactivitylogs x, trnrecord y where (y.userid=? OR y.originaluserid=?) AND x.recordid = y.id AND y.clientid=? AND y.mstorgnhirarchyid=? AND y.deleteflg=0 AND y.activeflg=1 group by y.id) group by d.recordid order by d.createddate desc limit 5"
var recent = "SELECT a.id,a.code,a.recordtitle,a.createdatetime,c.name,c.seqno,a.mstorgnhirarchyid FROM trnrecord a,maprecordtorecorddifferentiation b,mstrecorddifferentiation c,mstrecordactivitylogs d where (a.userid=? OR a.originaluserid=?) AND d.recordid = a.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.id=b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND b.recorddiffid=c.id and d.id in (select max(x.id) from mstrecordactivitylogs x, trnrecord y,maprecordtorecorddifferentiation c,mstrecorddifferentiation d where (y.userid=? OR y.originaluserid=?) AND x.recordid = y.id AND y.clientid=? AND y.mstorgnhirarchyid=? AND y.deleteflg=0 AND y.activeflg=1 AND y.id=c.recordid AND c.recorddifftypeid=2 AND c.islatest=1 AND c.recorddiffid=d.id AND d.seqno in (1,2) group by y.id) group by d.recordid order by d.createddate desc limit 5"

//var resolverrecent = "SELECT a.id,a.code,a.recordtitle,a.createdatetime,c.name,c.seqno FROM trnrecord a,maprecordtorecorddifferentiation b,mstrecorddifferentiation c,mstrecordactivitylogs d where d.createdid=? AND d.recordid = a.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.id=b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND b.recorddiffid=c.id and d.id in (select max(x.id) from mstrecordactivitylogs x, trnrecord y where x.createdid=? AND x.recordid = y.id AND y.clientid=? AND y.mstorgnhirarchyid=? AND y.deleteflg=0 AND y.activeflg=1 group by y.id) group by d.recordid order by d.createddate desc limit 10"
var recordlogs = "SELECT a.recordid, COALESCE(a.logValue,'') logvalue,a.createddate,concat(c.firstname,' ',c.lastname) name,b.activitydesc,d.supportgroupname FROM mstrecordactivitylogs a,mstrecordactivitymst b,mstclientuser c,mstclientsupportgroup d WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.activityseqno = b.seqno AND a.createdid = c.id AND a.createdgrpid=d.id order by a.id desc"
var frequent = "SELECT a.mstorgnhirarchyid,a.recorddiffid ,count(distinct a.recordid) cnt,c.parentcategoryids,c.parentcategorynames,b.fromrecorddifftypeid,b.fromrecorddiffid,c.seqno FROM maprecordtorecorddifferentiation a, mstrecordtype b,mstrecorddifferentiation c WHERE b.torecorddiffid = a.recorddiffid AND b.clientid=? AND b.mstorgnhirarchyid=? AND b.deleteflg=0 AND b.activeflg=1 AND b.fromrecorddifftypeid=? AND b.fromrecorddiffid=? AND b.torecorddifftypeid in (SELECT max(a.torecorddifftypeid ) FROM mstrecordtype a, mstrecorddifferentiationtype b WHERE a.torecorddifftypeid = b.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.fromrecorddifftypeid=? AND a.fromrecorddiffid=? AND b.parentid=1) AND a.recorddiffid=c.id  group by a.recorddiffid order by cnt desc limit 5"

//var frequentresolver = "SELECT a.recorddiffid,COUNT(DISTINCT a.recordid) cnt,c.parentcategoryids,c.parentcategorynames,b.fromrecorddifftypeid,b.fromrecorddiffid,c.seqno FROM maprecordtorecorddifferentiation a,mstrecordtype b,mstrecorddifferentiation c WHERE b.torecorddiffid = a.recorddiffid AND b.clientid = ? AND b.mstorgnhirarchyid in (SELECT distinct mstorgnhirarchyid FROM mstgroupmember where userid=?) AND b.deleteflg = 0 AND b.activeflg = 1 AND b.fromrecorddifftypeid = 2 AND b.fromrecorddiffid IN (SELECT distinct a.id FROM mstrecorddifferentiation a,mstgroupmember b WHERE b.userid=? AND b.mstorgnhirarchyid =a.mstorgnhirarchyid AND a.recorddifftypeid=2 AND a.seqno=1) AND b.torecorddifftypeid IN (SELECT MAX(a.torecorddifftypeid) FROM mstrecordtype a, mstrecorddifferentiationtype b WHERE a.torecorddifftypeid = b.id AND a.clientid = ? AND a.mstorgnhirarchyid in (SELECT distinct mstorgnhirarchyid FROM mstgroupmember where userid=?) AND a.deleteflg = 0 AND a.activeflg = 1 AND a.fromrecorddifftypeid = 2 AND a.fromrecorddiffid IN (SELECT distinct a.id FROM mstrecorddifferentiation a,mstgroupmember b WHERE b.userid=? AND b.mstorgnhirarchyid =a.mstorgnhirarchyid AND a.recorddifftypeid=2 AND a.seqno=1) AND b.parentid = 1 group by a.fromrecorddiffid) AND a.recorddiffid = c.id AND FIND_IN_SET(SUBSTRING_INDEX(SUBSTRING_INDEX(c.parentcategoryids, '->', 3), '->',- 1), (SELECT distinct GROUP_CONCAT(b.recorddiffid) FROM trnrecord a,maprecordtorecorddifferentiation b where (a.originalusergroupid IN (SELECT distinct groupid FROM mstgroupmember where userid=?) OR a.usergroupid IN (SELECT distinct groupid FROM mstgroupmember where userid=?)) AND a.id = b.recordid AND b.islatest=1 AND b.isworking=1)) GROUP BY a.recorddiffid order by cnt desc limit 10"
//var frequentresolver = "SELECT a.mstorgnhirarchyid,a.recorddiffid,COUNT(DISTINCT a.recordid) cnt,c.parentcategoryids,c.parentcategorynames,b.fromrecorddifftypeid,b.fromrecorddiffid,c.seqno FROM maprecordtorecorddifferentiation a,mstrecordtype b,mstrecorddifferentiation c WHERE b.torecorddiffid = a.recorddiffid AND b.clientid = ? AND b.mstorgnhirarchyid in (SELECT distinct mstorgnhirarchyid FROM mstgroupmember where userid=?) AND b.deleteflg = 0 AND b.activeflg = 1 AND b.fromrecorddifftypeid = 2 AND b.fromrecorddiffid IN (SELECT distinct a.id FROM mstrecorddifferentiation a,mstgroupmember b WHERE b.userid=? AND b.mstorgnhirarchyid =a.mstorgnhirarchyid AND a.recorddifftypeid=2 AND a.seqno=1) AND b.torecorddifftypeid IN (SELECT MAX(a.torecorddifftypeid) FROM mstrecordtype a, mstrecorddifferentiationtype b WHERE a.torecorddifftypeid = b.id AND a.clientid = ? AND a.mstorgnhirarchyid in (SELECT distinct mstorgnhirarchyid FROM mstgroupmember where userid=?) AND a.deleteflg = 0 AND a.activeflg = 1 AND b.deleteflg = 0 AND b.activeflg = 1 AND a.fromrecorddifftypeid = 2 AND a.fromrecorddiffid IN (SELECT distinct a.id FROM mstrecorddifferentiation a,mstgroupmember b WHERE b.userid=? AND b.mstorgnhirarchyid =a.mstorgnhirarchyid AND a.recorddifftypeid=2 AND a.seqno=1) AND b.parentid = 1 group by a.fromrecorddiffid) AND a.recorddiffid = c.id AND FIND_IN_SET(SUBSTRING_INDEX(SUBSTRING_INDEX(c.parentcategoryids, '->', 3), '->',- 1), (SELECT distinct GROUP_CONCAT(b.recorddiffid) FROM trnrecord a,maprecordtorecorddifferentiation b where (a.originalusergroupid IN (SELECT distinct groupid FROM mstgroupmember where userid=?) OR a.usergroupid IN (SELECT distinct groupid FROM mstgroupmember where userid=? AND deleteflg=0 AND activeflg=1)) AND a.id = b.recordid AND b.islatest=1 AND b.isworking=1)) GROUP BY a.recorddiffid order by cnt desc limit 10"

//var frequentresolver = "SELECT a.mstorgnhirarchyid,a.recorddiffid,COUNT(DISTINCT a.recordid) cnt,c.parentcategoryids,c.parentcategorynames,b.fromrecorddifftypeid,b.fromrecorddiffid,c.seqno FROM maprecordtorecorddifferentiation a,mstrecordtype b,mstrecorddifferentiation c WHERE b.torecorddiffid = a.recorddiffid AND b.clientid = ? AND b.mstorgnhirarchyid in (SELECT distinct mstorgnhirarchyid FROM mstgroupmember where userid=? AND deleteflg=0 AND activeflg=1) AND b.deleteflg = 0 AND b.activeflg = 1 AND b.fromrecorddifftypeid = 2 AND b.fromrecorddiffid IN (SELECT distinct a.id FROM mstrecorddifferentiation a,mstgroupmember b WHERE b.userid=? AND b.mstorgnhirarchyid =a.mstorgnhirarchyid AND a.recorddifftypeid=2 AND a.seqno=1 AND b.deleteflg=0 AND b.activeflg=1 AND a.deleteflg=0 AND a.activeflg=1) AND b.torecorddifftypeid IN (SELECT MAX(a.torecorddifftypeid) FROM mstrecordtype a, mstrecorddifferentiationtype b WHERE a.torecorddifftypeid = b.id AND a.clientid = ? AND a.mstorgnhirarchyid in (SELECT distinct mstorgnhirarchyid FROM mstgroupmember where userid=? AND deleteflg=0 AND activeflg=1) AND a.deleteflg = 0 AND a.activeflg = 1 AND b.deleteflg = 0 AND b.activeflg = 1 AND a.fromrecorddifftypeid = 2 AND a.fromrecorddiffid IN (SELECT distinct a.id FROM mstrecorddifferentiation a,mstgroupmember b WHERE b.userid=? AND b.mstorgnhirarchyid =a.mstorgnhirarchyid AND a.recorddifftypeid=2 AND a.seqno=1 AND b.deleteflg=0 AND b.activeflg=1 AND a.deleteflg=0 AND a.activeflg=1) AND b.parentid = 1 group by a.fromrecorddiffid) AND a.recorddiffid = c.id AND FIND_IN_SET(SUBSTRING_INDEX(SUBSTRING_INDEX(c.parentcategoryids, '->', 3), '->',- 1), (SELECT distinct GROUP_CONCAT(b.recorddiffid) FROM trnrecord a,maprecordtorecorddifferentiation b where (a.originalusergroupid IN (SELECT distinct groupid FROM mstgroupmember where userid=? AND deleteflg=0 AND activeflg=1) OR a.usergroupid IN (SELECT distinct groupid FROM mstgroupmember where userid=? AND deleteflg=0 AND activeflg=1)) AND a.id = b.recordid AND b.islatest=1 AND b.isworking=1)) GROUP BY a.recorddiffid order by cnt desc limit 10"


//var frequentresolver = "SELECT a.mstorgnhirarchyid,a.recorddiffid,COUNT(DISTINCT a.recordid) cnt,c.parentcategoryids,c.parentcategorynames,b.fromrecorddifftypeid,b.fromrecorddiffid,c.seqno FROM maprecordtorecorddifferentiation a,mstrecordtype b,mstrecorddifferentiation c WHERE b.torecorddiffid = a.recorddiffid AND b.clientid = ? AND b.mstorgnhirarchyid in (SELECT distinct mstorgnhirarchyid FROM mstgroupmember where userid=? AND deleteflg=0 AND activeflg=1) AND b.deleteflg = 0 AND b.activeflg = 1 AND b.fromrecorddifftypeid = 2 AND b.fromrecorddiffid IN (SELECT distinct a.id FROM mstrecorddifferentiation a,mstgroupmember b WHERE b.userid=? AND b.mstorgnhirarchyid =a.mstorgnhirarchyid AND a.recorddifftypeid=2 AND a.seqno=1 AND b.deleteflg=0 AND b.activeflg=1 AND a.deleteflg=0 AND a.activeflg=1) AND b.torecorddifftypeid IN (SELECT MAX(a.torecorddifftypeid) FROM mstrecordtype a, mstrecorddifferentiationtype b WHERE a.torecorddifftypeid = b.id AND a.clientid = ? AND a.mstorgnhirarchyid in (SELECT distinct mstorgnhirarchyid FROM mstgroupmember where userid=? AND deleteflg=0 AND activeflg=1) AND a.deleteflg = 0 AND a.activeflg = 1 AND b.deleteflg = 0 AND b.activeflg = 1 AND a.fromrecorddifftypeid = 2 AND a.fromrecorddiffid IN (SELECT distinct a.id FROM mstrecorddifferentiation a,mstgroupmember b WHERE b.userid=? AND b.mstorgnhirarchyid =a.mstorgnhirarchyid AND a.recorddifftypeid=2 AND a.seqno=1 AND b.deleteflg=0 AND b.activeflg=1 AND a.deleteflg=0 AND a.activeflg=1) AND b.parentid = 1 group by a.fromrecorddiffid) AND a.recorddiffid = c.id AND FIND_IN_SET(SUBSTRING_INDEX(SUBSTRING_INDEX(c.parentcategoryids, '->', 3), '->',- 1), (SELECT distinct GROUP_CONCAT(b.recorddiffid) FROM mstrecordactivitylogs a,maprecordtorecorddifferentiation b where a.createdgrpid IN (SELECT DISTINCT groupid FROM mstgroupmember WHERE userid = ? AND deleteflg=0 AND activeflg=1) AND a.recordid = b.recordid AND b.islatest=1 AND b.isworking=1)) GROUP BY a.recorddiffid order by cnt desc limit 10"

//var frequentresolver = "SELECT d.mstorgnhirarchyid,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid = d.recordid AND islatest =1 AND recorddifftypeid NOT IN (2,3,4,5,6,7,8,9,10,11,37) AND id > d.id order by id desc LIMIT 1) child,COUNT(DISTINCT d.recordid) cnt,(SELECT parentcategoryids FROM mstrecorddifferentiation where id=child) parentcategoryids,(SELECT parentcategorynames FROM mstrecorddifferentiation where id=child) parentcategorynames,c.recorddifftypeid,c.recorddiffid FROM mstgroupmember a,mstrecordactivitylogs b,maprecordtorecorddifferentiation c LEFT JOIN maprecordtorecorddifferentiation d ON c.recordid = d.recordid  AND d.isworking =1 AND d.islatest =1 WHERE a.userid = ? AND a.deleteflg = 0 AND a.activeflg = 1 AND a.groupid = b.createdgrpid AND b.recordid = c.recordid AND c.islatest=1 AND c.recorddifftypeid=2 AND c.recorddiffid IN (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=b.mstorgnhirarchyid AND recorddifftypeid=2 AND seqno=1) GROUP BY d.recorddiffid,child ORDER BY cnt DESC limit 10;"


var frequentresolver = "SELECT d.mstorgnhirarchyid,(SELECT recorddiffid FROM maprecordtorecorddifferentiation WHERE recordid = d.recordid AND islatest =1 AND recorddifftypeid NOT IN (2,3,4,5,6,7,8,9,10,11,37) AND id > d.id order by id desc LIMIT 1) child,COUNT(DISTINCT d.recordid) cnt,c.recorddifftypeid,c.recorddiffid FROM mstgroupmember a,mstrecordactivitylogs b,maprecordtorecorddifferentiation c LEFT JOIN maprecordtorecorddifferentiation d ON c.recordid = d.recordid  AND d.isworking =1 AND d.islatest =1 WHERE a.userid = ? AND a.deleteflg = 0 AND a.activeflg = 1 AND a.mstorgnhirarchyid = b.mstorgnhirarchyid AND a.groupid = b.createdgrpid AND b.recordid = c.recordid AND c.islatest=1 AND c.recorddifftypeid=2 AND c.recorddiffid IN (SELECT id FROM mstrecorddifferentiation WHERE clientid=? AND mstorgnhirarchyid=b.mstorgnhirarchyid AND recorddifftypeid=2 AND seqno=1) GROUP BY d.recorddiffid ORDER BY cnt DESC limit 10"

var parentrecords = "SELECT a.parentrecordid,b.code FROM mstparentchildmap a,trnrecord b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.childrecordid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.parentrecordid = b.id and a.isattached in ('Y','N')"
var termnmbyid = "SELECT termname FROM mstrecordterms where clientid=? AND mstorgnhirarchyid=? AND id=?"
var activitymst = "SELECT id as ID,activitydesc as Description,seqno as Seq  FROM mstrecordactivitymst WHERE clientid=? AND mstorgnhirarchyid=? AND activeflg=1 AND deleteflg=0"

//var newrecordlogs = "SELECT a.id, a.recordid, COALESCE(a.logValue,'') logvalue,a.createddate,concat(c.firstname,' ',c.lastname) name,b.activitydesc,d.supportgroupname,coalesce((SELECT termname FROM mstrecordterms where id =a.genericid),'') term,COALESCE(f.name,'') status FROM mstrecordactivitylogs a LEFT JOIN mststateterm e ON a.genericid = e.recordtermid and e.recorddifftypeid =3 and e.deleteflg=0 and e.activeflg=1 LEFT JOIN mstrecorddifferentiation f ON e.recorddiffid=f.id,mstrecordactivitymst b,mstclientuser c,mstclientsupportgroup d WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.activityseqno = b.seqno AND a.createdid = c.id AND a.createdgrpid=d.id and (a.genericid=0 or IF (a.genericid!=0 , a.genericid in (SELECT termid FROM mstsupportgrptermmap where grpid=? AND activeflg=1 AND deleteflg=0), 0 ) =1 ) order by a.id desc"

var newrecordlogs = "SELECT distinct a.id, a.recordid, COALESCE(a.logValue,'') logvalue,a.createddate,concat(c.firstname,' ',c.lastname) name,b.activitydesc,d.name,coalesce((SELECT termname FROM mstrecordterms where id =a.genericid),'') term,COALESCE(f.name,'') status FROM mstrecordactivitylogs a LEFT JOIN mststateterm e ON a.genericid = e.recordtermid and e.recorddifftypeid =3 and e.deleteflg=0 and e.activeflg=1 LEFT JOIN mstrecorddifferentiation f ON e.recorddiffid=f.id,mstrecordactivitymst b,mstclientuser c,mstsupportgrp d WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.activityseqno = b.seqno AND a.createdid = c.id AND a.createdgrpid=d.id and (a.genericid=0 or IF (a.genericid!=0 , a.genericid in (SELECT termid FROM mstsupportgrptermmap where grpid=? AND activeflg=1 AND deleteflg=0), 0 ) =1 ) order by a.id desc"

//var termsearchvalue = "SELECT a.id,a.recordid,a.logValue,a.createddate,b.activitydesc,concat(c.firstname,' ',c.lastname) name,d.supportgroupname,COALESCE(f.name,'') status,COALESCE(g.recordtrackdescription,'') description FROM mstrecordactivitymst b,mstclientuser c,mstclientsupportgroup d, mstrecordactivitylogs a LEFT JOIN mststateterm e ON a.genericid = e.recordtermid and e.recorddifftypeid =3 and e.deleteflg=0 and e.activeflg=1 LEFT JOIN mstrecorddifferentiation f ON e.recorddiffid=f.id LEFT JOIN trnreordtracking g ON g.recordid=a.recordid AND g.recordtermid=a.genericid where a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.genericid=? AND a.activityseqno=b.seqno AND a.createdid=c.id AND a.createdgrpid=d.id order by a.id desc"
var termsearchvalue = "SELECT distinct a.id,a.recordid,a.logValue,a.createddate,b.activitydesc,concat(c.firstname,' ',c.lastname) name,d.name,COALESCE(f.name,'') status FROM mstrecordactivitymst b,mstclientuser c,mstsupportgrp d, mstrecordactivitylogs a LEFT JOIN mststateterm e ON a.genericid = e.recordtermid and e.recorddifftypeid =3 and e.deleteflg=0 and e.activeflg=1 LEFT JOIN mstrecorddifferentiation f ON e.recorddiffid=f.id  where a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.genericid=? AND a.activityseqno=b.seqno AND a.createdid=c.id AND a.createdgrpid=d.id order by a.id desc"
var normalsearchvalue = "SELECT distinct a.id,a.recordid,a.logValue,a.createddate,b.activitydesc,concat(c.firstname,' ',c.lastname) name,d.name FROM mstrecordactivitymst b,mstclientuser c,mstsupportgrp d, mstrecordactivitylogs a where a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.activityseqno =? AND a.activityseqno=b.seqno AND a.createdid=c.id AND a.createdgrpid=d.id order by a.id desc"
var termseq = "SELECT id,seq FROM mstrecordterms WHERE clientid=? AND mstorgnhirarchyid=? AND deleteflg=0 AND activeflg=1 AND seq is not null"
var pendingstatustermvalue = "select distinct d.id,e.termname,d.recordtrackvalue,e.seq from mstrecorddifferentiation a,maprecordtorecorddifferentiation b,mststateterm c,trnreordtracking d,mstrecordterms e where a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddifftypeid=3 AND a.seqno=4 AND b.recordid=? AND b.recorddifftypeid=3 AND a.id= b.recorddiffid AND c.recorddifftypeid=3 AND b.recorddiffid=c.recorddiffid AND b.recordid=d.recordid AND d.recordtermid=c.recordtermid AND d.recordtermid=e.id order by d.id desc limit 4"
var attachfiles = "SELECT distinct b.id,b.recordid,b.clientid,b.mstorgnhirarchyid,b.recordtermid,COALESCE(b.recordtrackvalue,'') originalname,COALESCE(b.recordtrackdescription,'') uploadname,b.createddate,b.createdbyid,b.createdgrpid,COALESCE(concat(c.firstname,' ',c.lastname),'NA') name,d.supportgrouplevelid,e.userid,e.usergroupid,e.originaluserid,e.originalusergroupid FROM mstrecordterms a,trnreordtracking b,mstclientuser c,mstclientsupportgroup d,trnrecord e WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND b.activeflg=1 AND b.deleteflg=0 AND a.activeflg=1 AND a.termtypeid=3 AND b.recordid=? AND a.id= b.recordtermid AND b.createdbyid=c.id AND b.createdgrpid = d.grpid AND b.recordid=e.id"
var grplevel = "SELECT supportgrouplevelid FROM mstclientsupportgroup WHERE clientid=? AND mstorgnhirarchyid=? AND grpid=?"
var updatedoccount = "UPDATE mstdocumentdtls SET doccount =(doccount+1) WHERE clientid=? AND mstorgnhirarchyid=? AND id=?"

var termnamesbyseq = "select distinct a.id as ID,a.termname as Termname,b.recordtermvalue as Recordtermvalue,b.iscompulsory as Iscompulsory,c.termtypename as Termtypename,a.termtypeid as Termtypeid from mstrecordterms a,mststateterm b,msttermtype c,mstsupportgrptermmap d,mstclientsupportgroup e where b.clientid=? and b.mstorgnhirarchyid=? and b.recorddifftypeid=? and b.recorddiffid=? and a.id=b.recordtermid and a.termtypeid=c.id AND a.deleteflg=0 AND a.activeflg=1 AND e.id=? AND e.id=d.grpid AND d.termid=a.id AND a.seq=? AND d.deleteflg=0 AND d.activeflg=1 AND b.deleteflg=0 AND b.activeflg=1"
var rcount = "SELECT count(b.id) count FROM mstrecorddifferentiation a,maprecordtorecorddifferentiation b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recorddifftypeid=3 AND a.seqno=10 AND a.id = b.recorddiffid AND b.recordid=?"
var fcount = "SELECT count(a.id) count FROM trnreordtracking a,mstrecordterms b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.recordtermid =b.id AND b.seq=29"
var obcount = "SELECT count(a.id) count FROM trnreordtracking a,mstrecordterms b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.recordtermid =b.id AND b.seq=30"
var latestrecordcomment = "SELECT a.clientid,a.mstorgnhirarchyid,a.recordid,a.recordtrackvalue,FROM_UNIXTIME(a.createddate) createddate FROM trnreordtracking a,mstrecordterms b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.recordtermid =b.id AND b.seq=11 order by a.createddate desc limit 1"
var updateattach = "UPDATE trnreordtracking SET deleteflg=1 WHERE id=?"
var deletelogvalue = "UPDATE mstrecordactivitylogs SET deleteflg=1 WHERE recordid=? AND createdid=? AND createdgrpid=? AND genericid=? AND logValue like ?"
var termsseqvalue = "SELECT a.recordtrackvalue,concat(c.firstname,' ',c.lastname) name,a.createddate FROM trnreordtracking a,mstrecordterms b,mstclientuser c WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.recordid=? AND a.recordtermid = b.id AND b.seq=? AND a.createdbyid=c.id"
var getchildidswithcode = "select a.childrecordid,b.code from mstparentchildmap a,trnrecord b where a.parentrecordid=? and a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and a.deleteflg=0 and a.activeflg=1 and a.isattached='Y' AND a.childrecordid=b.id"
var termsearchvaluebyseq = "SELECT distinct a.id,a.recordid,a.logValue,a.createddate,b.activitydesc,concat(c.firstname,' ',c.lastname) name,d.name,COALESCE(f.name,'') status FROM mstrecordactivitymst b,mstclientuser c,mstsupportgrp d, mstrecordactivitylogs a LEFT JOIN mststateterm e ON a.genericid = e.recordtermid and e.recorddifftypeid =3 and e.deleteflg=0 and e.activeflg=1 LEFT JOIN mstrecorddifferentiation f ON e.recorddiffid=f.id  where a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.genericid in (SELECT id FROM mstrecordterms WHERE clientid=? AND mstorgnhirarchyid=? AND seq=?) AND a.activityseqno=b.seqno AND a.createdid=c.id AND a.createdgrpid=d.id order by a.id desc"
var aging_query = "SELECT DATEDIFF(NOW(),FROM_UNIXTIME(trnrecord.createdatetime)) age FROM trnrecord WHERE clientid = ? AND mstorgnhirarchyid = ? AND id = ?"
var recordcreatedt = "SELECT FROM_UNIXTIME(createdatetime) recordcreatedate FROM trnrecord WHERE clientid=? AND mstorgnhirarchyid=? AND id=?"


func (mdao DbConn) Getclientorgvalues(Recordid int64) (int64, int64, int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	var clientID int64
	var orgid int64
	var recordstageid int64
	rows, err := mdao.DB.Query("SELECT clientid,mstorgnhirarchyid,recordstageid FROM trnrecord where  id=? and deleteflg=0;", Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAllcommontermvalues Get Statement Prepare Error", err)
		return clientID, orgid, recordstageid, err
	}
	for rows.Next() {
		rows.Scan(&clientID, &orgid, &recordstageid)
	}
	return clientID, orgid, recordstageid, nil
}
func (mdao DbConn) InsertRecordTermvalues(rec *entities.RecordcommonEntity) (int64, error) {
	logger.Log.Println("trnreordtracking query -->", trnreordtracking)
	logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, rec.RecordID, rec.RecordstageID, rec.TermID, rec.Termvalue, rec.ForuserID, rec.Userid, rec.Usergroupid)

	stmt, err := mdao.DB.Prepare(trnreordtracking)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, rec.RecordID, rec.RecordstageID, rec.TermID, rec.Termvalue, rec.Userid, rec.Usergroupid, rec.Termdescription)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Last ID fetch Error")
	}

	return lastInsertedID, nil
}

func (mdao DbConn) InsertRecordTermvaluesForchilds(rec *entities.RecordcommonEntity, recordid int64) (int64, error) {
	logger.Log.Println("trnreordtracking query -->", trnreordtracking)
	logger.Log.Println("parameters -->", rec.ClientID, rec.Mstorgnhirarchyid, rec.RecordID, rec.RecordstageID, rec.TermID, rec.Termvalue, rec.ForuserID, rec.Userid)

	stmt, err := mdao.DB.Prepare(trnreordtracking)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(rec.ClientID, rec.Mstorgnhirarchyid, recordid, rec.RecordstageID, rec.TermID, rec.Termvalue, rec.Userid, rec.Usergroupid, rec.Termdescription)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Last ID fetch Error")
	}

	return lastInsertedID, nil
}

func InsertMultipleRecordTermvalues(tx *sql.Tx, rec *entities.RecordTermnamesEntity, ClientID int64, Mstorgnhirarchyid int64, RecordID int64, RecordstageID int64, ForuserID int64, Userid int64, Grpid int64) (int64, error) {
	logger.Log.Println("trnreordtracking query -->", trnreordtracking)
	logger.Log.Println("parameters -->", ClientID, Mstorgnhirarchyid, RecordID, RecordstageID, rec.ID, rec.Insertedvalue, Userid)

	stmt, err := tx.Prepare(trnreordtracking)

	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, RecordstageID, rec.ID, rec.Insertedvalue, Userid, Grpid, rec.Termdescription)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	lastInsertedID, err := res.LastInsertId()
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Last ID fetch Error")
	}

	return lastInsertedID, nil
}

func (mdao DbConn) GetAllcommontermvalues(page *entities.RecordcommonEntity) ([]entities.RecordcommonresponseEntity, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	values := []entities.RecordcommonresponseEntity{}
	rows, err := mdao.DB.Query(commontabvalues, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Userid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAllcommontermvalues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordcommonresponseEntity{}
		rows.Scan(&value.ID, &value.TermID, &value.Termvalue, &value.ForuserID, &value.Termname, &value.Recorddiffname, &value.Username, &value.Createddate)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetTermvalueagainsttermid(page *entities.RecordcommonEntity) ([]entities.RecordcommonresponseEntity, error) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	values := []entities.RecordcommonresponseEntity{}
	rows, err := mdao.DB.Query(termvalueagainsttermid, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.TermID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermvalueagainsttermid Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordcommonresponseEntity{}
		rows.Scan(&value.ID, &value.TermID, &value.Termvalue, &value.ForuserID, &value.Termname, &value.Recorddiffname, &value.Username, &value.Createddate)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetTermnamesbystate(page *entities.RecordcommonstateEntity) ([]entities.RecordTermnamesEntity, error) {
	logger.Log.Println("In side GetTermnamesbystate")
	values := []entities.RecordTermnamesEntity{}
	rows, err := mdao.DB.Query(termnamesbystate, page.ClientID, page.Mstorgnhirarchyid, page.Recordtickettypedifftypeid, page.Recordtickettypediffid, page.Recordstatusdifftypeid, page.Recordstatusdiffid, page.Userid, page.ClientID, page.Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermnamesbystate Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) Getchildrecordids(RecordID int64, ClientID int64, Mstorgnhirarchyid int64, Recorddifftypeid int64, Recorddiffid int64) ([]int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	logger.Log.Println("Query is --++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++---->", getchildids)
	logger.Log.Println("Parameter is ------>", RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid)

	var childids []int64
	rows, err := mdao.DB.Query(getchildids, RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid) //ClientID, Mstorgnhirarchyid,
	//rows, err := mdao.DB.Query(getchildids, RecordID, ClientID, Mstorgnhirarchyid, 2, 4)
	logger.Log.Println("Rows error is  ----+++++++++++++++++++++++++++++++++++++++++-->", err)

	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return childids, err
	}
	for rows.Next() {
		var ID int64
		err = rows.Scan(&ID)
		childids = append(childids, ID)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return childids, nil
}

func (mdao DbConn) GetchildrecordidForPendingVendorActions(RecordID int64, ClientID int64, Mstorgnhirarchyid int64, Recorddifftypeid int64, Recorddiffid int64) ([]int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	var sql = "select a.childrecordid from mstparentchildmap a,maprecordtorecorddifferentiation b,mstrecorddifferentiation c where a.parentrecordid=? and a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and a.deleteflg=0 and a.activeflg=1 and a.isattached='Y' AND a.childrecordid=b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND c.recorddifftypeid=3 AND c.seqno=5 AND c.deleteflg=0 and c.activeflg=1 AND b.recorddiffid=c.id"
	logger.Log.Println("GetchildrecordidForPendingVendorActions Query is --++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++---->", sql)
	logger.Log.Println("Parameter is ------>", RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid)
	var childids []int64
	rows, err := mdao.DB.Query(sql, RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid) //ClientID, Mstorgnhirarchyid,
	//rows, err := mdao.DB.Query(getchildids, RecordID, ClientID, Mstorgnhirarchyid, 2, 4)
	logger.Log.Println("Rows error is  ----+++++++++++++++++++++++++++++++++++++++++-->", err)

	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return childids, err
	}
	for rows.Next() {
		var ID int64
		err = rows.Scan(&ID)
		childids = append(childids, ID)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return childids, nil
}

func (mdao DbConn) Getchildrecordidswithcode(RecordID int64, ClientID int64, Mstorgnhirarchyid int64, Recorddifftypeid int64, Recorddiffid int64, Recordcode string) ([]entities.ChildRecordEntity, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	logger.Log.Println("Query is ------>", getchildids)
	logger.Log.Println("Parameter is ------>", RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid)

	values := []entities.ChildRecordEntity{}
	rows, err := mdao.DB.Query(getchildidswithcode, RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid)
	//rows, err := mdao.DB.Query(getchildids, RecordID, ClientID, Mstorgnhirarchyid, 2, 4)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.ChildRecordEntity{}
		err = rows.Scan(&value.ID, &value.Code)
		values = append(values, value)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	if len(values) > 0 {
		t := entities.ChildRecordEntity{}
		t.ID = RecordID
		t.Code = Recordcode
		values = append(values, t)
	}

	return values, nil
}

func (mdao DbConn) Checktermreleation(TermID int64, ClientID int64, Mstorgnhirarchyid int64, Recorddifftypeid int64, Recorddiffid int64) (int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	var ID int64
	rows, err := mdao.DB.Query(termreleation, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid, TermID)
	//rows, err := mdao.DB.Query(termreleation, ClientID, Mstorgnhirarchyid, 2, 4, TermID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return ID, err
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return ID, nil
}

func (mdao DbConn) GetPrioritycount(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	var count int64
	rows, err := mdao.DB.Query(pcount, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return count, nil
}

func (mdao DbConn) GetPendingvendorcount(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetAllcommontermvalues")
	var count int64
	rows, err := mdao.DB.Query(pvacount, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return count, nil
}

func (mdao DbConn) GetReopencount(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetReopencount")
	var count int64
	rows, err := mdao.DB.Query(rcount, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetReopencount Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("GetReopencount rows.next() Error", err)
	}
	return count, nil
}

func (mdao DbConn) GetFollowupcount(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetFollowupcount")
	var count int64
	rows, err := mdao.DB.Query(fcount, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetFollowupcount Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("GetFollowupcount rows.next() Error", err)
	}
	return count, nil
}

func (mdao DbConn) GetOutboundcount(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetOutboundcount")
	var count int64
	rows, err := mdao.DB.Query(obcount, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOutboundcount Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("GetOutboundcount rows.next() Error", err)
	}
	return count, nil
}

func (mdao DbConn) GetRecentrecords(page *entities.RecordcommonEntity, Timediff int64) ([]entities.RecentrecordEntity, error) {
	logger.Log.Println("In side GetRecentrecords")
	values := []entities.RecentrecordEntity{}
	rows, err := mdao.DB.Query(recent, page.Userid, page.Userid, page.ClientID, page.Mstorgnhirarchyid, page.Userid, page.Userid, page.ClientID, page.Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecentrecords Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecentrecordEntity{}
		rows.Scan(&value.ID, &value.Code, &value.Title, &value.Createdate, &value.Status, &value.Seq, &value.OrgnID)
		value.Showcreatedate = Convertdate(value.Createdate, Timediff)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetRecentrecordsForResolver(page *entities.RecordcommonEntity, Timediff int64, orgntype int64) ([]entities.RecentrecordEntity, error) {
	logger.Log.Println("In side GetRecentrecordsForResolver")
	values := []entities.RecentrecordEntity{}
	var params []interface{}
	var resolverrecent = ""
	if orgntype == 2 {
		resolverrecent = "SELECT a.id,a.code,a.recordtitle,a.createdatetime,c.name,c.seqno,a.mstorgnhirarchyid FROM trnrecord a,maprecordtorecorddifferentiation b,mstrecorddifferentiation c,mstrecordactivitylogs d where d.createdid=? AND d.recordid = a.id AND a.clientid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.id=b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND b.recorddiffid=c.id and d.id in (select max(x.id) from mstrecordactivitylogs x, trnrecord y where x.createdid=? AND x.recordid = y.id AND y.clientid=? AND y.deleteflg=0 AND y.activeflg=1 group by y.id) group by d.recordid order by d.createddate desc limit 10"
		params = append(params, page.Userid)
		params = append(params, page.ClientID)
		params = append(params, page.Userid)
		params = append(params, page.ClientID)
	} else {
		resolverrecent = "SELECT a.id,a.code,a.recordtitle,a.createdatetime,c.name,c.seqno,a.mstorgnhirarchyid FROM trnrecord a,maprecordtorecorddifferentiation b,mstrecorddifferentiation c,mstrecordactivitylogs d where d.createdid=? AND d.recordid = a.id AND a.clientid=? AND a.mstorgnhirarchyid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.id=b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND b.recorddiffid=c.id and d.id in (select max(x.id) from mstrecordactivitylogs x, trnrecord y where x.createdid=? AND x.recordid = y.id AND y.clientid=? AND y.mstorgnhirarchyid=? AND y.deleteflg=0 AND y.activeflg=1 group by y.id) group by d.recordid order by d.createddate desc limit 10"
		params = append(params, page.Userid)
		params = append(params, page.ClientID)
		params = append(params, page.Mstorgnhirarchyid)
		params = append(params, page.Userid)
		params = append(params, page.ClientID)
		params = append(params, page.Mstorgnhirarchyid)
	}

	rows, err := mdao.DB.Query(resolverrecent, params...)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecentrecordsForResolver Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecentrecordEntity{}
		rows.Scan(&value.ID, &value.Code, &value.Title, &value.Createdate, &value.Status, &value.Seq, &value.OrgnID)
		value.Showcreatedate = Convertdate(value.Createdate, Timediff)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetRecordlogs(page *entities.RecordcommonEntity) ([]entities.RecordlogsEntity, error) {
	logger.Log.Println("In side GetRecordlogs")
	values := []entities.RecordlogsEntity{}
	rows, err := mdao.DB.Query(recordlogs, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordlogs Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordlogsEntity{}
		rows.Scan(&value.ID, &value.Logvalue, &value.Createdate, &value.Name, &value.Activitydesc, &value.Supportgroupname)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) Getfrequentissues(page *entities.RecordcommonEntity) ([]entities.FrequentRecordEntity, error) {
	logger.Log.Println("In side Getfrequentissues")
	values := []entities.FrequentRecordEntity{}
	rows, err := mdao.DB.Query(frequent, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.FrequentRecordEntity{}
		rows.Scan(&value.Mstorgnhirarchyid, &value.LastlevelID, &value.Count, &value.ParentcatID, &value.Parentcatname, &value.Recorddifftypeid, &value.Recorddiffid, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

/*func (mdao DbConn) Getfrequentissuesresolver(page *entities.RecordcommonEntity) ([]entities.FrequentRecordEntity, error) {
	logger.Log.Println("In side Getfrequentissues")
	values := []entities.FrequentRecordEntity{}
	//rows, err := mdao.DB.Query(frequentresolver, page.ClientID, page.Userid, page.Userid, page.ClientID, page.Userid, page.Userid, page.Userid) //page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.Usergroupid
	rows, err := mdao.DB.Query(frequentresolver, page.Userid, page.ClientID) //, page.Userid, page.Userid, page.ClientID, page.Userid, page.Userid, page.Userid, page.Useriddefer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.FrequentRecordEntity{}
		rows.Scan(&value.Mstorgnhirarchyid, &value.LastlevelID, &value.Count, &value.ParentcatID, &value.Parentcatname, &value.Recorddifftypeid, &value.Recorddiffid, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}*/

func (mdao DbConn) GetParentrecord(page *entities.RecordcommonEntity) ([]entities.ParentticketEntity, error) {
	logger.Log.Println("In side GetParentrecord")
	values := []entities.ParentticketEntity{}
	rows, err := mdao.DB.Query(parentrecords, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetParentrecord Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.ParentticketEntity{}
		rows.Scan(&value.ID, &value.Recordnumber)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) Gettermnamebyid(Termid int64, ClientID int64, Mstorgnhirarchyid int64) (string, error) {
	logger.Log.Println("In side Gettermnamebyid")
	var tername string
	rows, err := mdao.DB.Query(termnmbyid, ClientID, Mstorgnhirarchyid, Termid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Gettermnamebyid Get Statement Prepare Error", err)
		return tername, err
	}
	for rows.Next() {
		err = rows.Scan(&tername)
		logger.Log.Println("Gettermnamebyid rows.next() Error", err)
	}
	return tername, nil
}

func Gettermnamebyid(tx *sql.Tx, Termid int64, ClientID int64, Mstorgnhirarchyid int64) (string, error) {
	logger.Log.Println("In side Gettermnamebyid")
	var tername string
	rows, err := tx.Query(termnmbyid, ClientID, Mstorgnhirarchyid, Termid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Gettermnamebyid Get Statement Prepare Error", err)
		return tername, err
	}
	for rows.Next() {
		err = rows.Scan(&tername)
		logger.Log.Println("Gettermnamebyid rows.next() Error", err)
	}
	return tername, nil
}

func (mdao DbConn) GetTermnames(page *entities.RecordcommonEntity) ([]entities.RecordTermnamesEntity, error) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	values := []entities.RecordTermnamesEntity{}
	rows, err := mdao.DB.Query(termnames, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.Userid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermvalueagainsttermid Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid)
		value.Seq = 100
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetActivitymstnames(page *entities.RecordcommonEntity) ([]entities.Recordactivitymst, error) {
	logger.Log.Println("In side Getfrequentissues")
	values := []entities.Recordactivitymst{}
	rows, err := mdao.DB.Query(activitymst, page.ClientID, page.Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.Recordactivitymst{}
		rows.Scan(&value.ID, &value.Description, &value.Seq)
		if value.Seq == 100 {
			val, _ := mdao.GetTermnames(page)
			value.Details = val
		}
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetNewActivitylogs(page *entities.RecordcommonEntity) ([]entities.NewActivitylogsEntity, error) {
	logger.Log.Println("In side Getfrequentissues")
	values := []entities.NewActivitylogsEntity{}
	rows, err := mdao.DB.Query(newrecordlogs, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Usergroupid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.NewActivitylogsEntity{}
		rows.Scan(&value.ID, &value.RecordID, &value.Logvalue, &value.Createddate, &value.Name, &value.Activitydesc, &value.Supportgroupname, &value.Termname, &value.Status)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) Searchtermlogs(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, TermID int64, values []entities.NewActivitylogsEntity, Timediff int64) ([]entities.NewActivitylogsEntity, error) {
	logger.Log.Println("In side Getfrequentissues------------->", termsearchvalue)
	//values := []entities.NewActivitylogsEntity{}
	rows, err := mdao.DB.Query(termsearchvalue, ClientID, Mstorgnhirarchyid, RecordID, TermID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.NewActivitylogsEntity{}
		err = rows.Scan(&value.ID, &value.RecordID, &value.Logvalue, &value.Createddate, &value.Activitydesc, &value.Name, &value.Supportgroupname, &value.Status)
		logger.Log.Println("Error is ------------->", err)
		value.Showcreatedate = Convertdate(value.Createddate, Timediff)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) Searchtermlogsbysequance(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, Termsequance int64, values []entities.NewActivitylogsEntity, Timediff int64, Code string) ([]entities.NewActivitylogsEntity, error) {
	logger.Log.Println("In side Getfrequentissues------------->", termsearchvaluebyseq)
	//values := []entities.NewActivitylogsEntity{}
	rows, err := mdao.DB.Query(termsearchvaluebyseq, ClientID, Mstorgnhirarchyid, RecordID, ClientID, Mstorgnhirarchyid, Termsequance)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	var str string
	for rows.Next() {
		value := entities.NewActivitylogsEntity{}
		err = rows.Scan(&value.ID, &value.RecordID, &value.Logvalue, &value.Createddate, &value.Activitydesc, &value.Name, &value.Supportgroupname, &value.Status)
		logger.Log.Println("Error is ------------->", err)
		if str != value.Logvalue {
			value.Code = Code
			value.Showcreatedate = Convertdate(value.Createddate, Timediff)
			values = append(values, value)
			str = value.Logvalue
		}

	}
	return values, nil
}

func (mdao DbConn) Searchnormallogs(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, Seqno int64, values []entities.NewActivitylogsEntity, Timediff int64) ([]entities.NewActivitylogsEntity, error) {
	//logger.Log.Println("values is ------------->", values)
	logger.Log.Println("In side Getfrequentissues---------------->", normalsearchvalue)
	//ids := []int{1, 4}
	logger.Log.Println("In side Getfrequentissues---------------->", ClientID, Mstorgnhirarchyid, RecordID, Seqno)
	//values := []entities.NewActivitylogsEntity{}
	//pq.Array(ids)
	rows, err := mdao.DB.Query(normalsearchvalue, ClientID, Mstorgnhirarchyid, RecordID, Seqno)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.NewActivitylogsEntity{}
		err = rows.Scan(&value.ID, &value.RecordID, &value.Logvalue, &value.Createddate, &value.Activitydesc, &value.Name, &value.Supportgroupname)
		logger.Log.Println("Error is ------------->", err)
		value.Showcreatedate = Convertdate(value.Createddate, Timediff)
		values = append(values, value)
	}
	//logger.Log.Println("values is ------------->", values)
	return values, nil
}

func (mdao DbConn) Termsequance(ClientID int64, Mstorgnhirarchyid int64) (map[int64]int64, error) {
	logger.Log.Println("In side Getfrequentissues---------------->", termseq)
	var t = make(map[int64]int64)
	var ID int64
	var Seq int64
	rows, err := mdao.DB.Query(termseq, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return t, err
	}
	for rows.Next() {
		err = rows.Scan(&ID, &Seq)
		t[Seq] = ID
	}
	logger.Log.Println("Hashmap value is---------------->", t)
	return t, nil
}

func (mdao DbConn) TermsIDs(ClientID int64, Mstorgnhirarchyid int64) (map[int64]int64, error) {
	logger.Log.Println("In side Getfrequentissues---------------->", termseq)
	var t = make(map[int64]int64)
	var ID int64
	var Seq int64
	rows, err := mdao.DB.Query(termseq, ClientID, Mstorgnhirarchyid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return t, err
	}
	for rows.Next() {
		err = rows.Scan(&ID, &Seq)
		t[ID] = Seq
	}
	logger.Log.Println("Hashmap value is---------------->", t)
	return t, nil
}

func (mdao DbConn) GetPendingstatustermvalue(page *entities.RecordcommonEntity) ([]entities.Pendingstatustermvalue, error) {
	logger.Log.Println("In side Getfrequentissues")
	values := []entities.Pendingstatustermvalue{}
	rows, err := mdao.DB.Query(pendingstatustermvalue, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.Pendingstatustermvalue{}
		rows.Scan(&value.ID, &value.Termname, &value.Termvalue, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetAttachmentfiles(page *entities.RecordcommonEntity) ([]entities.RecordAttachmentfiles, error) {
	logger.Log.Println("In side GetAttachmentfiles")
	values := []entities.RecordAttachmentfiles{}
	rows, err := mdao.DB.Query(attachfiles, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAttachmentfiles Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordAttachmentfiles{}
		rows.Scan(&value.ID, &value.RecordID, &value.ClientID, &value.Mstorgnhirarchyid, &value.RecordtermID, &value.Originalname, &value.Uploadname, &value.Createdate, &value.Createdbyid, &value.Createdgrpid, &value.Name, &value.Supportgrouplevelid, &value.RecorduserID, &value.RecordusergrpID, &value.RecordoriginaluserID, &value.RecordoriginalusergrpID)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetgrplevelID(ClientID int64, Mstorgnhirarchyid int64, GrpID int64) (int64, error) {
	logger.Log.Println("In side GetgrplevelID")
	var grplevelID int64
	rows, err := mdao.DB.Query(grplevel, ClientID, Mstorgnhirarchyid, GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetgrplevelID Get Statement Prepare Error", err)
		return grplevelID, err
	}
	for rows.Next() {
		err = rows.Scan(&grplevelID)
		logger.Log.Println("GetgrplevelID rows.next() Error", err)
	}
	return grplevelID, nil
}

func (mdao DbConn) Updatedocumentcount(ClientID int64, Mstorgnhirarchyid int64, ID int64) error {
	logger.Log.Println("Updatedocumentcount query -->", updatedoccount)
	logger.Log.Println("Updatedocumentcount parameters -->", ClientID, Mstorgnhirarchyid, ID)
	stmt, err := mdao.DB.Prepare(updatedoccount)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, ID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetTermnamesbysequance(page *entities.RecordcommonEntity, seq int64) (entities.RecordTermnamesEntity, error) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	values := entities.RecordTermnamesEntity{}
	rows, err := mdao.DB.Query(termnamesbyseq, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.Usergroupid, seq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermvalueagainsttermid Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		rows.Scan(&values.ID, &values.Termname, &values.Recordtermvalue, &values.Iscompulsory, &values.Termtypename, &values.Termtypeid)
	}
	logger.Log.Println(values)
	return values, nil
}

//latestrecordcomment

func (mdao DbConn) GetLastRecordcomment(page *entities.RecordcommonEntity) ([]entities.Customervisiblecomment, error) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	values := []entities.Customervisiblecomment{}
	rows, err := mdao.DB.Query(latestrecordcomment, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermvalueagainsttermid Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.Customervisiblecomment{}
		rows.Scan(&value.ClientID, &value.Mstorgnhirarchyid, &value.RecordID, &value.Recordtrackvalue, &value.Createddate)
		values = append(values, value)
	}
	logger.Log.Println(values)
	return values, nil
}

func Updateattachfiles(tx *sql.Tx, ID int64) error {
	stmt, err := tx.Prepare(updateattach)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func Deletefromactivitylogs(tx *sql.Tx, RecordID int64, Originalname string, Createdbyid int64, Createdgrpid int64, RecordtermID int64) error {
	stmt, err := tx.Prepare(deletelogvalue)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(RecordID, Createdbyid, Createdgrpid, RecordtermID, "%"+Originalname+"%")
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetTermvaluebysequance(page *entities.RecordcommonEntity) ([]entities.Recordtermseqvalue, error) {
	logger.Log.Println("In side GetTermvalueagainsttermid")
	values := []entities.Recordtermseqvalue{}
	rows, err := mdao.DB.Query(termsseqvalue, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Termseq)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetTermvalueagainsttermid Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.Recordtermseqvalue{}
		rows.Scan(&value.Recordtermvalue, &value.Name, &value.Createddate)
		values = append(values, value)
	}
	logger.Log.Println(values)
	return values, nil
}

func (mdao DbConn) GetAging(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetAging")
	var aging_count int64
	rows, err := mdao.DB.Query(aging_query, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAging Get Statement Prepare Error", err)
		return aging_count, err
	}
	for rows.Next() {
		err = rows.Scan(&aging_count)
		logger.Log.Println("GetAging rows.next() Error", err)
	}
	return aging_count, nil
}

func (mdao DbConn) GetRecordcreatedate(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (string, error) {
	logger.Log.Println("In side GetRecordcreatedate")
	var createdate string
	rows, err := mdao.DB.Query(recordcreatedt, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetRecordcreatedate Get Statement Prepare Error", err)
		return createdate, err
	}
	for rows.Next() {
		err = rows.Scan(&createdate)
		logger.Log.Println("GetRecordcreatedate rows.next() Error", err)
	}
	return createdate, nil
}

func (mdao DbConn) GetOrgnType(ClientID int64, OrgnID int64) (int64, error) {
	logger.Log.Println("In side GetOrgnType")
	var OrgnTypeID int64
	var sql = "SELECT mstorgnhierarchytypeid FROM mstorgnhierarchy WHERE clientid=? AND id=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return OrgnTypeID, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(ClientID, OrgnID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOrgnType Get Statement Prepare Error", err)
		return OrgnTypeID, err
	}
	for rows.Next() {
		err := rows.Scan(&OrgnTypeID)
		logger.Log.Println("Error is >>>>>>>", err)
	}
	return OrgnTypeID, nil
}

func (mdao DbConn) GetOrgnTypeForGrid(req map[string]interface{}) (int64, error) {
	logger.Log.Println("In side GetOrgnType")
	var OrgnTypeID int64
	var sql = "SELECT mstorgnhierarchytypeid FROM mstorgnhierarchy WHERE clientid=? AND id=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return OrgnTypeID, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(req["clientid"], req["mstorgnhirarchyid"])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOrgnType Get Statement Prepare Error", err)
		return OrgnTypeID, err
	}
	for rows.Next() {
		err := rows.Scan(&OrgnTypeID)
		logger.Log.Println("Error is >>>>>>>", err)
	}
	return OrgnTypeID, nil
}

// func (mdao DbConn) GetOrgnIDForGrid(req map[string]interface{}) (map[string]interface{}, error) {
// 	logger.Log.Println("In side GetOrgnType")
// 	values := map[string]interface{}{}
// 	t := []interface{}{}
// 	var sql = "SELECT distinct mstorgnhirarchyid FROM mstgroupmember WHERE clientid=? AND groupid=? AND userid=? AND deleteflg=0 AND activeflg=1"
// 	stmt, err := mdao.DB.Prepare(sql)
// 	if err != nil {
// 		logger.Log.Println(err)
// 		return values, err
// 	}
// 	defer stmt.Close()
// 	rows, err := stmt.Query(req["clientid"], req["supportgrpid"], req["userid"])
// 	defer rows.Close()
// 	if err != nil {
// 		logger.Log.Println("GetOrgnType Get Statement Prepare Error", err)
// 		return values, err
// 	}
// 	for rows.Next() {
// 		var OrgnTypeID int64
// 		err := rows.Scan(&OrgnTypeID)
// 		logger.Log.Println("Error is >>>>>>>", err)
// 		t = append(t, OrgnTypeID)
// 	}
// 	values["orgnids"] = t
// 	logger.Log.Println("values ------------------>", values)
// 	return values, nil
// }

func (mdao DbConn) GetOrgnIDForGrid(req map[string]interface{}) ([]int64, error) {
	logger.Log.Println("In side GetOrgnType")
	values := []int64{}
	var sql = "SELECT distinct mstorgnhirarchyid FROM mstgroupmember WHERE clientid=? AND groupid=? AND userid=? AND deleteflg=0 AND activeflg=1"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return values, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(req["clientid"], req["supportgrpid"], req["userid"])
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOrgnType Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		var OrgnTypeID int64
		err := rows.Scan(&OrgnTypeID)
		logger.Log.Println("Error is >>>>>>>", err)
		values = append(values, OrgnTypeID)
	}
	logger.Log.Println("values ------------------>", values)
	return values, nil
}

func (mdao DbConn) GetRecordID(ClientID int64, Mstorgnhirarchyid int64, Recordno string) (int64, error) {
	logger.Log.Println("In side GetAging")
	var ID int64
	var sql = "SELECT id FROM trnrecord WHERE clientid=? AND mstorgnhirarchyid=? AND code=?"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, Recordno)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetAging Get Statement Prepare Error", err)
		return ID, err
	}
	for rows.Next() {
		err = rows.Scan(&ID)
		logger.Log.Println("GetAging rows.next() Error", err)
	}
	return ID, nil
}

func (mdao DbConn) GetScheduleTabTermnames(page *entities.RecordcommonEntity) ([]entities.RecordScheduleTabTermnamesEntity, error) {
	logger.Log.Println("In side GetScheduleTabTermnames")
	var sql = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,a.iscompulsory as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,e.readpermission,e.writepermission,f.id,c.seq from mststateterm a,mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstrecordfield f where a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and a.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND e.grpid=? AND e.termid=c.id AND e.deleteflg=0 and e.activeflg=1 AND f.mstrecordfieldtype='SCHEDULE TAB' AND f.recordtermid=c.id AND f.deleteflg=0 and f.activeflg=1 ORDER BY f.displayseq"
	values := []entities.RecordScheduleTabTermnamesEntity{}
	rows, err := mdao.DB.Query(sql, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetScheduleTabTermnames Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordScheduleTabTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Readpermission, &value.Writepermission, &value.FieldID, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetScheduleTabTermnameswithStatus(page *entities.RecordcommonEntity) ([]entities.RecordScheduleTabTermnamesEntity, error) {
	logger.Log.Println("In side GetScheduleTabTermnames")
	var sql = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,b.iscompulsory as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,e.readpermission,e.writepermission,f.id,c.seq from mststateterm a,mststateterm b,mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstrecordfield f where a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and b.recorddifftypeid=? and b.recorddiffid=? and a.recordtermid=b.recordtermid and a.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND e.grpid=? AND e.termid=c.id AND e.deleteflg=0 and e.activeflg=1 AND f.mstrecordfieldtype='SCHEDULE TAB' AND f.recordtermid=c.id AND f.deleteflg=0 and f.activeflg=1 AND b.deleteflg=0 and b.activeflg=1 ORDER BY f.displayseq"
	values := []entities.RecordScheduleTabTermnamesEntity{}
	rows, err := mdao.DB.Query(sql, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.Recordstatustypeid, page.Recordstatusid, page.GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetScheduleTabTermnames Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordScheduleTabTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Readpermission, &value.Writepermission, &value.FieldID, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetPlanTabTermnames(page *entities.RecordcommonEntity) ([]entities.RecordPlanTabTermnamesEntity, error) {
	logger.Log.Println("In side GetScheduleTabTermnames")
	var sql = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,a.iscompulsory as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,e.readpermission,e.writepermission,f.id,c.seq from mststateterm a,mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstrecordfield f where a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and a.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND e.grpid=? AND e.termid=c.id AND e.deleteflg=0 and e.activeflg=1 AND f.mstrecordfieldtype='PLAN OF ACTION' AND f.recordtermid=c.id AND f.deleteflg=0 and f.activeflg=1 ORDER BY f.displayseq"
	values := []entities.RecordPlanTabTermnamesEntity{}
	rows, err := mdao.DB.Query(sql, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetScheduleTabTermnames Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordPlanTabTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Readpermission, &value.Writepermission, &value.FieldID, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetPlanTabTermnameswithStatus(page *entities.RecordcommonEntity) ([]entities.RecordPlanTabTermnamesEntity, error) {
	logger.Log.Println("In side GetScheduleTabTermnames")
	var sql = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,b.iscompulsory as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,e.readpermission,e.writepermission,f.id,c.seq from mststateterm a,mststateterm b,mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstrecordfield f where a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and b.recorddifftypeid=? and b.recorddiffid=? and a.recordtermid=b.recordtermid and a.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND e.grpid=? AND e.termid=c.id AND e.deleteflg=0 and e.activeflg=1 AND f.mstrecordfieldtype='PLAN OF ACTION' AND f.recordtermid=c.id AND f.deleteflg=0 and f.activeflg=1 AND b.deleteflg=0 and b.activeflg=1 ORDER BY f.displayseq"
	values := []entities.RecordPlanTabTermnamesEntity{}
	rows, err := mdao.DB.Query(sql, page.ClientID, page.Mstorgnhirarchyid, page.Recorddifftypeid, page.Recorddiffid, page.Recordstatustypeid, page.Recordstatusid, page.GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetScheduleTabTermnames Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordPlanTabTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Readpermission, &value.Writepermission, &value.FieldID, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetScheduleTabTermvalues(page *entities.RecordcommonEntity) ([]entities.RecordScheduleTabTermnamesEntity, error) {
	logger.Log.Println("In side GetScheduleTabTermnames")
	var sql = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,a.iscompulsory as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,e.readpermission,e.writepermission,f.id,g.recordtrackvalue,c.seq from mststateterm a,mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstrecordfield f,trnreordtracking g WHERE g.clientid=? and g.mstorgnhirarchyid=? AND g.recordid=? AND g.deleteflg=0 and g.activeflg=1 AND g.recordtermid=f.recordtermid AND f.mstrecordfieldtype='SCHEDULE TAB' AND f.deleteflg=0 and f.activeflg=1 AND f.recordtermid = a.recordtermid AND a.recorddifftypeid=? and a.recorddiffid=? AND a.deleteflg=0 and a.activeflg=1 AND a.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND c.id = e.termid AND e.grpid=?  AND e.deleteflg=0 and e.activeflg=1 ORDER BY f.displayseq"
	values := []entities.RecordScheduleTabTermnamesEntity{}
	rows, err := mdao.DB.Query(sql, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Recorddifftypeid, page.Recorddiffid, page.GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetScheduleTabTermnames Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordScheduleTabTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Readpermission, &value.Writepermission, &value.FieldID, &value.Val, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetScheduleTabTermvalueswithStatus(page *entities.RecordcommonEntity) ([]entities.RecordScheduleTabTermnamesEntity, error) {
	logger.Log.Println("In side GetScheduleTabTermnames")
	var sql = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,COALESCE(b.iscompulsory,'0') as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,e.readpermission,e.writepermission,f.id,g.recordtrackvalue,c.seq from mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstrecordfield f,trnreordtracking g,mststateterm a LEFT JOIN mststateterm b ON b.recorddifftypeid=? and b.recorddiffid=? and b.deleteflg=0 and b.activeflg=1 AND a.recordtermid=b.recordtermid WHERE g.clientid=? and g.mstorgnhirarchyid=? AND g.recordid=? AND g.deleteflg=0 and g.activeflg=1 AND g.recordtermid=f.recordtermid AND f.mstrecordfieldtype='SCHEDULE TAB' AND f.deleteflg=0 and f.activeflg=1 AND f.recordtermid = a.recordtermid AND a.recorddifftypeid=? and a.recorddiffid=? AND a.deleteflg=0 and a.activeflg=1 AND a.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND c.id = e.termid AND e.grpid=?  AND e.deleteflg=0 and e.activeflg=1 ORDER BY f.displayseq"
	values := []entities.RecordScheduleTabTermnamesEntity{}
	rows, err := mdao.DB.Query(sql, page.Recordstatustypeid, page.Recordstatusid, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Recorddifftypeid, page.Recorddiffid, page.GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetScheduleTabTermnames Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordScheduleTabTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Readpermission, &value.Writepermission, &value.FieldID, &value.Val, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetPlanTabTermvalues(page *entities.RecordcommonEntity) ([]entities.RecordPlanTabTermnamesEntity, error) {
	logger.Log.Println("In side GetScheduleTabTermnames")
	var sql = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,a.iscompulsory as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,e.readpermission,e.writepermission,f.id,g.recordtrackvalue,c.seq from mststateterm a,mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstrecordfield f,trnreordtracking g WHERE g.clientid=? and g.mstorgnhirarchyid=? AND g.recordid=? AND g.deleteflg=0 and g.activeflg=1 AND g.recordtermid=f.recordtermid AND f.mstrecordfieldtype='PLAN OF ACTION' AND f.deleteflg=0 and f.activeflg=1 AND f.recordtermid = a.recordtermid AND a.recorddifftypeid=? and a.recorddiffid=? AND a.deleteflg=0 and a.activeflg=1 AND a.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND c.id = e.termid AND e.grpid=?  AND e.deleteflg=0 and e.activeflg=1 ORDER BY f.displayseq"
	values := []entities.RecordPlanTabTermnamesEntity{}
	rows, err := mdao.DB.Query(sql, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Recorddifftypeid, page.Recorddiffid, page.GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetScheduleTabTermnames Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordPlanTabTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Readpermission, &value.Writepermission, &value.FieldID, &value.Val, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetPlanTabTermvalueswithStatus(page *entities.RecordcommonEntity) ([]entities.RecordPlanTabTermnamesEntity, error) {
	logger.Log.Println("In side GetScheduleTabTermnames")
	var sql = "select distinct c.id as ID,c.termname as Termname,c.termvalue as Recordtermvalue,COALESCE(b.iscompulsory,'0') as Iscompulsory,d.termtypename as Termtypename,c.termtypeid as Termtypeid,e.readpermission,e.writepermission,f.id,g.recordtrackvalue,c.seq from mstrecordterms c,msttermtype d,mstsupportgrptermmap e,mstrecordfield f,trnreordtracking g,mststateterm a LEFT JOIN mststateterm b ON b.recorddifftypeid=? and b.recorddiffid=? and b.deleteflg=0 and b.activeflg=1 AND a.recordtermid=b.recordtermid WHERE g.clientid=? and g.mstorgnhirarchyid=? AND g.recordid=? AND g.deleteflg=0 and g.activeflg=1 AND g.recordtermid=f.recordtermid AND f.mstrecordfieldtype='PLAN OF ACTION' AND f.deleteflg=0 and f.activeflg=1 AND f.recordtermid = a.recordtermid AND a.recorddifftypeid=? and a.recorddiffid=? AND a.deleteflg=0 and a.activeflg=1 AND a.recordtermid=c.id and c.termtypeid=d.id and a.deleteflg=0 and a.activeflg=1 and c.deleteflg=0 and c.activeflg=1 AND c.id = e.termid AND e.grpid=?  AND e.deleteflg=0 and e.activeflg=1 ORDER BY f.displayseq"
	values := []entities.RecordPlanTabTermnamesEntity{}
	rows, err := mdao.DB.Query(sql, page.Recordstatustypeid, page.Recordstatusid, page.ClientID, page.Mstorgnhirarchyid, page.RecordID, page.Recorddifftypeid, page.Recorddiffid, page.GrpID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetScheduleTabTermnames Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.RecordPlanTabTermnamesEntity{}
		rows.Scan(&value.ID, &value.Termname, &value.Recordtermvalue, &value.Iscompulsory, &value.Termtypename, &value.Termtypeid, &value.Readpermission, &value.Writepermission, &value.FieldID, &value.Val, &value.Seq)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) RemoverecordLink(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, LinkrecordID int64) error {
	var sql = "UPDATE maprecordtolinkrecords SET deleteflg=1 WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND linkrecordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, LinkrecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetLinkRecordsByID(page *entities.LinkRecordEntity) ([]entities.LinkRecordDetailsEntity, error) {
	logger.Log.Println("In side GetLinkRecordsByID")
	var sql = "SELECT a.linkrecordid,b.recordtitle,b.code,d.name FROM maprecordtolinkrecords a,trnrecord b,maprecordtorecorddifferentiation c,mstrecorddifferentiation d WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid=? AND a.activeflg=1 AND a.deleteflg=0 AND a.linkrecordid = b.id AND b.activeflg=1 AND b.deleteflg=0 AND b.id = c.recordid AND c.islatest=1 AND c.recorddifftypeid=2 AND c.activeflg=1 AND c.deleteflg=0 AND c.recorddiffid = d.id"
	values := []entities.LinkRecordDetailsEntity{}
	rows, err := mdao.DB.Query(sql, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLinkRecordsByID Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.LinkRecordDetailsEntity{}
		rows.Scan(&value.LinkrecordID, &value.Recordtitle, &value.Recordcode, &value.Recordtype)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) SaveLinkRecordsByID(ClientID int64, Mstorgnhirarchyid int64, RecordID int64, LinkrecordID int64) (int64, error) {
	//var sql = "UPDATE maprecordtolinkrecords SET deleteflg=1 WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=? AND linkrecordid=?"
	var sql = "INSERT INTO maprecordtolinkrecords(clientid,mstorgnhirarchyid,recordid,linkrecordid) VALUES (?,?,?,?)"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	res, err := stmt.Exec(ClientID, Mstorgnhirarchyid, RecordID, LinkrecordID)
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		logger.Log.Println(err)
		return 0, errors.New("SQL Execution Error")
	}
	return lastInsertedId, nil
}

func (mdao DbConn) CheckLinkRecordFlag(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (int64, error) {
	logger.Log.Println("In side GetAging")
	var countval int64
	var sql = "SELECT count(id) countval FROM maprecordtolinkrecords WHERE clientid=? AND mstorgnhirarchyid=? AND linkrecordid=? AND activeflg=1 AND deleteflg=0"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("CheckLinkRecordFlag Get Statement Prepare Error", err)
		return countval, err
	}
	for rows.Next() {
		err = rows.Scan(&countval)
		logger.Log.Println("CheckLinkRecordFlag rows.next() Error", err)
	}
	return countval, nil
}

func (mdao DbConn) GetParentRecordInfo(ClientID int64, Mstorgnhirarchyid int64, RecordID int64) (entities.ParentRecordInfoEntity, error) {
	logger.Log.Println("In side GetLinkRecordsByID")
	var sql = "SELECT a.parentrecordid,b.recordtitle,b.code,d.name,COALESCE(f.recordtrackvalue,'') palnnedstart,COALESCE(h.recordtrackvalue,'') plannedend FROM trnrecord b,maprecordtorecorddifferentiation c,mstrecorddifferentiation d,mstparentchildmap a LEFT JOIN mstrecordterms e ON e.seq IN (63) AND e.activeflg = 1 AND e.deleteflg = 0 AND e.clientid = ? AND e.mstorgnhirarchyid = ? LEFT JOIN  trnreordtracking f ON e.id = f.recordtermid AND a.parentrecordid = f.recordid LEFT JOIN mstrecordterms g ON g.seq IN (64) AND g.activeflg = 1 AND g.deleteflg = 0 AND g.clientid = ? AND g.mstorgnhirarchyid = ? LEFT JOIN  trnreordtracking h ON g.id =h.recordtermid AND a.parentrecordid = h.recordid WHERE a.clientid = ? AND a.mstorgnhirarchyid = ? AND a.childrecordid = ? AND isattached IN ('Y','N') AND a.activeflg = 1 AND a.deleteflg = 0 AND a.parentrecordid = b.id AND b.activeflg = 1 AND b.deleteflg = 0 AND b.id = c.recordid AND c.islatest = 1 AND c.recorddifftypeid = 3 AND c.activeflg = 1 AND c.deleteflg = 0 AND c.recorddiffid = d.id"
	value := entities.ParentRecordInfoEntity{}
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, ClientID, Mstorgnhirarchyid, ClientID, Mstorgnhirarchyid, RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetLinkRecordsByID Get Statement Prepare Error", err)
		return value, err
	}
	for rows.Next() {
		rows.Scan(&value.ID, &value.Recordtitle, &value.Recordcode, &value.Recordstatus, &value.PlannedStartDate, &value.PlannedEndDate)

	}
	return value, nil
}

func (mdao DbConn) Updatechildrecordflag(ClientID int64, Mstorgnhirarchyid int64, ParentID int64, ChildID int64) error {
	var sql = "UPDATE mstparentchildmap SET isattached='N' WHERE clientid=? AND mstorgnhirarchyid=? AND parentrecordid=? AND childrecordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(ClientID, Mstorgnhirarchyid, ParentID, ChildID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetParentrecordForIM(page *entities.RecordcommonEntity) ([]entities.ParentticketEntity, error) {
	logger.Log.Println("In side GetParentrecordForIM")
	values := []entities.ParentticketEntity{}
	var sql = "SELECT a.parentrecordid,b.code FROM mstparentchildmap a,trnrecord b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.childrecordid=? AND a.deleteflg=0 AND a.activeflg=1 AND a.parentrecordid = b.id and a.isattached in ('Y')"
	rows, err := mdao.DB.Query(sql, page.ClientID, page.Mstorgnhirarchyid, page.RecordID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetParentrecordForIM Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.ParentticketEntity{}
		rows.Scan(&value.ID, &value.Recordnumber)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetFollowupcountForSR(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetFollowupcount")
	var count int64
	var sql = "SELECT count(a.id) count FROM trnreordtracking a,mstrecordterms b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid IN (SELECT childrecordid FROM mstparentchildmap WHERE parentrecordid=? AND deleteflg=0 AND activeflg=1) AND a.recordtermid =b.id AND b.seq=29"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetFollowupcount Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("GetFollowupcount rows.next() Error", err)
	}
	return count, nil
}

func (mdao DbConn) GetOutboundcountForSR(ClientID int64, Mstorgnhirarchyid int64, Recordid int64) (int64, error) {
	logger.Log.Println("In side GetOutboundcount")
	var count int64
	var sql = "SELECT count(a.id) count FROM trnreordtracking a,mstrecordterms b WHERE a.clientid=? AND a.mstorgnhirarchyid=? AND a.recordid IN (SELECT childrecordid FROM mstparentchildmap WHERE parentrecordid=? AND deleteflg=0 AND activeflg=1) AND a.recordtermid =b.id AND b.seq=30"
	rows, err := mdao.DB.Query(sql, ClientID, Mstorgnhirarchyid, Recordid)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetOutboundcount Get Statement Prepare Error", err)
		return count, err
	}
	for rows.Next() {
		err = rows.Scan(&count)
		logger.Log.Println("GetOutboundcount rows.next() Error", err)
	}
	return count, nil
}

// New Addition on 28.04.2022

func UpdateResponseBreachCode(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, responsebreachcode string) error {
	var sql = "UPDATE recordfulldetails SET responsebreachcode=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(responsebreachcode, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateResponseBreachCode(ClientID int64, OrgnID int64, RecordID int64, responsebreachcode string) error {
	var sql = "UPDATE recordfulldetails SET responsebreachcode=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(responsebreachcode, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateResponseBreachComment(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, responsebreachcomment string) error {
	var sql = "UPDATE recordfulldetails SET responsebreachcomment=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(responsebreachcomment, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateResponseBreachComment(ClientID int64, OrgnID int64, RecordID int64, responsebreachcomment string) error {
	var sql = "UPDATE recordfulldetails SET responsebreachcomment=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(responsebreachcomment, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateResolutionBreachCode(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, resolutionbreachcode string) error {
	var sql = "UPDATE recordfulldetails SET resolutionbreachcode=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutionbreachcode, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateResolutionBreachCode(ClientID int64, OrgnID int64, RecordID int64, resolutionbreachcode string) error {
	var sql = "UPDATE recordfulldetails SET resolutionbreachcode=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutionbreachcode, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func UpdateResolutionBreachComment(tx *sql.Tx, ClientID int64, OrgnID int64, RecordID int64, resolutionbreachcomment string) error {
	var sql = "UPDATE recordfulldetails SET resolutionbreachcomment=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutionbreachcomment, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateResolutionBreachComment(ClientID int64, OrgnID int64, RecordID int64, resolutionbreachcomment string) error {
	var sql = "UPDATE recordfulldetails SET resolutionbreachcomment=?  WHERE clientid=? AND mstorgnhirarchyid=? AND recordid=?"
	stmt, err := mdao.DB.Prepare(sql)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(resolutionbreachcomment, ClientID, OrgnID, RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

//New Addition on 28.04.2022

func (mdao DbConn) UpdateRecordTermvalues(rec *entities.RecordcommonEntity) error {
	var sql = "Update trnreordtracking SET recordtrackvalue= ? WHERE recordid=? AND recordtermid=?"
	stmt, err := mdao.DB.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(rec.Termvalue, rec.RecordID, rec.TermID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) UpdateRecordFulldetailsVendorTicketID(rec *entities.RecordcommonEntity) error {
	var sql = "UPDATE recordfulldetails SET vendorticketid=?  WHERE recordid=?"
	stmt, err := mdao.DB.Prepare(sql)

	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Prepare Error")
	}
	defer stmt.Close()
	_, err = stmt.Exec(rec.Termvalue, rec.RecordID)
	if err != nil {
		logger.Log.Println(err)
		return errors.New("SQL Execution Error")
	}

	return nil
}

func (mdao DbConn) GetVendorTicketID(rec *entities.RecordcommonEntity) (int64, error) {
	logger.Log.Println("In side GetVendorTicketID")
	var id int64
	var sql = "SELECT id FROM trnreordtracking WHERE recordid=? AND recordtermid=? AND deleteflg=0 AND activeflg=1"
	rows, err := mdao.DB.Query(sql, rec.RecordID, rec.TermID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetVendorTicketID Get Statement Prepare Error", err)
		return id, err
	}
	for rows.Next() {
		err = rows.Scan(&id)
		logger.Log.Println("GetVendorTicketID rows.next() Error", err)
	}
	return id, nil
}

func (mdao DbConn) GetchildrecordidsForCommon(RecordID int64, ClientID int64, Mstorgnhirarchyid int64, Recorddifftypeid int64, Recorddiffid int64) ([]int64, error) {
	logger.Log.Println("In side GetchildrecordidsForCommon")
	var sql = "select a.childrecordid from mstparentchildmap a,maprecordtorecorddifferentiation b,mstrecorddifferentiation c where a.parentrecordid=? and a.clientid=? and a.mstorgnhirarchyid=? and a.recorddifftypeid=? and a.recorddiffid=? and a.deleteflg=0 and a.activeflg=1 and a.isattached='Y' AND a.childrecordid=b.recordid AND b.recorddifftypeid=3 AND b.islatest=1 AND b.recorddiffid=c.id AND c.seqno NOT IN (3,8,11,14)"
	logger.Log.Println("Query is --++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++---->", sql)
	logger.Log.Println("Parameter is ------>", RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid)

	var childids []int64
	rows, err := mdao.DB.Query(sql, RecordID, ClientID, Mstorgnhirarchyid, Recorddifftypeid, Recorddiffid) //ClientID, Mstorgnhirarchyid,
	//rows, err := mdao.DB.Query(getchildids, RecordID, ClientID, Mstorgnhirarchyid, 2, 4)
	logger.Log.Println("Rows error is  ----+++++++++++++++++++++++++++++++++++++++++-->", err)

	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getchildrecordids Get Statement Prepare Error", err)
		return childids, err
	}
	for rows.Next() {
		var ID int64
		err = rows.Scan(&ID)
		childids = append(childids, ID)
		logger.Log.Println("Getchildrecordids rows.next() Error", err)
	}
	return childids, nil
}
func (mdao DbConn) GetRecordidsToSlaUpdate(recordno []string) ([]entities.RecordInfoEntity, error) {
	logger.Log.Println("In side GetRecordidsToSlaUpdate")
	values := []entities.RecordInfoEntity{}
	var rcnos string
	for i := 0; i < len(recordno); i++ {
		if i > 0 {
			rcnos = rcnos + ","
		}
		rcnos = rcnos + "'" + recordno[i] + "'"

	}
	logger.Log.Println(rcnos)
	sql := "SELECT id as recordid,clientid as clientid,mstorgnhirarchyid as mstorgnhirarchyid,lastupdateddate as datetime FROM trnrecord where code in (" + rcnos + ") and deleteflg=0;"
	rows, err := mdao.DB.Query(sql)
	if err != nil {
		logger.Log.Println("Exception in GetTickeType Prepare Statement..", err)
		return values, err
	}
	defer rows.Close()

	logger.Log.Println("hiiiiiii", rcnos)
	for rows.Next() {
		logger.Log.Println(rows)
		var value entities.RecordInfoEntity
		rows.Scan(&value.RecordID, &value.ClientID, &value.Mstorgnhirarchyid, &value.Datetime)
		values = append(values, value)
	}
	return values, nil
}
func (mdao DbConn) FetchFirstGrpID(ID int64) (int64, error) {
	logger.Log.Println("In side FetchCurrentGrpID")
	var sql = "SELECT a.mstgroupid FROM mstrequesthistory a,maprequestorecord b WHERE b.recordid=? AND a.mainrequestid=b.mstrequestid AND b.activeflg=1 AND b.deleteflg=0 AND a.activeflg=1 AND a.deleteflg=0 order by b.id desc limit 1"
	var grpID int64
	rows, err := mdao.DB.Query(sql, ID)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("FetchCurrentGrpID Get Statement Prepare Error", err)
		return grpID, err
	}
	for rows.Next() {
		rows.Scan(&grpID)

	}
	return grpID, nil
}



func (mdao DbConn) Getfrequentissuesresolver(page *entities.RecordcommonEntity) ([]entities.FrequentRecordEntity, error) {
	logger.Log.Println("In side Getfrequentissues")
	values := []entities.FrequentRecordEntity{}
	rows, err := mdao.DB.Query(frequentresolver, page.Userid, page.ClientID) //, page.Userid, page.Userid, page.ClientID, page.Userid, page.Userid, page.Userid, page.Userid
	logger.Log.Println(frequentresolver)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("Getfrequentissues Get Statement Prepare Error", err)
		return values, err
	}
	for rows.Next() {
		value := entities.FrequentRecordEntity{}
		rows.Scan(&value.Mstorgnhirarchyid, &value.LastlevelID, &value.Count, &value.Recorddifftypeid, &value.Recorddiffid)
		values = append(values, value)
	}
	return values, nil
}

func (mdao DbConn) GetCategoryInfo(CatId int64) (string, string, error) {
	var parentcategoryids string
	var parentcategorynames string
	var sql = "SELECT parentcategoryids,parentcategorynames FROM mstrecorddifferentiation WHERE id=?"
	rows, err := mdao.DB.Query(sql, CatId)
	defer rows.Close()
	if err != nil {
		logger.Log.Println("GetCategoryInfo Get Statement Prepare Error", err)
		return parentcategoryids, parentcategorynames, err
	}
	for rows.Next() {
		err = rows.Scan(&parentcategoryids, &parentcategorynames)
		logger.Log.Println(err)

	}
	return parentcategoryids, parentcategorynames, nil
}
