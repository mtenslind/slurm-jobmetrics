package main

import (
	"fmt"
	"os/user"
	"strconv"
)

type JobInfo struct {
	jobId    int
	userId   int
	userName string
}

func SetJobInfo(jobid int) JobInfo {
	gotid := GetJobUid(jobid)
	userinfo, _ := user.LookupId(strconv.Itoa(gotid))
	gotname := userinfo.Username
	return JobInfo{jobid, gotid, gotname}
}

func main() {

	job := SetJobInfo(66)
	fmt.Printf("Uid %d\n", job.userId)
	fmt.Printf("Username %q\n", job.userName)
}
