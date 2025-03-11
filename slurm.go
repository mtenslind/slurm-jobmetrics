package main

/*
#cgo LDFLAGS: -lslurm
#include <stdlib.h>
#include <slurm/slurm.h>

int get_uid_from_job(uint32_t job_id) {
    slurm_init(NULL);
    job_info_msg_t *job_buffer = NULL;
    slurm_job_info_t *job_info = NULL;
    int rc = slurm_load_job(&job_buffer, job_id, SHOW_LOCAL);
    if (rc != 0) {
        return rc;  // Slurm API call failed
    }
    job_info = job_buffer->job_array;
    return job_info->user_id;

    slurm_free_job_info_msg(job_buffer);
    return 1; // Job ID not found
}
*/
import "C"

func GetJobUid(jobid int) int {

	uid := C.get_uid_from_job(C.uint32_t(jobid))
	return int(uid)
}
