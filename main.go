package main

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"strings"
)

const cgroupRoot = "/sys/fs/cgroup/system.slice/slurmstepd.scope/"

type JobInfo struct {
	jobID    int
	userId   int
	userName string
}

func SetJobInfo(jobid int) JobInfo {
	gotid := GetJobUid(jobid)
	userinfo, _ := user.LookupId(strconv.Itoa(gotid))
	gotname := userinfo.Username
	return JobInfo{jobid, gotid, gotname}
}

func FindJobs(jobInfoMap map[int]JobInfo) {
	path, err := os.OpenRoot(cgroupRoot)
	defer path.Close()
	if err != nil {
		panic(err)
	}
	rootfile, err := path.Open("./")
	if err != nil {
		panic(err)
	}
	files, err := rootfile.ReadDir(0)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "job_") {
			jobstring, _ := strings.CutPrefix(file.Name(), "job_")
			jobid, _ := strconv.Atoi(jobstring)
			jobInfoMap[jobid] = SetJobInfo(jobid)

		}
	}
}

func main() {

	jobMap := make(map[int]JobInfo)
	FindJobs(jobMap)

	for jobID, jobinfo := range jobMap {
		fmt.Printf("Job id: %d | User name: %s\n", jobID, jobinfo.userName)
	}

}
