package main

import (
	"fmt"
	"os/user"
	"strconv"
	"strings"

	// "log"
	// "time"
	"os"
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

func FindJobs(cgroupPath string) []JobInfo {
	path, err := os.OpenRoot(cgroupPath)
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
	var jobs []JobInfo
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "job_") {
			jobstring, _ := strings.CutPrefix(file.Name(), "job_")
			jobid, _ := strconv.Atoi(jobstring)
			jobs = append(jobs, SetJobInfo(jobid))
		}
	}
	return jobs
}

func main() {
	jobs := FindJobs("/sys/fs/cgroup/system.slice/slurmstepd.scope/")

	for _, job := range jobs {
		fmt.Printf("Job id: %d | User name: %s\n", job.jobId, job.userName)
	}

}
