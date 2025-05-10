package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
)

const cgroupRoot = "/sys/fs/cgroup/system.slice/slurmstepd.scope/"

type JobInfo struct {
	jobId          int
	userId         int
	userName       string
	cgroupPath     string
	fileDescriptor os.File
	stats          JobStats
}

type JobStats struct {
	memory string
}

// Job struct builder with empty stats
func NewJobStruct(path string, jobId int) JobInfo {
	var jobinfo JobInfo
	var stats JobStats

	userId, userName := GetUserInfo(jobId)

	jobinfo.jobId = jobId
	jobinfo.userId = userId
	jobinfo.userName = userName
	jobinfo.cgroupPath = path
	jobinfo.stats = stats
	return jobinfo

}

func GetUserInfo(jobID int) (int, string) {
	// GetJobUid function calls a CGO function that uses slurm.h
	// Uses an RPC call to the Slurm controller to retreive the jobs' user based on Job ID
	uid := GetJobUid(jobID)
	userinfo, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		panic(err)
	}
	return uid, userinfo.Username
}

func GetStats(root *os.Root, cgroupPath string) JobStats {
	// Takes the CGroup root and parses relevant files for information
	var stats JobStats
	stats.memory = ReadStatFromFile(root, "memory.current")
	return stats

}

func ReadStatFromFile(cgroupFile *os.Root, fileName string) string {
	file, err := cgroupFile.Open(fileName)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	statValue, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(statValue)
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
			fileDescriptor, err := path.OpenRoot(file.Name())
			if err != nil {
				panic(err)
			}
			jobstring, _ := strings.CutPrefix(file.Name(), "job_")
			jobid, _ := strconv.Atoi(jobstring)
			// Check if the job already exists
			// Initializes job with empty stats if not
			if _, ok := jobInfoMap[jobid]; !ok {
				jobInfoMap[jobid] = NewJobStruct(fileDescriptor.Name(), jobid)
			}
			// Set jobinfo to the created struct
			jobinfo := jobInfoMap[jobid]
			jobinfo.stats = GetStats(fileDescriptor, file.Name())
			jobInfoMap[jobid] = jobinfo
		}
	}
}

func main() {

	jobMap := make(map[int]JobInfo)
	FindJobs(jobMap)

	for jobID, jobinfo := range jobMap {
		fmt.Printf("Job id: %d | User name: %s | Memory: %s | path: %s\n", jobID, jobinfo.userName, jobinfo.stats.memory, jobinfo.cgroupPath)
	}

}
