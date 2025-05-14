# slurm-jobmetrics

Slurm-jobmetrics is an agent that reports actual job resource usage for Slurm.

The program looks for Slurm relevant cGroups in `/sys/fs/cgroup/system.slice/slurmstepd.scope` and recovers information periodically. 

The user information relevant to the job is gathered from the Slurm controller via an RPC call. This is done only once for new jobs.

## Requirements

- `slurm-devel` package for `slurm.h`
- `task/cgroup` needs to be enabled in Slurms' `cgroup.conf`

## Usage

- Compile with `go build .`
- Transfer the executable to relevant nodes
- Start the executable in the background
- Logs will be written to `logfile.log` in the current direcotry

## Example output

```
2025/05/14 12:09:07 Job_ID=55153156;User=testuser;User_ID=12345;MemReq=68719476736;MemUsed=4898816;CpuReq=28;CpuStat=1080311
2025/05/14 12:09:17 Job_ID=55153156;User=testuser;User_ID=12345;MemReq=68719476736;MemUsed=4403200;CpuReq=28;CpuStat=1082721
2025/05/14 12:09:27 Job_ID=55153156;User=testuser;User_ID=12345;MemReq=68719476736;MemUsed=4426018816;CpuReq=28;CpuStat=142723141
2025/05/14 12:09:37 Job_ID=55153156;User=testuser;User_ID=12345;MemReq=68719476736;MemUsed=4152467456;CpuReq=28;CpuStat=418001796
```
