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
	userId     int
	userName   string
	cgroupPath string
	memory     string
}

func GetUserInfo(jobID int) (int, string) {
	uid := GetJobUid(jobID)
	userinfo, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		panic(err)
	}
	return uid, userinfo.Username
}

func GetStats(root *os.Root, name string) string {
	mem, err := root.Open(name + "/memory.current")
	defer mem.Close()
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(mem)
	if err != nil {
		panic(err)
	}
	return string(content)

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
			if _, ok := jobInfoMap[jobid]; !ok {
				userID, userName := GetUserInfo(jobid)
				path := cgroupRoot + file.Name()
				jobInfoMap[jobid] = JobInfo{userID, userName, path, ""}
			}
			jobinfo := jobInfoMap[jobid]
			jobinfo.memory = GetStats(path, file.Name())
			jobInfoMap[jobid] = jobinfo
		}
	}
}

func main() {

	jobMap := make(map[int]JobInfo)
	FindJobs(jobMap)

	for jobID, jobinfo := range jobMap {
		fmt.Printf("Job id: %d | User name: %s | Memory: %s | path: %s\n", jobID, jobinfo.userName, jobinfo.memory, jobinfo.cgroupPath)
	}

}
